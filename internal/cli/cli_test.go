package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/report"
	"github.com/kataras/server-benchmarks/internal/runner"
	"github.com/kataras/server-benchmarks/internal/spec"
)

const sampleSpec = `
- Name: Static
  Description: Fires {{.NumberOfRequests}} requests.
  NumberOfRequests: 1000
  Method: GET
  URL: http://localhost:5000
  Envs:
    - Repo: kataras/iris
    - Repo: expressjs/express
      Language: Javascript
      NotYetImplemented: true
`

func writeSpec(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "tests.yml")
	if err := os.WriteFile(path, []byte(sampleSpec), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func run(t *testing.T, args ...string) (code int, stdout, stderr string) {
	t.Helper()

	var out, errOut strings.Builder
	code = Main(context.Background(), args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestUnknownCommand(t *testing.T) {
	code, _, stderr := run(t, "bogus")
	if code != ExitUsage {
		t.Errorf("exit = %d, want %d", code, ExitUsage)
	}
	if !strings.Contains(stderr, "unknown command") || !strings.Contains(stderr, "Usage:") {
		t.Errorf("stderr = %q, want the usage text", stderr)
	}
}

func TestHelp(t *testing.T) {
	code, stdout, _ := run(t, "help")
	if code != ExitOK {
		t.Errorf("exit = %d, want 0", code)
	}
	for _, want := range []string{"run", "list", "export", "version"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help does not mention %q", want)
		}
	}
}

func TestVersion(t *testing.T) {
	code, stdout, _ := run(t, "version")
	if code != ExitOK {
		t.Errorf("exit = %d, want 0", code)
	}
	if !strings.Contains(stdout, "server-benchmarks") || !strings.Contains(stdout, "go1.") {
		t.Errorf("version output = %q", stdout)
	}
}

func TestList(t *testing.T) {
	code, stdout, stderr := run(t, "list", "-i", writeSpec(t))
	if code != ExitOK {
		t.Fatalf("exit = %d, stderr = %q", code, stderr)
	}

	for _, want := range []string{
		"Static: GET http://localhost:5000 (1000 requests, 125 connections)",
		"Iris", "Express", "not yet implemented",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("list output does not contain %q:\n%s", want, stdout)
		}
	}
}

func TestListSelector(t *testing.T) {
	code, stdout, _ := run(t, "list", "-i", writeSpec(t), "-t", "static.iris")
	if code != ExitOK {
		t.Fatalf("exit = %d", code)
	}
	if strings.Contains(stdout, "Express") {
		t.Errorf("selector did not filter envs:\n%s", stdout)
	}
}

func TestRunUsageErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"bad flag", []string{"run", "-bogus"}, ""},
		{"missing spec", []string{"run", "-i", filepath.Join(t.TempDir(), "nope.yml")}, "open spec"},
		{"unmatched selector", []string{"run", "-i", writeSpec(t), "-t", "nope"}, "no test or env matches"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, _, stderr := run(t, tc.args...)
			if code != ExitUsage {
				t.Errorf("exit = %d, want %d (stderr: %q)", code, ExitUsage, stderr)
			}
			if tc.want != "" && !strings.Contains(stderr, tc.want) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tc.want)
			}
		})
	}
}

func TestRunMissingBombardier(t *testing.T) {
	code, _, stderr := run(t,
		"run", "-i", writeSpec(t),
		"-bombardier", filepath.Join(t.TempDir(), "definitely-not-bombardier"),
	)
	if code != ExitFailure {
		t.Errorf("exit = %d, want %d", code, ExitFailure)
	}
	if !strings.Contains(stderr, "not found") {
		t.Errorf("stderr = %q, want a not-found error with install hint", stderr)
	}
}

func TestDefaultSubcommandIsRun(t *testing.T) {
	// Bare flags (no subcommand) must dispatch to run — the Docker
	// ENTRYPOINT depends on it. The missing spec proves run handled it.
	code, _, stderr := run(t, "-i", filepath.Join(t.TempDir(), "nope.yml"))
	if code != ExitUsage || !strings.Contains(stderr, "open spec") {
		t.Errorf("exit = %d, stderr = %q; want run to handle bare flags", code, stderr)
	}
}

func TestExport(t *testing.T) {
	dir := t.TempDir()

	data := report.Data{
		GeneratedAt: time.Date(2026, 8, 12, 10, 0, 0, 0, time.UTC),
		Tests: []runner.TestReport{{
			Test: spec.Test{Name: "Static", Description: "d"},
			Envs: []runner.EnvReport{{
				Env:    spec.Env{Name: "Iris", Language: "Go", Link: "https://github.com/kataras/iris"},
				Status: runner.StatusOK,
				Result: &bombardier.Result{
					TimeTakenSeconds:  1,
					Req2XX:            1000,
					Latency:           bombardier.Stats{Mean: 500},
					RequestsPerSecond: bombardier.Stats{Mean: 1000},
				},
			}},
		}},
	}
	if err := report.Write(data, report.Options{Dir: dir}); err != nil {
		t.Fatal(err)
	}

	exportDir := t.TempDir()
	code, _, stderr := run(t, "export", "-i", filepath.Join(dir, "results.json"), "-o", exportDir)
	if code != ExitOK {
		t.Fatalf("exit = %d, stderr = %q", code, stderr)
	}

	for _, want := range []string{"RESULTS.md", "static.csv", "static.svg", "results.json"} {
		if _, err := os.Stat(filepath.Join(exportDir, want)); err != nil {
			t.Errorf("export did not produce %s: %v", want, err)
		}
	}
}

func TestExportBadInput(t *testing.T) {
	code, _, _ := run(t, "export", "-i", filepath.Join(t.TempDir(), "nope.json"))
	if code != ExitUsage {
		t.Errorf("exit = %d, want %d", code, ExitUsage)
	}
}

func TestRunInterrupted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already canceled: the run must bail out as interrupted.

	var out, errOut strings.Builder
	code := Main(ctx, []string{"run", "-i", writeSpec(t), "-bombardier", "go"}, &out, &errOut)
	if code != ExitInterrupted {
		t.Errorf("exit = %d, want %d (stderr: %q)", code, ExitInterrupted, errOut.String())
	}
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatal("sanity: ctx not canceled")
	}
}
