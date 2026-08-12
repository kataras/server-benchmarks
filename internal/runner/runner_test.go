package runner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/spec"
)

// Paths of the helper binaries compiled once in TestMain.
var (
	fakeServerBin     string
	fakeBombardierBin string
)

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "runner-testbins-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	build := func(name, pkg string) (string, error) {
		out := filepath.Join(tmp, name)
		if runtime.GOOS == "windows" {
			out += ".exe"
		}
		cmd := exec.Command("go", "build", "-o", out, pkg)
		if outBytes, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v: %s", pkg, err, outBytes)
		}
		return out, nil
	}

	code := func() int {
		if fakeServerBin, err = build("fakeserver", "./testdata/fakeserver"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if fakeBombardierBin, err = build("fakebombardier", "./testdata/fakebombardier"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return m.Run()
	}()

	os.RemoveAll(tmp)
	os.Exit(code)
}

// newRunner returns a Runner wired to the fake bombardier, with timeouts
// tightened for tests.
func newRunner(t *testing.T) *Runner {
	t.Helper()

	return New(Options{
		Bombardier:   bombardier.Client{Bin: fakeBombardierBin},
		ReadyTimeout: 15 * time.Second,
		StopGrace:    2 * time.Second,
	})
}

// serverEnv builds an env that runs the fake server in the given mode.
// The binary path is quoted: temp dirs may contain spaces, which also
// exercises the quote-aware command parser on every test.
func serverEnv(t *testing.T, name, addr, mode string) spec.Env {
	t.Helper()

	return spec.Env{
		Name:     name,
		Language: "Go",
		Dir:      filepath.Dir(fakeServerBin),
		Exec:     fmt.Sprintf("%q -addr %s -mode %s", fakeServerBin, addr, mode),
	}
}

// staticTest builds a small request-count test targeting addr.
func staticTest(addr string, envs ...spec.Env) spec.Test {
	return spec.Test{
		Name:                "Static",
		NumberOfConnections: 10,
		NumberOfRequests:    1000,
		Timeout:             spec.Duration(5 * time.Second),
		Method:              "GET",
		URL:                 "http://" + addr,
		Envs:                envs,
	}
}

// withCounter points the fake bombardier at a fresh invocation counter,
// making successive environments report decreasing RPS.
func withCounter(t *testing.T) {
	t.Helper()
	t.Setenv("FAKEBOMB_COUNTER", filepath.Join(t.TempDir(), "counter"))
}

func TestRunHappyPathSortingAndSkipped(t *testing.T) {
	withCounter(t)
	addr := freeAddr(t)

	test := staticTest(addr,
		serverEnv(t, "Alpha", addr, "serve"), // counter 0: rps 200000
		serverEnv(t, "Beta", addr, "serve"),  // counter 1: rps 150000
		spec.Env{Name: "Gamma", Dir: ".", Exec: "unused", NotSupported: true},
	)

	reports, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err != nil {
		t.Fatal(err)
	}

	if len(reports) != 1 {
		t.Fatalf("got %d reports, want 1", len(reports))
	}

	envs := reports[0].Envs
	if len(envs) != 3 {
		t.Fatalf("got %d env reports, want 3", len(envs))
	}

	// Sorted: OK by descending RPS, skipped last. No nil Result is ever
	// dereferenced even with a skipped env present (regression: the old
	// sort comparator panicked here).
	for i, want := range []struct {
		name   string
		status Status
	}{{"Alpha", StatusOK}, {"Beta", StatusOK}, {"Gamma", StatusSkipped}} {
		if envs[i].Env.Name != want.name || envs[i].Status != want.status {
			t.Errorf("envs[%d] = %s/%s, want %s/%s", i, envs[i].Env.Name, envs[i].Status, want.name, want.status)
		}
	}

	if envs[0].Result.RequestsPerSecond.Mean <= envs[1].Result.RequestsPerSecond.Mean {
		t.Error("envs are not sorted by descending RPS")
	}

	winner, ok := reports[0].Winner()
	if !ok || winner.Env.Name != "Alpha" {
		t.Errorf("winner = %v/%v, want Alpha", winner.Env.Name, ok)
	}

	// The server must be gone: the port has to be free again.
	if portBusy(addr) {
		t.Error("server still listening after the run")
	}
}

func TestRunContinuesAfterFailure(t *testing.T) {
	withCounter(t)
	t.Setenv("FAKEBOMB_FAIL_FIRST", "1")
	addr := freeAddr(t)

	test := staticTest(addr,
		serverEnv(t, "Failing", addr, "serve"), // first invocation reports 4xx.
		serverEnv(t, "Passing", addr, "serve"),
	)

	reports, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil {
		t.Fatal("expected an aggregate error")
	}

	var vErr *ValidationError
	if !errors.As(err, &vErr) {
		t.Fatalf("error chain %v does not contain a *ValidationError", err)
	}
	if vErr.Result.Req4XX == 0 {
		t.Errorf("validation error result = %+v, want 4xx > 0", vErr.Result)
	}

	envs := reports[0].Envs
	if len(envs) != 2 {
		t.Fatalf("got %d env reports, want 2", len(envs))
	}

	// Successful env sorts first, the failed one after it — and the
	// failure did not abort the rest of the run.
	if envs[0].Env.Name != "Passing" || envs[0].Status != StatusOK {
		t.Errorf("envs[0] = %s/%s, want Passing/ok", envs[0].Env.Name, envs[0].Status)
	}
	if envs[1].Env.Name != "Failing" || envs[1].Status != StatusFailed {
		t.Errorf("envs[1] = %s/%s, want Failing/failed", envs[1].Env.Name, envs[1].Status)
	}
	if envs[1].Err == nil {
		t.Error("failed env has no error attached")
	}
}

func TestRunServerCrashSurfacesOutput(t *testing.T) {
	addr := freeAddr(t)
	test := staticTest(addr, serverEnv(t, "Crasher", addr, "crash"))

	reports, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil {
		t.Fatal("expected an error")
	}

	if !errors.Is(err, errServerExited) {
		t.Errorf("error = %v, want errServerExited in the chain", err)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error = %v, want it to include the process output", err)
	}

	if got := reports[0].Envs[0].Status; got != StatusFailed {
		t.Errorf("status = %s, want failed", got)
	}
}

func TestRunReadyTimeout(t *testing.T) {
	addr := freeAddr(t)
	test := staticTest(addr, serverEnv(t, "Staller", addr, "stall"))

	r := New(Options{
		Bombardier:   bombardier.Client{Bin: fakeBombardierBin},
		ReadyTimeout: 700 * time.Millisecond,
		StopGrace:    2 * time.Second,
	})

	_, err := r.Run(context.Background(), []spec.Test{test})
	if err == nil || !strings.Contains(err.Error(), "did not become ready") {
		t.Fatalf("error = %v, want a readiness timeout", err)
	}
}

func TestRunPortConflictFailsFast(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	addr := listener.Addr().String()

	test := staticTest(addr, serverEnv(t, "Blocked", addr, "serve"))

	_, err = newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil || !strings.Contains(err.Error(), "already in use") {
		t.Fatalf("error = %v, want a port-conflict error", err)
	}
}

func TestRunKillsWholeProcessTree(t *testing.T) {
	withCounter(t)
	addr := freeAddr(t)

	// spawn mode: the listening server is a *grandchild* of the runner,
	// exactly like `go run .`. After the run, the whole tree must be dead
	// (regression: on Linux the old code killed only the direct child and
	// the real server kept the port forever).
	test := staticTest(addr, serverEnv(t, "Wrapped", addr, "spawn"))

	if _, err := newRunner(t).Run(context.Background(), []spec.Test{test}); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for portBusy(addr) {
		if time.Now().After(deadline) {
			t.Fatal("grandchild server still listening after the run — process tree was not killed")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestRunContextCancelStopsEverything(t *testing.T) {
	addr := freeAddr(t)
	test := staticTest(addr, serverEnv(t, "Staller", addr, "stall"))

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := New(Options{
		Bombardier:   bombardier.Client{Bin: fakeBombardierBin},
		ReadyTimeout: time.Minute,
		StopGrace:    2 * time.Second,
	}).Run(ctx, []spec.Test{test})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(start); elapsed > 15*time.Second {
		t.Fatalf("cancellation took %s, teardown is stuck", elapsed)
	}
}

func TestRunMissingDirFailsEnvOnly(t *testing.T) {
	withCounter(t)
	addr := freeAddr(t)

	missing := spec.Env{
		Name: "Ghost",
		Dir:  filepath.Join(t.TempDir(), "does-not-exist"),
		Exec: "irrelevant",
	}

	test := staticTest(addr, missing, serverEnv(t, "Alive", addr, "serve"))

	reports, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil || !strings.Contains(err.Error(), "not accessible") {
		t.Fatalf("error = %v, want a missing-directory error", err)
	}

	envs := reports[0].Envs
	if envs[0].Env.Name != "Alive" || envs[0].Status != StatusOK {
		t.Errorf("envs[0] = %s/%s, want Alive/ok — a missing dir must not abort the others", envs[0].Env.Name, envs[0].Status)
	}
	if envs[1].Env.Name != "Ghost" || envs[1].Status != StatusFailed {
		t.Errorf("envs[1] = %s/%s, want Ghost/failed", envs[1].Env.Name, envs[1].Status)
	}
}

func TestRunBadBombardierOutput(t *testing.T) {
	for _, mode := range []string{"noresult", "garbage"} {
		t.Run(mode, func(t *testing.T) {
			t.Setenv("FAKEBOMB_MODE", mode)
			addr := freeAddr(t)
			test := staticTest(addr, serverEnv(t, "Server", addr, "serve"))

			reports, err := newRunner(t).Run(context.Background(), []spec.Test{test})
			if err == nil || !strings.Contains(err.Error(), "parse bombardier output") {
				t.Fatalf("error = %v, want a parse error", err)
			}
			// Regression: nil Result must never be dereferenced afterwards.
			if got := reports[0].Envs[0].Status; got != StatusFailed {
				t.Errorf("status = %s, want failed", got)
			}
		})
	}
}

func TestRunTransportErrors(t *testing.T) {
	t.Setenv("FAKEBOMB_MODE", "transport-error")
	addr := freeAddr(t)
	test := staticTest(addr, serverEnv(t, "Flaky", addr, "serve"))

	_, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil || !strings.Contains(err.Error(), "forcibly closed") {
		t.Fatalf("error = %v, want the transport error description", err)
	}
}

func TestRunRequestCountMismatch(t *testing.T) {
	t.Setenv("FAKEBOMB_MODE", "mismatch")
	addr := freeAddr(t)
	test := staticTest(addr, serverEnv(t, "Short", addr, "serve"))

	_, err := newRunner(t).Run(context.Background(), []spec.Test{test})
	if err == nil || !strings.Contains(err.Error(), "expected 1000 successful requests, got 999") {
		t.Fatalf("error = %v, want a request-count mismatch", err)
	}
}
