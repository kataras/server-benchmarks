// Package spec defines the benchmark suite specification: the tests to run
// and the environments (frameworks) each test runs against.
//
// A specification is loaded once from a YAML document via Load or LoadFile,
// which applies defaults, normalizes paths and validates the result. After
// loading, Test and Env values are treated as immutable: no other package
// mutates them, which keeps them safe to share across goroutines.
package spec

import (
	"fmt"
	"strings"
	"time"
)

// Defaults applied by Load when the corresponding field is not set.
const (
	// DefaultLanguage is the language assumed for an Env without one.
	DefaultLanguage = "Go"
	// DefaultConnections is the number of concurrent connections used
	// when a Test does not declare NumberOfConnections.
	DefaultConnections = 125
	// DefaultTimeout is the per-request timeout used when a Test does
	// not declare one.
	DefaultTimeout = 15 * time.Second
	// DefaultDuration is the test duration used when a Test declares
	// neither NumberOfRequests nor Duration.
	DefaultDuration = 5 * time.Second
)

// DefaultCodeDir is the directory, relative to the specification file,
// that hosts the per-test, per-framework server applications
// (e.g. _code/static/gin).
const DefaultCodeDir = "_code"

// Duration is a time.Duration that unmarshals from human-readable YAML
// scalars such as "5s", "150ms" or "1m30s".
type Duration time.Duration

// String returns the duration in time.Duration string form, e.g. "5s".
func (d Duration) String() string { return time.Duration(d).String() }

// Std returns the standard library representation of d.
func (d Duration) Std() time.Duration { return time.Duration(d) }

// UnmarshalYAML parses scalars like "5s" or "150ms" via time.ParseDuration.
func (d *Duration) UnmarshalYAML(b []byte) error {
	s := strings.TrimSpace(string(b))
	s = strings.Trim(s, `"'`)
	if s == "" {
		*d = 0
		return nil
	}

	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q (use forms like 5s, 150ms or 1m30s)", s)
	}

	*d = Duration(v)
	return nil
}

// Test describes a single benchmark group: one workload (method, URL,
// payload and load characteristics) executed against every environment
// listed in Envs.
type Test struct {
	// Name is the unique test name, e.g. "Static". Its lowercase form
	// doubles as the default code sub-directory and the output file slug.
	Name string `yaml:"Name"`
	// Description is rendered as Markdown in the generated reports.
	// It may reference the test itself as template data,
	// e.g. "Fires {{.NumberOfRequests}} requests".
	Description string `yaml:"Description"`

	// NumberOfConnections is the number of concurrent connections the
	// load generator keeps open. Defaults to DefaultConnections.
	NumberOfConnections uint64 `yaml:"NumberOfConnections"`
	// NumberOfRequests is the total number of requests to fire.
	// Mutually complementary with Duration: when both are zero,
	// Duration defaults to DefaultDuration.
	NumberOfRequests uint64 `yaml:"NumberOfRequests"`
	// Duration is how long to fire requests for, e.g. 5s.
	Duration Duration `yaml:"Duration"`
	// Timeout is the per-request timeout. Defaults to DefaultTimeout.
	Timeout Duration `yaml:"Timeout"`
	// Headers holds extra HTTP headers to send with every request.
	// For POST and PUT requests with a BodyFile, a Content-Type header
	// is inferred from the file extension when not declared here.
	Headers map[string]string `yaml:"Headers"`
	// Method is the HTTP method, e.g. GET or POST. Empty means the load
	// generator's default (GET).
	Method string `yaml:"Method"`
	// URL is the target endpoint, e.g. http://localhost:5000/hello/world.
	// Every environment's server is expected to listen on its host:port.
	URL string `yaml:"URL"`
	// BodyFile is a local file path or an http(s) URL whose content is
	// sent as the request body. URLs are downloaded once per test into a
	// temporary directory and removed afterwards. Relative paths are
	// resolved against the specification file's directory.
	BodyFile string `yaml:"BodyFile"`

	// Envs lists the environments (frameworks) this test runs against.
	Envs []Env `yaml:"Envs"`
}

// Env describes one environment (framework) a test runs against: where its
// server application lives and how to start it.
//
// An Env needs at least a Repo or a Name. Local or private code bases —
// for example an unreleased framework checked out at C:/github/iris-private —
// are declared with a Name and a Dir pointing anywhere on disk; no public
// repository is required.
type Env struct {
	// Name is the display name. When empty it is derived from the last
	// segment of Repo, title-cased (kataras/iris -> "Iris").
	Name string `yaml:"Name"`
	// Link is the URL behind the name in reports. When empty it falls
	// back to the GitHub page of Repo, or to no link at all for local
	// code bases without a Repo.
	Link string `yaml:"Link"`
	// Repo is the GitHub repository in owner/name form, e.g. kataras/iris.
	Repo string `yaml:"Repo"`
	// Dir is the directory containing the server application. It may be
	// absolute (e.g. C:/github/iris-private/_benchmarks/static) or
	// relative to the specification file. When empty it defaults to
	// <spec dir>/_code/<lowercase test name>/<lowercase env name>.
	Dir string `yaml:"Dir"`
	// Exec holds the commands that build and start the server, one per
	// line; the last line must be the (blocking) server process. When
	// empty it is inferred from Language:
	//
	//	Go          go run .
	//	C#          dotnet run -c Release
	//	Javascript  npm install
	//	            node .
	Exec string `yaml:"Exec"`
	// Language is the implementation language, e.g. Go, Javascript or C#.
	// Defaults to DefaultLanguage. It selects the default Exec and is
	// shown in reports.
	Language string `yaml:"Language"`
	// NotYetImplemented marks an environment whose application has not
	// been written yet; it is listed in reports but not benchmarked.
	NotYetImplemented bool `yaml:"NotYetImplemented"`
	// NotSupported marks an environment that cannot implement the test;
	// it is listed in reports but not benchmarked.
	NotSupported bool `yaml:"NotSupported"`
}

// CanBenchmark reports whether the environment participates in a benchmark
// run, i.e. it is neither NotSupported nor NotYetImplemented.
func (e Env) CanBenchmark() bool {
	return !e.NotSupported && !e.NotYetImplemented
}
