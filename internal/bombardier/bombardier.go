// Package bombardier drives the external bombardier load generator
// (github.com/codesenberg/bombardier): it builds its command-line
// arguments, executes it and parses its JSON output.
package bombardier

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// DefaultBin is the executable name used when Client.Bin is empty.
const DefaultBin = "bombardier"

// Stats holds a latency or requests-per-second distribution as reported
// by bombardier.
type Stats struct {
	Mean   float64 `json:"mean"`
	Stddev float64 `json:"stddev"`
	Max    float64 `json:"max"`
}

// ResultError is a transport-level error (timeouts, connection resets...)
// bombardier encountered while firing requests.
type ResultError struct {
	Description string `json:"description"`
	Count       uint64 `json:"count"`
}

// Result holds the measurements of a single bombardier run.
// Field names and units mirror bombardier's own JSON output:
// latency values are microseconds, TimeTakenSeconds is seconds.
type Result struct {
	BytesRead        int64   `json:"bytesRead"`
	BytesWritten     int64   `json:"bytesWritten"`
	TimeTakenSeconds float64 `json:"timeTakenSeconds"`

	Req1XX uint64 `json:"req1xx"`
	Req2XX uint64 `json:"req2xx"`
	Req3XX uint64 `json:"req3xx"`
	Req4XX uint64 `json:"req4xx"`
	Req5XX uint64 `json:"req5xx"`
	Others uint64 `json:"others"`

	Errors []ResultError `json:"errors,omitempty"`

	Latency           Stats `json:"latency"`
	RequestsPerSecond Stats `json:"rps"`
}

// Throughput returns the total throughput (read + write) in bytes per second.
func (r Result) Throughput() float64 {
	if r.TimeTakenSeconds == 0 {
		return 0
	}

	return float64(r.BytesRead+r.BytesWritten) / r.TimeTakenSeconds
}

// ErrorDescriptions returns the descriptions of all transport errors.
func (r Result) ErrorDescriptions() []string {
	if len(r.Errors) == 0 {
		return nil
	}

	descriptions := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		descriptions[i] = e.Description
	}

	return descriptions
}

// Client executes the bombardier binary.
// The zero value is ready to use and invokes DefaultBin from PATH.
type Client struct {
	// Bin is the bombardier executable: a name looked up in PATH or a
	// full path. Empty means DefaultBin.
	Bin string
}

// bin returns the executable to invoke.
func (c Client) bin() string {
	if c.Bin != "" {
		return c.Bin
	}

	return DefaultBin
}

// LookPath verifies that the bombardier executable can be found,
// returning its resolved path.
func (c Client) LookPath() (string, error) {
	path, err := exec.LookPath(c.bin())
	if err != nil {
		return "", fmt.Errorf("bombardier executable %q not found (install: go install github.com/codesenberg/bombardier@latest): %w", c.bin(), err)
	}

	return path, nil
}

// String returns the command line c would execute for the given arguments,
// for logging purposes.
func (c Client) String(args []string) string {
	return c.bin() + " " + strings.Join(args, " ")
}

// Run executes bombardier with the given arguments and parses its JSON
// output. The process is killed when ctx is canceled.
func (c Client) Run(ctx context.Context, args []string) (Result, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, c.bin(), args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = strings.TrimSpace(stdout.String())
		}

		return Result{}, fmt.Errorf("bombardier: %w: %s", err, msg)
	}

	return Parse(stdout.Bytes())
}
