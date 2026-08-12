package runner

import (
	"cmp"
	"slices"
	"strings"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/spec"
)

// Status classifies the outcome of one environment in a benchmark run.
type Status string

const (
	// StatusOK means the benchmark ran and its results validated.
	StatusOK Status = "ok"
	// StatusFailed means the environment failed to start, to be measured
	// or to validate.
	StatusFailed Status = "failed"
	// StatusSkipped means the environment was not benchmarked because it
	// is marked NotSupported or NotYetImplemented.
	StatusSkipped Status = "skipped"
)

// EnvReport is the outcome of a single environment within a test.
type EnvReport struct {
	Env    spec.Env
	Status Status
	// Result holds the measurements; non-nil only when Status is StatusOK.
	Result *bombardier.Result
	// Err describes the failure; non-nil only when Status is StatusFailed.
	Err error
}

// TestReport is the outcome of one test across all its environments,
// sorted best-first: successful runs by descending requests per second,
// then failed, then skipped environments; ties break by name.
type TestReport struct {
	Test spec.Test
	Envs []EnvReport
}

// Winner returns the best-performing environment of the test, if any
// environment succeeded at all.
func (r TestReport) Winner() (EnvReport, bool) {
	if len(r.Envs) > 0 && r.Envs[0].Status == StatusOK {
		return r.Envs[0], true
	}

	return EnvReport{}, false
}

// statusRank orders env reports: successful first, then failed, then skipped.
func statusRank(s Status) int {
	switch s {
	case StatusOK:
		return 0
	case StatusFailed:
		return 1
	default:
		return 2
	}
}

// sortEnvReports sorts reports best-first. It never touches Result of
// non-successful entries, whose Result is nil.
func sortEnvReports(reports []EnvReport) {
	slices.SortStableFunc(reports, func(a, b EnvReport) int {
		if c := cmp.Compare(statusRank(a.Status), statusRank(b.Status)); c != 0 {
			return c
		}

		if a.Status == StatusOK {
			if c := cmp.Compare(b.Result.RequestsPerSecond.Mean, a.Result.RequestsPerSecond.Mean); c != 0 {
				return c
			}
		}

		return strings.Compare(strings.ToLower(a.Env.Name), strings.ToLower(b.Env.Name))
	})
}
