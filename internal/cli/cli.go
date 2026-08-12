// Package cli implements the server-benchmarks command-line interface:
// subcommand dispatch, flag parsing, logging setup and process exit codes.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// Exit codes returned by Main.
const (
	// ExitOK means the command completed successfully.
	ExitOK = 0
	// ExitFailure means a benchmark or rendering failure.
	ExitFailure = 1
	// ExitUsage means invalid flags, an unknown command or a bad spec.
	ExitUsage = 2
	// ExitInterrupted means the run was canceled by SIGINT/SIGTERM.
	ExitInterrupted = 130
)

// Main runs the command line interface and returns the process exit code.
// args are the program arguments without the program name. With no
// subcommand (or with flags only), the run subcommand is executed — so a
// plain `server-benchmarks` benchmarks everything, like it always did.
func Main(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	var command string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		command, args = args[0], args[1:]
	}

	switch command {
	case "", "run":
		return runCmd(ctx, args, stderr)
	case "list":
		return listCmd(args, stdout, stderr)
	case "export":
		return exportCmd(args, stderr)
	case "version":
		return versionCmd(stdout)
	case "help":
		printUsage(stdout)
		return ExitOK
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", command)
		printUsage(stderr)
		return ExitUsage
	}
}

// printUsage writes the global help text.
func printUsage(w io.Writer) {
	fmt.Fprint(w, `server-benchmarks stress-tests web servers with the bombardier load
generator and renders Markdown, CSV, JSON and SVG chart reports.

Usage:

  server-benchmarks [run] [flags]    benchmark and write the reports (default)
  server-benchmarks list [flags]     print the tests of a spec file
  server-benchmarks export [flags]   re-render reports from a results.json
  server-benchmarks version          print version information
  server-benchmarks help             print this text

Run 'server-benchmarks <command> -h' to see a command's flags.
`)
}

// newFlagSet builds a subcommand flag set that reports to stderr and
// never exits the process on its own.
func newFlagSet(name string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	return fs
}

// parseCode translates a flag parse error into an exit code, treating
// -h/-help as success.
func parseCode(err error) int {
	if errors.Is(err, flag.ErrHelp) {
		return ExitOK
	}

	return ExitUsage
}

// newLogger builds the CLI logger; verbose raises it to debug level.
func newLogger(w io.Writer, verbose bool) *slog.Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: level}))
}
