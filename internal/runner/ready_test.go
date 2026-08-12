package runner

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

func neverExited() <-chan struct{} { return make(chan struct{}) }

func TestWaitReadyImmediate(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	if err := waitReady(context.Background(), listener.Addr().String(), neverExited(), 5*time.Second); err != nil {
		t.Fatalf("waitReady against live listener: %v", err)
	}
}

func TestWaitReadyTimeout(t *testing.T) {
	// Reserve a port and close it so nothing listens there.
	addr := freeAddr(t)

	err := waitReady(context.Background(), addr, neverExited(), 600*time.Millisecond)
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	if !strings.Contains(err.Error(), "did not become ready") {
		t.Fatalf("error = %v, want a readiness timeout", err)
	}
}

func TestWaitReadyProcessExit(t *testing.T) {
	addr := freeAddr(t)

	exited := make(chan struct{})
	close(exited)

	err := waitReady(context.Background(), addr, exited, 5*time.Second)
	if !errors.Is(err, errServerExited) {
		t.Fatalf("error = %v, want errServerExited", err)
	}
}

func TestWaitReadyContextCancel(t *testing.T) {
	addr := freeAddr(t)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err := waitReady(ctx, addr, neverExited(), 10*time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
}

func TestWaitPortFree(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()

	// Busy port times out.
	if err := waitPortFree(context.Background(), addr, 600*time.Millisecond); err == nil {
		t.Fatal("expected an error while the listener is alive")
	}

	listener.Close()

	if err := waitPortFree(context.Background(), addr, 5*time.Second); err != nil {
		t.Fatalf("waitPortFree after close: %v", err)
	}
}

// freeAddr reserves a TCP port on the loopback interface and releases it,
// returning an address that (very likely) nothing listens on.
func freeAddr(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	addr := listener.Addr().String()
	listener.Close()
	return addr
}
