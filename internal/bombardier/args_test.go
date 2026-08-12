package bombardier

import (
	"reflect"
	"testing"
	"time"

	"github.com/kataras/server-benchmarks/internal/spec"
)

func TestArgs(t *testing.T) {
	cases := []struct {
		name     string
		test     spec.Test
		bodyPath string
		want     []string
	}{
		{
			name: "request count run",
			test: spec.Test{
				NumberOfConnections: 125,
				NumberOfRequests:    100000,
				Timeout:             spec.Duration(15 * time.Second),
				Method:              "GET",
				URL:                 "http://localhost:5000",
			},
			want: []string{
				"-c", "125", "-n", "100000", "-t", "15s", "-m", "GET",
				"--format=json", "--print=result", "http://localhost:5000",
			},
		},
		{
			name: "duration run",
			test: spec.Test{
				NumberOfConnections: 50,
				Duration:            spec.Duration(5 * time.Second),
				Timeout:             spec.Duration(2 * time.Second),
				URL:                 "http://localhost:5000",
			},
			want: []string{
				"-c", "50", "-d", "5s", "-t", "2s",
				"--format=json", "--print=result", "http://localhost:5000",
			},
		},
		{
			name: "post with downloaded json body infers content type",
			test: spec.Test{
				NumberOfConnections: 125,
				NumberOfRequests:    1000,
				Timeout:             spec.Duration(15 * time.Second),
				Method:              "POST",
				BodyFile:            "https://example.com/1MB.json",
				URL:                 "http://localhost:5000/42",
			},
			bodyPath: "C:\\tmp\\body.json",
			want: []string{
				"-c", "125", "-n", "1000", "-t", "15s", "-m", "POST",
				"-f", "C:\\tmp\\body.json",
				"-H", "Content-Type: application/json",
				"--format=json", "--print=result", "http://localhost:5000/42",
			},
		},
		{
			name: "declared content type wins over inference",
			test: spec.Test{
				NumberOfConnections: 1,
				NumberOfRequests:    1,
				Timeout:             spec.Duration(time.Second),
				Method:              "POST",
				BodyFile:            "body.json",
				Headers:             map[string]string{"Content-Type": "application/vnd.custom"},
				URL:                 "http://localhost:5000",
			},
			want: []string{
				"-c", "1", "-n", "1", "-t", "1s", "-m", "POST",
				"-f", "body.json",
				"-H", "Content-Type: application/vnd.custom",
				"--format=json", "--print=result", "http://localhost:5000",
			},
		},
		{
			name: "get with body does not infer content type",
			test: spec.Test{
				NumberOfConnections: 1,
				NumberOfRequests:    1,
				Timeout:             spec.Duration(time.Second),
				Method:              "GET",
				BodyFile:            "body.json",
				URL:                 "http://localhost:5000",
			},
			want: []string{
				"-c", "1", "-n", "1", "-t", "1s", "-m", "GET",
				"-f", "body.json",
				"--format=json", "--print=result", "http://localhost:5000",
			},
		},
		{
			name: "headers are sorted deterministically",
			test: spec.Test{
				NumberOfConnections: 1,
				NumberOfRequests:    1,
				Timeout:             spec.Duration(time.Second),
				Headers: map[string]string{
					"X-Zeta":  "1",
					"Accept":  "application/json",
					"X-Alpha": "2",
				},
				URL: "http://localhost:5000",
			},
			want: []string{
				"-c", "1", "-n", "1", "-t", "1s",
				"-H", "Accept: application/json",
				"-H", "X-Alpha: 2",
				"-H", "X-Zeta: 1",
				"--format=json", "--print=result", "http://localhost:5000",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := tc.test // shallow copy to detect mutation of scalar fields.

			got := Args(tc.test, tc.bodyPath)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Args() =\n  %q\nwant:\n  %q", got, tc.want)
			}

			// Run twice: result must be identical (map iteration must not leak).
			if second := Args(tc.test, tc.bodyPath); !reflect.DeepEqual(got, second) {
				t.Errorf("Args() is not deterministic:\n  %q\nvs\n  %q", got, second)
			}

			if tc.test.Method != before.Method || len(tc.test.Headers) != len(before.Headers) {
				t.Error("Args() mutated its input")
			}
		})
	}
}
