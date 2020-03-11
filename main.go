package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	timeLayout = "Jan 2, 2006 at 3:04pm (UTC)"

	inFilename  = "./tests.yml"
	outFilename = "./README.md"
)

func main() {
	tests, err := readTests(inFilename)
	catch(err)

	// TESTS
	for _, t := range tests {
		err = benchmark(t)
		catch(err)
	}

	// TODO:
	// 1. CSV format as well to transfer them to google docs for graphs.
	// 2. Template: Mark with bold the winners of each test, find the higher reqs/sec
	//    (time-to-complete can be fixed if Test Duration is provided so we can't depend on that).
	err = writeReadme(outFilename, tests)
	catch(err)
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
	readmeFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer readmeFile.Close()

	return readmeTmpl.Execute(readmeFile, struct {
		Datetime string
		System   systemInfo
		Tests    []*Test
	}{
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
