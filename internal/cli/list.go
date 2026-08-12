package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/kataras/server-benchmarks/internal/spec"
)

// listCmd implements `server-benchmarks list`: print the tests and
// environments of a spec file without running anything.
func listCmd(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("list", stderr)
	specFile := fs.String("i", "./tests.yml", "path of the YAML spec file describing the tests")

	var selectors stringSlice
	fs.Var(&selectors, "t", "list only this test or test.env (repeatable)")

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

	w := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
	for i, t := range tests {
		if i > 0 {
			fmt.Fprintln(w)
		}

		method := t.Method
		if method == "" {
			method = "GET"
		}

		load := fmt.Sprintf("%d requests", t.NumberOfRequests)
		if t.NumberOfRequests == 0 {
			load = fmt.Sprintf("%s duration", t.Duration)
		}

		fmt.Fprintf(w, "%s: %s %s (%s, %d connections)\n", t.Name, method, t.URL, load, t.NumberOfConnections)

		for _, env := range t.Envs {
			marker := ""
			switch {
			case env.NotSupported:
				marker = "not supported"
			case env.NotYetImplemented:
				marker = "not yet implemented"
			}

			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", env.Name, env.Language, env.Dir, marker)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFailure
	}

	return ExitOK
}
