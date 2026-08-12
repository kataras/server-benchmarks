// Package report renders benchmark results into human- and
// machine-readable artifacts: a Markdown report, per-test CSV files,
// SVG bar charts, a results.json document and, optionally, a standalone
// README.
//
// Every file is written whole (truncate-then-write), so re-running a
// benchmark into the same directory always produces a fresh report
// instead of appending to a stale one.
package report

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/runner"
	"github.com/kataras/server-benchmarks/internal/sysinfo"
)

// TimeLayout is the human-readable timestamp format used in reports.
const TimeLayout = "Jan 2, 2006 at 3:04pm (UTC)"

// Filenames of the rendered artifacts inside Options.Dir. Per-test CSV
// and SVG files derive from the test name: static.csv, static_latency.csv,
// static.svg and static_latency.svg for a test named "Static".
const (
	ResultsMarkdownFile = "RESULTS.md"
	ReadmeFile          = "README.md"
	ResultsJSONFile     = "results.json"
)

// Data is everything a report is rendered from.
type Data struct {
	GeneratedAt time.Time
	System      sysinfo.Info
	Tests       []runner.TestReport
}

// Options configures Write.
type Options struct {
	// Dir is the output directory; it is created when missing.
	// Empty means the current directory.
	Dir string
	// Readme also renders a publishable README.md alongside RESULTS.md.
	Readme bool
}

// Write renders all artifacts for d into opts.Dir.
func Write(d Data, opts Options) error {
	if opts.Dir == "" {
		opts.Dir = "."
	}

	if err := os.MkdirAll(opts.Dir, 0o755); err != nil {
		return fmt.Errorf("output directory: %w", err)
	}

	v := buildView(d)
	var errs []error

	write := func(name string, data []byte) {
		if err := os.WriteFile(filepath.Join(opts.Dir, name), data, 0o644); err != nil {
			errs = append(errs, fmt.Errorf("write %s: %w", name, err))
		}
	}

	if data, err := executeTemplate("results.md.tmpl", v); err != nil {
		errs = append(errs, err)
	} else {
		write(ResultsMarkdownFile, data)
	}

	if opts.Readme {
		if data, err := executeTemplate("readme.md.tmpl", v); err != nil {
			errs = append(errs, err)
		} else {
			write(ReadmeFile, data)
		}
	}

	for _, tv := range v.Tests {
		if data, err := renderCSV("Reqs/sec", tv, func(r *bombardier.Result) string {
			return fmt.Sprintf("%.0f", r.RequestsPerSecond.Mean)
		}); err != nil {
			errs = append(errs, err)
		} else {
			write(tv.Slug+".csv", data)
		}

		if data, err := renderCSV("Latency", tv, func(r *bombardier.Result) string {
			return fmt.Sprintf("%.2f", r.Latency.Mean)
		}); err != nil {
			errs = append(errs, err)
		} else {
			write(tv.Slug+"_latency.csv", data)
		}

		if tv.HasChart {
			write(tv.ChartFile, tv.chartSVG)
			write(tv.LatencyChartFile, tv.latencySVG)
		}
	}

	if data, err := marshalResults(d); err != nil {
		errs = append(errs, err)
	} else {
		write(ResultsJSONFile, data)
	}

	return errors.Join(errs...)
}

// view is the fully pre-formatted template model: templates only print,
// they never compute.
type view struct {
	Datetime string
	System   sysinfo.Info
	Tests    []testView
}

// testView is one test's section of the report.
type testView struct {
	Name        string
	Slug        string
	Description string

	HasChart         bool
	ChartFile        string
	LatencyChartFile string
	chartSVG         []byte
	latencySVG       []byte

	Rows     []rowView
	Failures []failureView

	// results carries the successful envs for CSV rendering, in report
	// order.
	results []runner.EnvReport
}

// rowView is one Markdown table row. Non-successful environments render
// "-" cells, keeping the table shape intact.
type rowView struct {
	Name     string
	Link     string
	Language string

	ReqsPerSec string
	Latency    string
	Throughput string
	TimeTaken  string
}

// failureView is a short failure note rendered below the table.
type failureView struct {
	Name  string
	Error string
}

// buildView pre-formats d for the templates and renders the charts.
func buildView(d Data) view {
	slots := languageSlots(d.Tests)

	v := view{
		Datetime: d.GeneratedAt.UTC().Format(TimeLayout),
		System:   d.System,
	}

	for _, tr := range d.Tests {
		tv := testView{
			Name:        tr.Test.Name,
			Slug:        slug(tr.Test.Name),
			Description: tr.Test.Description,
		}

		for _, er := range tr.Envs {
			row := rowView{
				Name:     er.Env.Name,
				Link:     er.Env.Link,
				Language: er.Env.Language,

				ReqsPerSec: "-",
				Latency:    "-",
				Throughput: "-",
				TimeTaken:  "-",
			}

			switch er.Status {
			case runner.StatusOK:
				r := er.Result
				row.ReqsPerSec = fmt.Sprintf("%.0f", r.RequestsPerSecond.Mean)
				row.Latency = bombardier.FormatTimeUs(r.Latency.Mean)
				row.Throughput = bombardier.FormatBinary(r.Throughput())
				row.TimeTaken = fmt.Sprintf("%.2fs", r.TimeTakenSeconds)
				tv.results = append(tv.results, er)
			case runner.StatusFailed:
				tv.Failures = append(tv.Failures, failureView{Name: er.Env.Name, Error: errorString(er.Err)})
			}

			tv.Rows = append(tv.Rows, row)
		}

		if len(tv.results) > 0 {
			tv.HasChart = true
			tv.ChartFile = tv.Slug + ".svg"
			tv.LatencyChartFile = tv.Slug + "_latency.svg"
			tv.chartSVG = renderRPSChart(tr.Test.Name, tv.results, slots)
			tv.latencySVG = renderLatencyChart(tr.Test.Name, tv.results, slots)
		}

		v.Tests = append(v.Tests, tv)
	}

	return v
}

// executeTemplate renders a named template into memory, so files are only
// ever written whole.
func executeTemplate(name string, v view) ([]byte, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, v); err != nil {
		return nil, fmt.Errorf("render %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

// renderCSV renders a test's two-column CSV (Name plus one metric) for
// its successful environments.
func renderCSV(header string, tv testView, value func(*bombardier.Result) string) ([]byte, error) {
	var buf bytes.Buffer

	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"Name", header})
	for _, er := range tv.results {
		_ = w.Write([]string{er.Env.Name, value(er.Result)})
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("render %s.csv: %w", tv.Slug, err)
	}

	return buf.Bytes(), nil
}

// slug converts a test name into a safe lowercase file name fragment.
func slug(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

// errorString formats an error for display, tolerating nil.
func errorString(err error) string {
	if err == nil {
		return "unknown error"
	}

	return err.Error()
}
