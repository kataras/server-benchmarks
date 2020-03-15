package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kataras/sheets"
)

func writeResults(markdownFilename, spreadsheetID, googleSecretFile string, tests []*Test) error {
	resultsMarkdownFile, err := os.Create(markdownFilename)
	if err != nil {
		return err
	}
	err = rootTmpl.ExecuteTemplate(resultsMarkdownFile, "results", templateData{
		Datetime: time.Now().UTC().Format(timeLayout),
		System:   getSystemInfo(),
		Tests:    tests,
	})
	resultsMarkdownFile.Close()
	if err != nil {
		return err
	}

	if spreadsheetID != "" && googleSecretFile != "" {
		client := sheets.NewClient(sheets.ServiceAccount(context.TODO(), googleSecretFile, sheets.ScopeReadWrite))
		// TODO:
		//
		// Server |      Test     | Reqs/sec | Latency
		// ___________________________________________
		// F1     | Static        | 1        | 1
		// F2     | Static        | 1        | 1
		// F1     | Parameterized | 1        | 1
		// F2     | Parameterized | 1        | 1
		//
		// OR
		//
		//      Test     | Server | Reqs/sec | Lantency
		// ____________________________________________
		// Static        | F1     | 1        | 1
		// Static        | F2     | 1        | 1
		// Parameterized | F1     | 1        | 1
		// Parameterized | F2     | 1        | 1
		_ = client
	}

	for _, t := range tests {
		reqsPerSecond := func(r *TestResult) string {
			return fmt.Sprintf("%.0f", r.RequestsPerSecond.Mean)
		}
		if err = writeCSV(t, "Reqs/sec", reqsPerSecond, ""); err != nil {
			return err
		}

		latency := func(r *TestResult) string {
			return fmt.Sprintf("%.2f", r.Latency.Mean) // formatTimeUs(r.Latency.Mean)
		}
		if err = writeCSV(t, "Latency", latency, "_latency"); err != nil {
			return err
		}
	}

	return nil
}

func writeCSV(t *Test, header string, valueFn func(t *TestResult) string, fileSuffix string) error {
	filename := filepath.Join(*outputDir, strings.ToLower(t.Name+fileSuffix)+".csv")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(f)
	csvWriter.Write([]string{"Name", header})

	for _, env := range t.Envs {
		if !env.CanBenchmark() {
			continue
		}

		csvWriter.Write([]string{env.GetName(), valueFn(env.Result)})
	}

	csvWriter.Flush()
	return f.Close()
}

func writeReadme(filename string, tests []*Test) error {
	readmeFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	return rootTmpl.ExecuteTemplate(readmeFile, "readme", templateData{
		Datetime: time.Now().UTC().Format(timeLayout),
		System:   getSystemInfo(),
		Tests:    tests,
	})
}
