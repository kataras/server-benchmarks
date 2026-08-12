// Command server-benchmarks stress-tests web servers (frameworks) with
// the bombardier load generator and renders Markdown, CSV, JSON and SVG
// chart reports.
//
// Usage:
//
//	server-benchmarks [run] [-i tests.yml] [-o dir] [-t test[.env]]...
//	server-benchmarks list [-i tests.yml]
//	server-benchmarks export [-i results.json] [-o dir]
//	server-benchmarks version
//
// Run `server-benchmarks help` for details.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/server-benchmarks/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Once the first signal cancels ctx, restore the default handling so
	// a second Ctrl-C terminates immediately instead of waiting for the
	// graceful teardown.
	go func() {
		<-ctx.Done()
		stop()
	}()

	os.Exit(cli.Main(ctx, os.Args[1:], os.Stdout, os.Stderr))
}
