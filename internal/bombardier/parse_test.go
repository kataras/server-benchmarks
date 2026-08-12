package bombardier

import (
	"strings"
	"testing"
)

// sampleOutput mirrors real bombardier --format=json --print=result output.
const sampleOutput = `{
  "spec": {"numberOfConnections": 125, "testType": "number-of-requests", "numberOfRequests": 100000, "method": "GET", "url": "http://localhost:5000"},
  "result": {
    "bytesRead": 15200000, "bytesWritten": 4100000, "timeTakenSeconds": 0.352,
    "req1xx": 0, "req2xx": 100000, "req3xx": 0, "req4xx": 0, "req5xx": 0, "others": 0,
    "errors": [],
    "latency": {"mean": 438.94, "stddev": 197.67, "max": 12508.0},
    "rps": {"mean": 284059.44, "stddev": 32294.13, "max": 344467.44}
  }
}`

func TestParse(t *testing.T) {
	res, err := Parse([]byte(sampleOutput))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := res.Req2XX, uint64(100000); got != want {
		t.Errorf("Req2XX = %d, want %d", got, want)
	}
	if got, want := res.RequestsPerSecond.Mean, 284059.44; got != want {
		t.Errorf("RPS mean = %f, want %f", got, want)
	}
	if got, want := res.Latency.Max, 12508.0; got != want {
		t.Errorf("latency max = %f, want %f", got, want)
	}
	if len(res.Errors) != 0 {
		t.Errorf("errors = %v, want none", res.Errors)
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name    string
		data    string
		wantSub string
	}{
		{"empty output", ``, "parse bombardier output"},
		{"malformed json", `{"result": nope}`, "parse bombardier output"},
		{"missing result object", `{"spec": {}}`, `missing "result"`},
		{"null result", `{"result": null}`, `missing "result"`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.data))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Fatalf("error = %q, want it to contain %q", err, tc.wantSub)
			}
		})
	}
}

func TestThroughput(t *testing.T) {
	r := Result{BytesRead: 1000, BytesWritten: 500, TimeTakenSeconds: 2}
	if got, want := r.Throughput(), 750.0; got != want {
		t.Errorf("Throughput = %f, want %f", got, want)
	}

	var zero Result
	if got := zero.Throughput(); got != 0 {
		t.Errorf("zero-value Throughput = %f, want 0 (no division by zero)", got)
	}
}
