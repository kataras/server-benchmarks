// Package runner executes benchmark tests: it starts each environment's
// server application, waits for it to accept connections, measures it with
// the bombardier load generator, validates the outcome and tears the
// server's whole process tree down again.
//
// A failing environment never aborts the run: it is recorded as
// StatusFailed and the runner continues with the next one, returning an
// aggregate error at the end.
package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/spec"
)

// Default option values used by New when the corresponding Options field
// is zero.
const (
	// DefaultReadyTimeout is how long a server may take to accept
	// connections after being started.
	DefaultReadyTimeout = 30 * time.Second
	// DefaultStopGrace is the grace period between asking a server to
	// terminate and force-killing its process tree.
	DefaultStopGrace = 5 * time.Second
	// DefaultPortFreeTimeout is how long to wait after stopping a server
	// for its port to be released.
	DefaultPortFreeTimeout = 10 * time.Second
)

// Options configures a Runner.
type Options struct {
	// Bombardier executes the load generator. The zero value uses the
	// "bombardier" binary from PATH.
	Bombardier bombardier.Client
	// WaitBetweenRuns is an idle pause between two environments of the
	// same test, letting the machine settle. Zero means no pause.
	WaitBetweenRuns time.Duration
	// ReadyTimeout overrides DefaultReadyTimeout when positive.
	ReadyTimeout time.Duration
	// StopGrace overrides DefaultStopGrace when positive.
	StopGrace time.Duration
	// Logger receives progress and diagnostics. Nil discards everything.
	Logger *slog.Logger
}

// Runner executes benchmark tests sequentially — measurements must not
// compete with each other for CPU, so there is deliberately no
// parallelism across environments.
type Runner struct {
	opts Options
	log  *slog.Logger
}

// New returns a Runner with defaults applied.
func New(opts Options) *Runner {
	if opts.ReadyTimeout <= 0 {
		opts.ReadyTimeout = DefaultReadyTimeout
	}

	if opts.StopGrace <= 0 {
		opts.StopGrace = DefaultStopGrace
	}

	log := opts.Logger
	if log == nil {
		log = slog.New(slog.DiscardHandler)
	}

	return &Runner{opts: opts, log: log}
}

// Run benchmarks every test and returns one report per test. Reports are
// returned even when some environments failed; the aggregate error joins
// all per-environment failures. When ctx is canceled the run stops early
// and the error includes ctx's cause.
func (r *Runner) Run(ctx context.Context, tests []spec.Test) ([]TestReport, error) {
	reports := make([]TestReport, 0, len(tests))
	var errs []error

	for _, t := range tests {
		report, err := r.runTest(ctx, t)
		reports = append(reports, report)
		if err != nil {
			errs = append(errs, err)
		}

		if ctx.Err() != nil {
			errs = append(errs, ctx.Err())
			break
		}
	}

	return reports, errors.Join(errs...)
}

// runTest benchmarks one test across all its environments.
func (r *Runner) runTest(ctx context.Context, t spec.Test) (TestReport, error) {
	report := TestReport{Test: t, Envs: make([]EnvReport, 0, len(t.Envs))}
	var errs []error

	failAll := func(err error) (TestReport, error) {
		for _, env := range t.Envs {
			if !env.CanBenchmark() {
				report.Envs = append(report.Envs, EnvReport{Env: env, Status: StatusSkipped})
				continue
			}
			report.Envs = append(report.Envs, EnvReport{Env: env, Status: StatusFailed, Err: err})
		}
		sortEnvReports(report.Envs)
		return report, err
	}

	var bodyPath string
	if spec.IsURL(t.BodyFile) {
		tmpDir, err := os.MkdirTemp("", "server-benchmarks-")
		if err != nil {
			return failAll(fmt.Errorf("test %s: temp dir: %w", t.Name, err))
		}
		defer os.RemoveAll(tmpDir)

		r.log.Info("downloading request body", "test", t.Name, "url", t.BodyFile)
		bodyPath, err = fetchBodyFile(ctx, t.BodyFile, tmpDir)
		if err != nil {
			return failAll(fmt.Errorf("test %s: %w", t.Name, err))
		}
	}

	args := bombardier.Args(t, bodyPath)

	ranAny := false
	for _, env := range t.Envs {
		if ctx.Err() != nil {
			break
		}

		if !env.CanBenchmark() {
			r.log.Debug("skipping", "test", t.Name, "env", env.Name)
			report.Envs = append(report.Envs, EnvReport{Env: env, Status: StatusSkipped})
			continue
		}

		if ranAny && r.opts.WaitBetweenRuns > 0 {
			r.log.Debug("settling", "wait", r.opts.WaitBetweenRuns.String())
			if sleepCtx(ctx, r.opts.WaitBetweenRuns) != nil {
				break
			}
		}
		ranAny = true

		r.log.Info("benchmarking", "test", t.Name, "env", env.Name, "dir", env.Dir)

		result, err := r.runEnv(ctx, t, env, args)
		if err != nil {
			if ctx.Err() != nil {
				break // canceled mid-run: the partial result is meaningless.
			}

			err = fmt.Errorf("%s: %s: %w", t.Name, env.Name, err)
			r.log.Error("benchmark failed", "test", t.Name, "env", env.Name, "error", err.Error())
			report.Envs = append(report.Envs, EnvReport{Env: env, Status: StatusFailed, Err: err})
			errs = append(errs, err)
			continue
		}

		r.log.Info("done", "test", t.Name, "env", env.Name,
			"rps", fmt.Sprintf("%.0f", result.RequestsPerSecond.Mean),
			"latency", bombardier.FormatTimeUs(result.Latency.Mean))
		report.Envs = append(report.Envs, EnvReport{Env: env, Status: StatusOK, Result: &result})
	}

	sortEnvReports(report.Envs)

	if winner, ok := report.Winner(); ok {
		r.log.Info("winner", "test", t.Name, "env", winner.Env.Name,
			"rps", fmt.Sprintf("%.0f", winner.Result.RequestsPerSecond.Mean))
	}

	return report, errors.Join(errs...)
}

// runEnv executes the full lifecycle for a single environment:
// port check, setup commands, server start, readiness wait, measurement,
// validation, teardown and port release.
func (r *Runner) runEnv(ctx context.Context, t spec.Test, env spec.Env, args []string) (bombardier.Result, error) {
	var zero bombardier.Result

	addr, err := targetAddr(t.URL)
	if err != nil {
		return zero, err
	}

	if info, statErr := os.Stat(env.Dir); statErr != nil {
		return zero, fmt.Errorf("code directory is not accessible: %w", statErr)
	} else if !info.IsDir() {
		return zero, fmt.Errorf("code directory %q is not a directory", env.Dir)
	}

	if portBusy(addr) {
		return zero, fmt.Errorf("%s is already in use before starting the server; is a previous one still running?", addr)
	}

	lines := execLines(env.Exec)
	if len(lines) == 0 {
		return zero, errors.New("empty Exec")
	}

	output := newTail(8 << 10)
	var out io.Writer = output
	if r.log.Enabled(ctx, slog.LevelDebug) {
		out = io.MultiWriter(output, logWriter{log: r.log, tag: env.Name})
	}

	for _, line := range lines[:len(lines)-1] {
		r.log.Debug("setup", "env", env.Name, "cmd", line)
		if err := runCommand(ctx, env.Dir, line, r.opts.StopGrace, out); err != nil {
			return zero, fmt.Errorf("setup command %q: %w%s", line, err, outputSuffix(output))
		}
	}

	serverLine := lines[len(lines)-1]
	r.log.Debug("starting server", "env", env.Name, "cmd", serverLine)

	proc, err := startProcess(ctx, env.Dir, serverLine, r.opts.StopGrace, out)
	if err != nil {
		return zero, fmt.Errorf("start server %q: %w", serverLine, err)
	}

	result, benchErr := func() (bombardier.Result, error) {
		if err := waitReady(ctx, addr, proc.exited, r.opts.ReadyTimeout); err != nil {
			if errors.Is(err, errServerExited) {
				return zero, fmt.Errorf("%w (%v)%s", err, proc.waitErr, outputSuffix(output))
			}
			return zero, err
		}

		r.log.Debug("measuring", "env", env.Name, "cmd", r.opts.Bombardier.String(args))

		result, err := r.opts.Bombardier.Run(ctx, args)
		if err != nil {
			return zero, err
		}

		return result, validate(t, env, result)
	}()

	var stopErr error
	if err := proc.stop(); err != nil {
		stopErr = fmt.Errorf("stop server: %w", err)
	}

	var portErr error
	if ctx.Err() == nil {
		if err := waitPortFree(ctx, addr, DefaultPortFreeTimeout); err != nil && ctx.Err() == nil {
			portErr = err
		}
	}

	if err := errors.Join(benchErr, stopErr, portErr); err != nil {
		return zero, err
	}

	return result, nil
}

// ValidationError reports a benchmark whose HTTP results did not meet the
// test's expectations (transport errors, non-2xx responses or a request
// count mismatch).
type ValidationError struct {
	Test   string
	Env    string
	Reason string
	Result bombardier.Result
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s [2xx: %d, 3xx: %d, 4xx: %d, 5xx: %d, others: %d]",
		e.Reason, e.Result.Req2XX, e.Result.Req3XX, e.Result.Req4XX, e.Result.Req5XX, e.Result.Others)
}

// validate checks a bombardier result against the test's expectations.
func validate(t spec.Test, env spec.Env, result bombardier.Result) error {
	fail := func(reason string) error {
		return &ValidationError{Test: t.Name, Env: env.Name, Reason: reason, Result: result}
	}

	if descriptions := result.ErrorDescriptions(); len(descriptions) > 0 {
		return fail("transport errors: " + strings.Join(descriptions, ", "))
	}

	if n := result.Req4XX + result.Req5XX + result.Others; n > 0 {
		return fail(fmt.Sprintf("%d non-successful responses", n))
	}

	if want := t.NumberOfRequests; want > 0 && result.Req2XX != want {
		return fail(fmt.Sprintf("expected %d successful requests, got %d", want, result.Req2XX))
	}

	return nil
}

// outputSuffix renders captured process output for inclusion in an error
// message, or nothing when there was no output.
func outputSuffix(t *tail) string {
	if s := t.String(); s != "" {
		return " — output: " + s
	}

	return ""
}

// sleepCtx sleeps for d unless ctx is canceled first.
func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
