// Command fakebombardier mimics `bombardier --format=json --print=result`
// for the runner tests, without generating any load. It performs a single
// real GET against the target URL (the last argument) to prove the server
// was actually reachable, then prints a canned JSON result.
//
// Environment variables:
//
//	FAKEBOMB_MODE     "" (ok) | http4xx | mismatch | transport-error |
//	                  noresult | garbage
//	FAKEBOMB_COUNTER  path to a counter file; each invocation increments
//	                  it and lowers the reported RPS, giving successive
//	                  environments distinct, decreasing numbers
//	FAKEBOMB_FAIL_FIRST  when "1", the first invocation (counter 0)
//	                  reports http4xx and later ones succeed
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type stats struct {
	Mean   float64 `json:"mean"`
	Stddev float64 `json:"stddev"`
	Max    float64 `json:"max"`
}

type result struct {
	BytesRead        int64         `json:"bytesRead"`
	BytesWritten     int64         `json:"bytesWritten"`
	TimeTakenSeconds float64       `json:"timeTakenSeconds"`
	Req1XX           uint64        `json:"req1xx"`
	Req2XX           uint64        `json:"req2xx"`
	Req3XX           uint64        `json:"req3xx"`
	Req4XX           uint64        `json:"req4xx"`
	Req5XX           uint64        `json:"req5xx"`
	Others           uint64        `json:"others"`
	Errors           []resultError `json:"errors"`
	Latency          stats         `json:"latency"`
	RPS              stats         `json:"rps"`
}

type resultError struct {
	Description string `json:"description"`
	Count       uint64 `json:"count"`
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "fakebombardier: missing target URL")
		os.Exit(2)
	}

	url := args[len(args)-1]

	// The number of requests the runner asked for (-n), if any.
	var requested uint64 = 1000
	for i, a := range args {
		if a == "-n" && i+1 < len(args) {
			if v, err := strconv.ParseUint(args[i+1], 10, 64); err == nil {
				requested = v
			}
		}
	}

	count := bumpCounter()

	mode := os.Getenv("FAKEBOMB_MODE")
	if os.Getenv("FAKEBOMB_FAIL_FIRST") == "1" && count == 0 {
		mode = "http4xx"
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fakebombardier: target unreachable: %v\n", err)
		os.Exit(1)
	}
	resp.Body.Close()

	switch mode {
	case "noresult":
		fmt.Println(`{"spec": {}}`)
		return
	case "garbage":
		fmt.Println(`this is not json`)
		return
	}

	res := result{
		BytesRead:        int64(requested) * 150,
		BytesWritten:     int64(requested) * 40,
		TimeTakenSeconds: 0.5,
		Req2XX:           requested,
		Latency:          stats{Mean: 500 + float64(count)*100, Stddev: 50, Max: 2000},
		RPS:              stats{Mean: 200000 - float64(count)*50000, Stddev: 1000, Max: 250000},
		Errors:           []resultError{},
	}

	switch mode {
	case "http4xx":
		res.Req4XX = 5
		res.Req2XX = requested - 5
	case "mismatch":
		res.Req2XX = requested - 1
	case "transport-error":
		res.Errors = append(res.Errors, resultError{Description: "the connection was forcibly closed", Count: 3})
	}

	out, err := json.Marshal(map[string]any{"result": res})
	if err != nil {
		fmt.Fprintln(os.Stderr, "fakebombardier:", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}

// bumpCounter returns the current invocation index and persists the next
// one, when FAKEBOMB_COUNTER is set. Invocations are sequential by design
// (the runner never runs two benchmarks at once), so a plain file is safe.
func bumpCounter() int {
	path := os.Getenv("FAKEBOMB_COUNTER")
	if path == "" {
		return 0
	}

	count := 0
	if data, err := os.ReadFile(path); err == nil {
		if v, err := strconv.Atoi(string(data)); err == nil {
			count = v
		}
	}

	_ = os.WriteFile(path, []byte(strconv.Itoa(count+1)), 0o644)
	return count
}
