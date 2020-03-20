package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	timeLayout = "Jan 2, 2006 at 3:04pm (UTC)"

	readmeMarkdownFilename  = "README.md"
	resultsMarkdownFilename = "RESULTS.md"
)

var (
	waitRunDur = flag.Duration("wait-run", 3*time.Second, "wait time between tests")
	testsFile  = flag.String("i", "./tests.yml", "yaml file path contains the tests to run")
	outputDir  = flag.String("o", "./", "directory to save generaged RESULTS.md and README.md files")

	spreadsheetID = flag.String("g-spreadsheet", "", "Google Spreadsheet ID to send results")
	// Note: we can support user credentials and token generated by user's input
	// and saved on disk at first ran (kataras/sheets already supports it)
	// but let's accept only a service account unless otherwise requested.
	googleSecretFile = flag.String("g-secret", "client_secret.json", "if g-spreadsheet is set, service account credentials file should be provided")
)

// server-benchmarks --wait-run=3s -i ./tests.dev.yml -o ./dev -g-spreadsheet $GoogleSpreadsheetID -g-secret client_secret.json
func main() {
	flag.Parse()

	if _, err := os.Stat("/.dockerenv"); err == nil || os.IsExist(err) {
		fmt.Fprintf(os.Stdout, "Running through docker container...")
	}

	tests, err := readTests(*testsFile)
	catch(err)

	// TESTS
	for _, t := range tests {
		err = benchmark(t)
		catch(err)
	}

	os.MkdirAll(*outputDir, os.ModeDir)

	filename := filepath.Join(*outputDir, resultsMarkdownFilename)
	err = writeResults(filename, *spreadsheetID, *googleSecretFile, tests)
	catch(err)
	fmt.Fprintf(os.Stdout, "Generate: %s\n", filename)

	filename = filepath.Join(*outputDir, readmeMarkdownFilename)
	err = writeReadme(filename, tests)
	catch(err)
	fmt.Fprintf(os.Stdout, "Generate: %s\n", filename)
}

func readTests(filename string) ([]*Test, error) {
	testsFile, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var tests []*Test
	if err = yaml.NewDecoder(testsFile).Decode(&tests); err != nil {
		return nil, err
	}

	if err = testsFile.Close(); err != nil {
		return nil, err
	}

	if len(tests) == 0 {
		return nil, fmt.Errorf("no tests to run")
	}

	return tests, nil
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
