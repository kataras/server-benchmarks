package cli

import (
	"fmt"
	"io"

	"github.com/kataras/server-benchmarks/internal/report"
)

// exportCmd implements `server-benchmarks export`: re-render every report
// artifact (Markdown, CSV, SVG charts, optionally README) from a saved
// results.json, without re-running any benchmark.
func exportCmd(args []string, stderr io.Writer) int {
	fs := newFlagSet("export", stderr)
	input := fs.String("i", "./results.json", "path of a results.json produced by a previous run")
	outDir := fs.String("o", "./", "directory for the re-rendered reports")
	readme := fs.Bool("readme", false, "also generate a README.md next to RESULTS.md")

	if err := fs.Parse(args); err != nil {
		return parseCode(err)
	}

	data, err := report.LoadResultsFile(*input)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitUsage
	}

	if err := report.Write(data, report.Options{Dir: *outDir, Readme: *readme}); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFailure
	}

	return ExitOK
}
