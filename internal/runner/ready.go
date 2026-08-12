package runner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"
)

// errServerExited reports that the server process died before it started
// listening.
var errServerExited = errors.New("server exited before becoming ready")

// dialProbe is how long a single readiness dial may take.
const dialProbe = 250 * time.Millisecond

// targetAddr extracts the host:port to probe from a test URL, deriving
// the port from the scheme when absent.
func targetAddr(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid test URL %q: %w", rawURL, err)
	}

	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("test URL %q has no host", rawURL)
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		default:
			port = "80"
		}
	}

	return net.JoinHostPort(host, port), nil
}

// portBusy reports whether something is currently accepting connections
// on addr.
func portBusy(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, dialProbe)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

// waitReady blocks until a TCP connection to addr succeeds, the server
// process exits, ctx is canceled or the timeout elapses.
func waitReady(ctx context.Context, addr string, exited <-chan struct{}, timeout time.Duration) error {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	ticker := time.NewTicker(dialProbe)
	defer ticker.Stop()

	for {
		conn, err := net.DialTimeout("tcp", addr, dialProbe)
		if err == nil {
			conn.Close()
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-exited:
			return errServerExited
		case <-deadline.C:
			return fmt.Errorf("server at %s did not become ready within %s", addr, timeout)
		case <-ticker.C:
		}
	}
}

// waitPortFree blocks until nothing accepts connections on addr anymore,
// so the next environment's server can bind it.
func waitPortFree(ctx context.Context, addr string, timeout time.Duration) error {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	ticker := time.NewTicker(dialProbe)
	defer ticker.Stop()

	for {
		if !portBusy(addr) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("%s is still accepting connections after %s", addr, timeout)
		case <-ticker.C:
		}
	}
}
