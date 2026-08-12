// Command fakeserver is a minimal HTTP server used by the runner tests as
// a stand-in for a framework's benchmark application.
//
// Modes:
//
//	serve  listen on -addr and answer 200 "Index" (default)
//	stall  never listen, sleep forever (tests the readiness timeout)
//	crash  print to stderr and exit non-zero (tests early-exit detection)
//	spawn  re-exec itself as a serving child and block, emulating
//	       wrappers like `go run .` whose real server is a grandchild
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:0", "listen address")
	mode := flag.String("mode", "serve", "serve|stall|crash|spawn")
	flag.Parse()

	switch *mode {
	case "crash":
		fmt.Fprintln(os.Stderr, "boom: intentional crash")
		os.Exit(3)
	case "stall":
		for {
			time.Sleep(time.Hour) // never listen; a bare select{} would trip the deadlock detector.
		}
	case "spawn":
		child := exec.Command(os.Args[0], "-addr", *addr, "-mode", "serve")
		child.Stdout = os.Stdout
		child.Stderr = os.Stderr
		if err := child.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "spawn:", err)
			os.Exit(1)
		}
		_ = child.Wait()
	default:
		serve(*addr)
	}
}

func serve(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Index")
	})

	if err := http.Serve(listener, mux); err != nil {
		fmt.Fprintln(os.Stderr, "serve:", err)
		os.Exit(1)
	}
}
