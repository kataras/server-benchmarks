package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/report"
	"github.com/kataras/server-benchmarks/internal/runner"
	"github.com/kataras/server-benchmarks/internal/spec"
	"github.com/kataras/server-benchmarks/internal/sysinfo"
)

// runCmd implements `server-benchmarks run` (the default subcommand):
// load the spec, benchmark every selected test and render the reports.
func runCmd(ctx context.Context, args []string, stderr io.Writer) int {
	fs := newFlagSet("run", stderr)
	specFile := fs.String("i", "./tests.yml", "path of the YAML spec file describing the tests")
	outDir := fs.String("o", "./", "directory for the generated reports")
	waitRun := fs.Duration("wait-run", 3*time.Second, "idle time between two benchmarks, letting the machine settle")
	readme := fs.Bool("readme", false, "also generate a README.md next to RESULTS.md")
	bin := fs.String("bombardier", "", `bombardier executable (defaults to "bombardier" from PATH)`)
	verbose := fs.Bool("v", false, "verbose (debug) logging")

	var selectors stringSlice
	fs.Var(&selectors, "t", "run only this test or test.env (repeatable), e.g. -t rest -t static.iris")

	if err := fs.Parse(args); err != nil {
		return parseCode(err)
	}

	tests, err := spec.LoadFile(*specFile)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitUsage
	}

	tests, unmatched := spec.Filter(tests, selectors)
	if len(unmatched) > 0 {
		fmt.Fprintf(stderr, "no test or env matches: %s\n", strings.Join(unmatched, ", "))
		return ExitUsage
	}

	client := bombardier.Client{Bin: *bin}
	if _, err := client.LookPath(); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFailure
	}

	logger := newLogger(stderr, *verbose)

	started := time.Now()
	reports, runErr := runner.New(runner.Options{
		Bombardier:      client,
		WaitBetweenRuns: *waitRun,
		Logger:          logger,
	}).Run(ctx, tests)

	if ctx.Err() != nil {
		// The user aborted; leave no half-written reports behind.
		fmt.Fprintln(stderr, "interrupted")
		return ExitInterrupted
	}

	data := report.Data{
		GeneratedAt: time.Now(),
		System:      sysinfo.Collect(ctx, logger),
		Tests:       reports,
	}

	// Render even when some environments failed: a partially failed run
	// (recorded in the report itself) must not evaporate.
	if err := report.Write(data, report.Options{Dir: *outDir, Readme: *readme}); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFailure
	}

	logger.Info("reports written", "dir", *outDir, "took", time.Since(started).Round(time.Second).String())

	if runErr != nil {
		fmt.Fprintln(stderr, runErr)
		return ExitFailure
	}

	return ExitOK
}
