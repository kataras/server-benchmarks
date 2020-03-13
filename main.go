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

	readmeFilename  = "README.md"
	resultsFilename = "RESULTS.md"
)

var (
	waitServerDur = flag.Duration("wait-server", 6*time.Second, "wait time for server readiness")
	waitRunDur    = flag.Duration("wait-run", 3*time.Second, "wait time between tests")
	testsFile     = flag.String("i", "./tests.yml", "yaml file path contains the tests to run")
	outputDir     = flag.String("o", "./", "directory to save generaged RESULTS.md and README.md files")
)

// server-benchmarks --wait-server=6s --wait-run=3s -i ./tests.yml -o ./results
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

	// TODO: CSV format as well to transfer them to google docs for graphs.
	filename := filepath.Join(*outputDir, resultsFilename)
	err = writeResults(filename, tests)
	catch(err)
	fmt.Fprintf(os.Stdout, "Generate: %s\n", filename)
	filename = filepath.Join(*outputDir, readmeFilename)
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

func writeReadme(filename string, tests []*Test) error {
	os.MkdirAll(filepath.Dir(filename), os.ModeDir)

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

func writeResults(filename string, tests []*Test) error {
	os.MkdirAll(filepath.Dir(filename), os.ModeDir)

	resultsFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	return rootTmpl.ExecuteTemplate(resultsFile, "results", templateData{
		Datetime: time.Now().UTC().Format(timeLayout),
		System:   getSystemInfo(),
		Tests:    tests,
	})
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
