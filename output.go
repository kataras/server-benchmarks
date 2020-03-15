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
	// Write results as table to the RESULTS.md file.
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

	// Write results to CSV files per test.
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

	// Push results to a remote google spreadsheet.
	if spreadsheetID != "" && googleSecretFile != "" {
		client := sheets.NewClient(sheets.ServiceAccount(context.TODO(), googleSecretFile, sheets.ScopeReadWrite))
		// Server |      Test     | Reqs/sec | Latency
		// ___________________________________________
		// F1     | Static        | X        | X
		// F2     | Static        | X        | X
		// F1     | Parameterized | X        | X
		// F2     | Parameterized | X        | X

		records := [][]interface{}{
			{"Server", "Test", "Reqs/sec", "Latency"},
		}

		_, err := client.ClearSpreadsheet(context.TODO(), spreadsheetID, "*")
		if err != nil {
			return err
		}

		for _, t := range tests {
			for _, env := range t.Envs {
				if !env.CanBenchmark() {
					continue
				}

				records = append(records, []interface{}{
					env.Name, t.Name, env.Result.RequestsPerSecond.Mean, env.Result.Latency.Mean,
				})
			}
		}

		_, err = client.UpdateSpreadsheet(context.TODO(), spreadsheetID, sheets.ValueRange{Values: records})
		if err != nil {
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
