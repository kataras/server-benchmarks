package main

import (
	"fmt"
	"text/template"
)

type templateData struct {
	Datetime string
	System   systemInfo
	Tests    []*Test
}

var (
	rootTmpl *template.Template
)

func init() {
	rootTmpl = template.New("")

	template.Must(rootTmpl.New("results").Funcs(template.FuncMap{
		"formatTimeUs": formatTimeUs,
		"formatBinary": formatBinary,
	}).Parse(`## System

|    |    |
|----|:---|
| Processor | {{.System.Processor}} |
| RAM | {{.System.RAM}} |
| OS | {{.System.OS}} |
| [Bombardier](https://github.com/codesenberg/bombardier) | {{.System.Bombardier}} |
| [Go](https://golang.org) | {{.System.Go}} |
{{- if .System.Dotnet }}
| [.Net Core](https://dotnet.microsoft.com/) | {{.System.Dotnet}} |
{{ end -}}
{{- if .System.Node }}
| [Node.js](https://nodejs.org/) | {{.System.Node}} |
{{- end }}

> Last updated: {{.Datetime}}

## Terminology

**Name** is the name of the framework(or router) used under a particular test.

**Reqs/sec** is the avg number of total requests could be processed per second (the higher the better).

**Latency** is the amount of time it takes from when a request is made by the client to the time it takes for the response to get back to that client (the smaller the better).

**Throughput** is the rate of production or the rate at which data are transferred (the higher the better, it depends from response length (body + headers).

**Time To Complete** is the total time (in seconds) the test completed (the smaller the better).

## Results
{{ range $test := .Tests}}
### Test:{{ $test.Name}}

{{ if $test.Description -}}
📖 {{ $test.ParseDescription $test -}}
{{ end }}

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
{{ range $env := $test.Envs -}}
| [{{ $env.GetName }}]({{ $env.GetLink }}) | {{ $env.Language }} | 
{{- if $env.CanBenchmark }}
	{{- printf "%.0f" $env.Result.RequestsPerSecond.Mean }} | 
	{{- formatTimeUs $env.Result.Latency.Mean }} | 
	{{- formatBinary $env.Result.Throughput }} | 
	{{- printf "%.2f" $env.Result.TimeTakenSeconds }}s | 
{{- else -}}
- | - | - | - | - | - |
{{- end}}
{{end -}}
{{ end -}}
`))

	template.Must(rootTmpl.New("readme").Parse(`# Server Benchmarks

A benchmark suite which, **transparently**, stress-tests web servers and generates a report in markdown. It measures the requests per second, data transferred and time between requests and responses.

## Why YABS (Yet Another Benchmark Suite)

It's true, there already enough of benchmark suites to play around. However, most of them don't even contain real-life test applications to benchmark, therefore the results are not always accurate e.g. a route handler executes SQL queries or reads and sends JSON. This benchmark suite is a fresh start, it can contain any type of tests as the tests are running as self-executables and the measuring is done by a popular and trusted 3rd-party software which acts as a real HTTP Client (one more reason of transparency). [Contributions](CONTRIBUTING.md) and improvements are always welcomed here.

## Use case

Measure the performance of application(s) between different versions or implementations (or web frameworks).

This suite can be further customized, through its [tests.yml](https://github.com/kataras/server-benchmarks/blob/master/tests.yml) file, in order to test personal or internal web applications before their public releases.

## Installation

The only requirement for the benchmark tool is the [Go Programming Language](https://golang.org/dl).

` + "```sh" + `
$ go get -u github.com/kataras/server-benchmarks
$ go get -u github.com/codesenberg/bombardier
` + "```sh" + `

Depending on your test cases you may want to install [Node.js](https://nodejs.org/en/download/current/) and [.NET Core](https://dotnet.microsoft.com/download) too.

## How to run

1. Navigate to your tests directory, the one which includes a  **tests.yml** file
1. Open a temrinal and execute: ` + "`server-benchmarks`" + `
2. Wait for the executable _server-benchmarks_ (or _server-benchmarks.exe_ for windows) to finish
3. That's all, now open the the results from the generated **RESULTS.md** file.

### Advanced usage

- Read the tests from the _./tests.dev.yml_ file
- Wait 3 seconds between tests
- Output the results to the _./dev_ directory
- Write the results to a remote [google spreadsheet](https://www.google.com/sheets/about/) table, which you can convert to a graph later on (as shown above).

` + "```sh" + `
$ server-benchmarks --wait-run=3s -i ./tests.dev.yml -o ./dev -g-spreadsheet $GoogleSpreadsheetID -g-secret client_secret.json
` + "```" + `

### Run using Docker

The only requirement is [Docker](https://docs.docker.com/).

` + "```sh" + `
$ docker run -v ${PWD}:/data kataras/server-benchmarks
` + "```" + `

## Benchmarks

The following generated README contains benchmark results from builtin tests between popular **HTTP/2 web frameworks as of 2020**.

_Note:_ it's possible that the contents of this file will be updated regularly to accept even more tests cases and frameworks.

{{ template "results" .}}

## License

This project is licensed under the [MIT License](LICENSE).
`))

}

// copied from bombardier's source code itself to display identical results.
type units struct {
	scale uint64
	base  string
	units []string
}

var (
	binaryUnits = &units{
		scale: 1024,
		base:  "",
		units: []string{"KB", "MB", "GB", "TB", "PB"},
	}
	timeUnitsUs = &units{
		scale: 1000,
		base:  "us",
		units: []string{"ms", "s"},
	}
	timeUnitsS = &units{
		scale: 60,
		base:  "s",
		units: []string{"m", "h"},
	}
)

func formatUnits(n float64, m *units, prec int) string {
	amt := n
	unit := m.base

	scale := float64(m.scale) * 0.85

	for i := 0; i < len(m.units) && amt >= scale; i++ {
		amt /= float64(m.scale)
		unit = m.units[i]
	}
	return fmt.Sprintf("%.*f%s", prec, amt, unit)
}

func formatBinary(n float64) string {
	return formatUnits(n, binaryUnits, 2)
}

func formatTimeUs(n float64) string {
	units := timeUnitsUs
	if n >= 1000000.0 {
		n /= 1000000.0
		units = timeUnitsS
	}
	return formatUnits(n, units, 2)
}
