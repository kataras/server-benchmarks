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
	resultsMarkdownFile, err := os.OpenFile(markdownFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("markdown file: %w", err)
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
		_, err := client.ClearSpreadsheet(context.TODO(), spreadsheetID, "*")
		if err != nil {
			return err
		}

		// Server |      Test     | Reqs/sec | Latency
		// ___________________________________________
		// F1     | Static        | X        | X
		// F2     | Static        | X        | X
		// F1     | Parameterized | X        | X
		// F2     | Parameterized | X        | X
		records := [][]interface{}{
			{"Server", "Test", "Reqs/sec", "Latency"},
		}

		for _, t := range tests {
			for _, env := range t.Envs {
				if !env.CanBenchmark() {
					continue
				}

				records = append(records, []interface{}{
					env.Name, t.Name, fmt.Sprintf("%.0f", env.Result.RequestsPerSecond.Mean), fmt.Sprintf("%.2f", env.Result.Latency.Mean),
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

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("csv file: %w", err)
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
	readmeFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("readme file: %w", err)
	}
	defer readmeFile.Close()

	return rootTmpl.ExecuteTemplate(readmeFile, "readme", templateData{
		Datetime: time.Now().UTC().Format(timeLayout),
		System:   getSystemInfo(),
		Tests:    tests,
	})
}
