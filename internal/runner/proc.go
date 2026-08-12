package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode"
)

// splitCommand splits a command line into the executable name and its
// arguments. Double and single quotes group words (useful for paths with
// spaces, e.g. "C:\Program Files\dotnet\dotnet" run); there is no further
// escape processing.
func splitCommand(line string) (name string, args []string, err error) {
	var (
		fields  []string
		current strings.Builder
		quote   rune // active quote character, 0 when outside quotes.
		inField bool
	)

	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote = r
			inField = true
		case unicode.IsSpace(r):
			if inField {
				fields = append(fields, current.String())
				current.Reset()
				inField = false
			}
		default:
			current.WriteRune(r)
			inField = true
		}
	}

	if quote != 0 {
		return "", nil, fmt.Errorf("unterminated quote in command %q", line)
	}

	if inField {
		fields = append(fields, current.String())
	}

	if len(fields) == 0 {
		return "", nil, errors.New("empty command")
	}

	return fields[0], fields[1:], nil
}

// execLines splits a (possibly multiline) Exec value into trimmed,
// non-empty command lines.
func execLines(s string) []string {
	var lines []string
	for line := range strings.Lines(s) {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

// processPid returns the cmd's PID, or 0 when it has not started.
func processPid(cmd *exec.Cmd) int {
	if cmd.Process == nil {
		return 0
	}

	return cmd.Process.Pid
}

// newCommand builds an exec.Cmd whose whole process tree is terminated
// when ctx is canceled or, past WaitDelay, when its output pipes are
// still open.
func newCommand(ctx context.Context, dir, line string, grace time.Duration, out io.Writer) (*exec.Cmd, error) {
	name, args, err := splitCommand(line)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out
	setSysProcAttr(cmd)
	// The default Cancel kills only the direct child; commands like
	// `go run .` re-spawn the real server as a grandchild, so the whole
	// tree must go.
	cmd.Cancel = func() error { return terminateTree(cmd) }
	cmd.WaitDelay = grace

	return cmd, nil
}

// runCommand runs a setup command (e.g. `npm install`) to completion.
func runCommand(ctx context.Context, dir, line string, grace time.Duration, out io.Writer) error {
	cmd, err := newCommand(ctx, dir, line, grace, out)
	if err != nil {
		return err
	}

	return cmd.Run()
}

// process is a started, long-running server process.
type process struct {
	cmd    *exec.Cmd
	grace  time.Duration
	exited chan struct{}
	// waitErr is written once before exited is closed and must only be
	// read after <-exited.
	waitErr error
}

// startProcess launches line as the blocking server process inside dir.
func startProcess(ctx context.Context, dir, line string, grace time.Duration, out io.Writer) (*process, error) {
	cmd, err := newCommand(ctx, dir, line, grace, out)
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	p := &process{
		cmd:    cmd,
		grace:  grace,
		exited: make(chan struct{}),
	}

	go func() {
		p.waitErr = cmd.Wait()
		close(p.exited)
	}()

	return p, nil
}

// stop terminates the server's process tree: it asks politely first
// (SIGTERM to the process group on unix, a forced tree kill on Windows,
// which knows no gentler tree-wide signal) and escalates to a hard kill
// after the grace period.
func (p *process) stop() error {
	select {
	case <-p.exited:
		return nil // already gone.
	default:
	}

	signalErr := signalTree(p.cmd)

	select {
	case <-p.exited:
		return nil
	case <-time.After(p.grace):
	}

	terminateErr := terminateTree(p.cmd)

	select {
	case <-p.exited:
		return nil
	case <-time.After(p.grace):
		return errors.Join(
			fmt.Errorf("server process %d is still running after kill", processPid(p.cmd)),
			signalErr, terminateErr,
		)
	}
}

// tail is a concurrency-safe writer that keeps only the last max bytes,
// used to attach recent process output to error messages.
type tail struct {
	mu  sync.Mutex
	max int
	buf []byte
}

func newTail(max int) *tail { return &tail{max: max} }

// Write implements io.Writer.
func (t *tail) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.buf = append(t.buf, p...)
	if over := len(t.buf) - t.max; over > 0 {
		t.buf = append(t.buf[:0], t.buf[over:]...)
	}

	return len(p), nil
}

// String returns the captured output, trimmed.
func (t *tail) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return strings.TrimSpace(string(t.buf))
}

// logWriter forwards process output lines to the logger at debug level.
type logWriter struct {
	log *slog.Logger
	tag string
}

// Write implements io.Writer.
func (w logWriter) Write(p []byte) (int, error) {
	for line := range strings.Lines(string(p)) {
		if trimmed := strings.TrimRight(line, "\r\n"); trimmed != "" {
			w.log.Debug("server output", "env", w.tag, "line", trimmed)
		}
	}

	return len(p), nil
}
