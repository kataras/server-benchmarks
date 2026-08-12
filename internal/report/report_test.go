package report

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/runner"
	"github.com/kataras/server-benchmarks/internal/spec"
	"github.com/kataras/server-benchmarks/internal/sysinfo"
)

var update = flag.Bool("update", false, "rewrite the golden files")

// goldenFiles lists every artifact sampleData produces.
var goldenFiles = []string{
	"RESULTS.md",
	"README.md",
	"results.json",
	"static.csv",
	"static_latency.csv",
	"static.svg",
	"static_latency.svg",
	"rest.csv",
	"rest_latency.csv",
	"rest.svg",
	"rest_latency.svg",
}

func okResult(rps, latency float64, requests uint64) *bombardier.Result {
	return &bombardier.Result{
		BytesRead:        int64(requests) * 150,
		BytesWritten:     int64(requests) * 40,
		TimeTakenSeconds: float64(requests) / rps,
		Req2XX:           requests,
		Latency:          bombardier.Stats{Mean: latency, Stddev: latency / 4, Max: latency * 20},
		RequestsPerSecond: bombardier.Stats{
			Mean: rps, Stddev: rps / 10, Max: rps * 1.2,
		},
	}
}

func env(name, language, repo string) spec.Env {
	e := spec.Env{Name: name, Language: language, Repo: repo}
	if repo != "" {
		e.Link = "https://github.com/" + repo
	}
	return e
}

// sampleData is a fixed, deterministic report input covering every row
// kind: successful, failed and skipped environments, and an env without a
// public repository (no link).
func sampleData() Data {
	return Data{
		GeneratedAt: time.Date(2026, 8, 12, 14, 30, 0, 0, time.UTC),
		System: sysinfo.Info{
			Processor:  "Snapdragon X Elite X1E80100",
			RAM:        "31.69 GB",
			OS:         "Microsoft Windows 11 Home",
			Bombardier: "v2.0.2",
			Go:         "go1.26.5",
			Dotnet:     "10.0.100",
			Node:       "v24.4.1",
		},
		Tests: []runner.TestReport{
			{
				Test: spec.Test{
					Name:        "Static",
					Description: "Fires 100000 requests, receives a static message as response.",
				},
				Envs: []runner.EnvReport{
					{Env: env("Iris", "Go", "kataras/iris"), Status: runner.StatusOK, Result: okResult(284059.44, 438.94, 100000)},
					{Env: env("Fastify", "Javascript", "fastify/fastify"), Status: runner.StatusOK, Result: okResult(180010.12, 690.12, 100000)},
					{Env: env("ASP.NET Core", "C#", "dotnet/aspnetcore"), Status: runner.StatusOK, Result: okResult(150000.55, 830.77, 100000)},
					{Env: env("Iris (next)", "Go", ""), Status: runner.StatusOK, Result: okResult(120044.32, 990.01, 100000)},
					{Env: env("Gin", "Go", "gin-gonic/gin"), Status: runner.StatusFailed, Err: errors.New("Static: Gin: 5 non-successful responses [2xx: 99995, 3xx: 0, 4xx: 5, 5xx: 0, others: 0]")},
					{Env: spec.Env{Name: "Koa", Language: "Javascript", NotSupported: true, Link: "https://github.com/koajs/koa"}, Status: runner.StatusSkipped},
				},
			},
			{
				Test: spec.Test{
					Name:        "REST",
					Description: "Fires 100000 requests with a 1MB JSON body.",
				},
				Envs: []runner.EnvReport{
					{Env: env("Iris", "Go", "kataras/iris"), Status: runner.StatusOK, Result: okResult(238954.12, 520.33, 100000)},
					{Env: env("Fastify", "Javascript", "fastify/fastify"), Status: runner.StatusOK, Result: okResult(160123.99, 780.45, 100000)},
				},
			},
		},
	}
}

func TestWriteGolden(t *testing.T) {
	dir := t.TempDir()
	if err := Write(sampleData(), Options{Dir: dir, Readme: true}); err != nil {
		t.Fatal(err)
	}

	goldenDir := filepath.Join("testdata", "golden")
	if *update {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	for _, name := range goldenFiles {
		got, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Errorf("missing artifact %s: %v", name, err)
			continue
		}

		goldenPath := filepath.Join(goldenDir, name+".golden")
		if *update {
			if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}

		want, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Errorf("missing golden for %s (run: go test ./internal/report -update): %v", name, err)
			continue
		}

		if !bytes.Equal(normalize(got), normalize(want)) {
			t.Errorf("%s differs from golden (run with -update after intentional changes)\n--- got ---\n%s", name, got)
		}
	}
}

// normalize strips CR so goldens compare identically across git eol settings.
func normalize(b []byte) []byte { return bytes.ReplaceAll(b, []byte("\r"), nil) }

func TestWriteTwiceIsIdentical(t *testing.T) {
	dir := t.TempDir()
	d := sampleData()

	if err := Write(d, Options{Dir: dir, Readme: true}); err != nil {
		t.Fatal(err)
	}

	first := make(map[string][]byte, len(goldenFiles))
	for _, name := range goldenFiles {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		first[name] = data
	}

	// Regression: the old writers opened files in append mode, so a
	// second run into the same directory doubled every report.
	if err := Write(d, Options{Dir: dir, Readme: true}); err != nil {
		t.Fatal(err)
	}

	for _, name := range goldenFiles {
		second, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(first[name], second) {
			t.Errorf("%s changed after a second Write into the same directory", name)
		}
	}
}

func TestMarkdownTableShape(t *testing.T) {
	dir := t.TempDir()
	if err := Write(sampleData(), Options{Dir: dir}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "RESULTS.md"))
	if err != nil {
		t.Fatal(err)
	}

	// Regression: skipped envs used to render 8 cells into the 6-column
	// results table. Every row of a results table must have exactly
	// 7 pipes (6 columns).
	inTable := false
	for line := range strings.Lines(string(content)) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "| Name | Language |") {
			inTable = true
		}
		if !inTable || !strings.HasPrefix(line, "|") {
			inTable = inTable && line != ""
			continue
		}
		if got := strings.Count(line, "|"); got != 7 {
			t.Errorf("table row has %d pipes, want 7: %s", got, line)
		}
		if line == "" {
			inTable = false
		}
	}

	// The skipped and failed envs render as "-" cells but keep their name.
	for _, want := range []string{
		"[Koa](https://github.com/koajs/koa) | Javascript | - | - | - | - |",
		"[Gin](https://github.com/gin-gonic/gin) | Go | - | - | - | - |",
		"> ⚠️ **Gin**:",
		"| Iris (next) | Go |", // no repo -> plain name, no link.
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("RESULTS.md does not contain %q", want)
		}
	}
}

func TestResultsJSONRoundTrip(t *testing.T) {
	dir := t.TempDir()
	d := sampleData()

	if err := Write(d, Options{Dir: dir, Readme: true}); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadResultsFile(filepath.Join(dir, ResultsJSONFile))
	if err != nil {
		t.Fatal(err)
	}

	// Re-rendering from the loaded document (what `export` does) must
	// reproduce every artifact byte for byte.
	exportDir := t.TempDir()
	if err := Write(loaded, Options{Dir: exportDir, Readme: true}); err != nil {
		t.Fatal(err)
	}

	for _, name := range goldenFiles {
		want, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(exportDir, name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(want, got) {
			t.Errorf("%s differs after a results.json round trip", name)
		}
	}
}

func TestLoadResultsFileRejectsUnknownSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.json")
	if err := os.WriteFile(path, []byte(`{"schemaVersion": 999}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadResultsFile(path); err == nil || !strings.Contains(err.Error(), "schemaVersion") {
		t.Fatalf("error = %v, want a schema version error", err)
	}
}

func TestChartSVGBasics(t *testing.T) {
	dir := t.TempDir()
	if err := Write(sampleData(), Options{Dir: dir}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "static.svg"))
	if err != nil {
		t.Fatal(err)
	}
	svg := string(content)

	for _, want := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		"prefers-color-scheme: dark", // theme-aware.
		">Iris<",                     // names present...
		">284,059<",                  // ...and values labeled at the tips.
		">Iris (next)<",
		"ASP.NET Core", // escaped safely.
		`class="s0"`,   // Go's color slot.
		`class="s1"`,   // Javascript's color slot.
		`class="s2"`,   // C#'s color slot.
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("static.svg does not contain %q", want)
		}
	}

	// Failed and skipped envs never chart.
	for _, banned := range []string{">Gin<", ">Koa<"} {
		if strings.Contains(svg, banned) {
			t.Errorf("static.svg contains %q, which did not run successfully", banned)
		}
	}

	// The latency chart ranks ascending: Iris (438.94us) stays first.
	latency, err := os.ReadFile(filepath.Join(dir, "static_latency.svg"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(latency), "438.94us") {
		t.Error("static_latency.svg does not label the best latency")
	}
	if irisIdx, nextIdx := strings.Index(string(latency), ">Iris<"), strings.Index(string(latency), ">Iris (next)<"); irisIdx < 0 || nextIdx < 0 || irisIdx > nextIdx {
		t.Error("static_latency.svg is not ranked by ascending latency")
	}
}

func TestNoChartWithoutSuccessfulEnvs(t *testing.T) {
	d := Data{
		GeneratedAt: time.Date(2026, 8, 12, 14, 30, 0, 0, time.UTC),
		Tests: []runner.TestReport{{
			Test: spec.Test{Name: "Static"},
			Envs: []runner.EnvReport{
				{Env: env("Gin", "Go", "gin-gonic/gin"), Status: runner.StatusFailed, Err: errors.New("boom")},
			},
		}},
	}

	dir := t.TempDir()
	if err := Write(d, Options{Dir: dir}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, "static.svg")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("static.svg exists for a test with zero successful envs (stat err: %v)", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "RESULTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "static.svg") {
		t.Error("RESULTS.md links a chart that was not generated")
	}
}
