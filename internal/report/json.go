package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kataras/server-benchmarks/internal/bombardier"
	"github.com/kataras/server-benchmarks/internal/runner"
	"github.com/kataras/server-benchmarks/internal/spec"
	"github.com/kataras/server-benchmarks/internal/sysinfo"
)

// SchemaVersion is the current results.json schema version.
const SchemaVersion = 1

// resultsDocument is the schema of results.json: the machine-readable
// counterpart of RESULTS.md and the input of the `export` subcommand.
type resultsDocument struct {
	SchemaVersion int           `json:"schemaVersion"`
	GeneratedAt   time.Time     `json:"generatedAt"`
	System        sysinfo.Info  `json:"system"`
	Tests         []resultsTest `json:"tests"`
}

type resultsTest struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Envs        []resultsEnv `json:"envs"`
}

type resultsEnv struct {
	Name     string             `json:"name"`
	Language string             `json:"language,omitempty"`
	Link     string             `json:"link,omitempty"`
	Status   runner.Status      `json:"status"`
	Error    string             `json:"error,omitempty"`
	Result   *bombardier.Result `json:"result,omitempty"`
}

// marshalResults encodes d as an indented results.json document.
func marshalResults(d Data) ([]byte, error) {
	doc := resultsDocument{
		SchemaVersion: SchemaVersion,
		GeneratedAt:   d.GeneratedAt.UTC(),
		System:        d.System,
		Tests:         make([]resultsTest, 0, len(d.Tests)),
	}

	for _, tr := range d.Tests {
		rt := resultsTest{
			Name:        tr.Test.Name,
			Description: tr.Test.Description,
			Envs:        make([]resultsEnv, 0, len(tr.Envs)),
		}

		for _, er := range tr.Envs {
			re := resultsEnv{
				Name:     er.Env.Name,
				Language: er.Env.Language,
				Link:     er.Env.Link,
				Status:   er.Status,
				Result:   er.Result,
			}
			if er.Err != nil {
				re.Error = er.Err.Error()
			}

			rt.Envs = append(rt.Envs, re)
		}

		doc.Tests = append(doc.Tests, rt)
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", ResultsJSONFile, err)
	}

	return append(data, '\n'), nil
}

// LoadResultsFile reads a results.json document back into Data, so all
// reports can be re-rendered (`export`) without re-running the benchmarks.
func LoadResultsFile(path string) (Data, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Data{}, fmt.Errorf("read results: %w", err)
	}

	var doc resultsDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Data{}, fmt.Errorf("%s: %w", path, err)
	}

	if doc.SchemaVersion != SchemaVersion {
		return Data{}, fmt.Errorf("%s: unsupported schemaVersion %d (this build understands %d)",
			path, doc.SchemaVersion, SchemaVersion)
	}

	d := Data{
		GeneratedAt: doc.GeneratedAt,
		System:      doc.System,
		Tests:       make([]runner.TestReport, 0, len(doc.Tests)),
	}

	for _, rt := range doc.Tests {
		tr := runner.TestReport{
			Test: spec.Test{Name: rt.Name, Description: rt.Description},
			Envs: make([]runner.EnvReport, 0, len(rt.Envs)),
		}

		for _, re := range rt.Envs {
			er := runner.EnvReport{
				Env: spec.Env{
					Name:     re.Name,
					Language: re.Language,
					Link:     re.Link,
				},
				Status: re.Status,
				Result: re.Result,
			}
			if re.Error != "" {
				er.Err = errors.New(re.Error)
			}

			tr.Envs = append(tr.Envs, er)
		}

		d.Tests = append(d.Tests, tr)
	}

	return d, nil
}
