# Server Benchmarks

[![CI](https://github.com/kataras/server-benchmarks/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/kataras/server-benchmarks/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kataras/server-benchmarks)](https://goreportcard.com/report/github.com/kataras/server-benchmarks)
[![Go Reference](https://pkg.go.dev/badge/github.com/kataras/server-benchmarks.svg)](https://pkg.go.dev/github.com/kataras/server-benchmarks)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A benchmark suite that stress-tests real web framework applications and writes the results as Markdown tables, SVG charts, CSV files and a machine-readable `results.json`.

The measuring is done by [bombardier](https://github.com/codesenberg/bombardier), a third-party HTTP load generator, not by this suite. Every framework runs as a self-contained program under [`_code/`](_code), written the way its own documentation recommends. Anyone can read the app code, run the same command and compare numbers.

## Frameworks

| Framework | Language | Version source |
|-----------|:---------|:---------------|
| [Iris](https://github.com/kataras/iris) | Go | `_code/*/iris/go.mod` |
| [Gin](https://github.com/gin-gonic/gin) | Go | `_code/*/gin/go.mod` |
| [Echo](https://github.com/labstack/echo) | Go | `_code/*/echo/go.mod` |
| [Chi](https://github.com/go-chi/chi) | Go | `_code/*/chi/go.mod` |
| [Fiber](https://github.com/gofiber/fiber) | Go | `_code/*/fiber/go.mod` |
| net/http (standard library baseline) | Go | `_code/*/nethttp/go.mod` |
| [Express](https://github.com/expressjs/express) | Javascript | `_code/*/express/package.json` |
| [Koa](https://github.com/koajs/koa) | Javascript | `_code/*/koa/package.json` |
| [Fastify](https://github.com/fastify/fastify) | Javascript | `_code/*/fastify/package.json` |
| [ASP.NET Core](https://github.com/dotnet/aspnetcore) minimal APIs | C# | `_code/*/aspnetcore/*.csproj` |

Versions are pinned in each app's manifest and kept current by Dependabot. Node apps run one worker per CPU through the builtin `node:cluster` module, since a single Node process is single-threaded; Go and .NET runtimes schedule across cores on their own.

## Tests

| Test | Request | What it exercises |
|------|:--------|:------------------|
| Static | `GET /` | routing, plain text response |
| Parameterized | `GET /hello/{name}` | dynamic path parameters |
| REST | `POST /{id}` with a ~1MB JSON array | routing constraints, JSON decode and encode (2MB body limit) |

Each test fires 100000 requests over 125 concurrent connections. All of it is configuration, not code: see the schema comment at the top of [tests.yml](tests.yml).

## Install

The tool itself only needs the [Go Programming Language](https://go.dev/dl/):

```sh
go install github.com/kataras/server-benchmarks@latest
go install github.com/codesenberg/bombardier@latest
```

[Node.js](https://nodejs.org/) 24+ and the [.NET SDK](https://dotnet.microsoft.com/download) 10 are needed only for the Javascript and C# environments. Without them, select the Go environments with `-t`, or let those environments fail and read the partial report.

## Usage

```sh
server-benchmarks                # run everything in ./tests.yml, write reports to ./
server-benchmarks list           # print the tests and environments without running
server-benchmarks version        # print version information
server-benchmarks help           # global help
```

The `run` subcommand (implied when omitted) accepts:

| Flag | Default | Meaning |
|------|:--------|:--------|
| `-i` | `./tests.yml` | spec file to read |
| `-o` | `./` | output directory for the reports |
| `-t` | (all) | run only a test or a single `test.env`; repeatable |
| `-wait-run` | `3s` | idle time between two benchmarks |
| `-readme` | off | also generate a publishable README.md |
| `-bombardier` | from PATH | path of the bombardier executable |
| `-v` | off | debug logging, including the servers' own output |

```sh
# only the REST test, plus one specific pairing:
server-benchmarks run -t rest -t static.iris

# a development spec into a separate directory:
server-benchmarks run -i ./tests.dev.yml -o ./dev -wait-run 5s
```

A finished run leaves behind `RESULTS.md`, one `<test>.csv` and `<test>_latency.csv` per test, one `<test>.svg` and `<test>_latency.svg` chart per test, and `results.json`. The charts are embedded in `RESULTS.md` and adapt to light and dark GitHub themes.

`export` re-renders every artifact from a saved `results.json`, so you can regenerate reports after a template change without benchmarking again:

```sh
server-benchmarks export -i ./results.json -o ./
```

### Test a private or unreleased framework

An environment does not need a public repository. Point it at any local directory that contains a runnable server implementing the test's contract:

```yaml
Envs:
  - Name: Iris (next)
    Dir: C:/github/iris-private/_benchmarks/static
```

For an unreleased version of a framework that already has apps here, copy the app instead (e.g. `_code/static/iris` to `_code/static/iris-next`) and add a `replace` directive to its `go.mod`:

```
replace github.com/kataras/iris/v12 => C:/github/iris-private
```

Then run just that pairing: `server-benchmarks run -t "static.iris (next)"`.

## Run with Docker

The image bundles Go, Node.js, the .NET SDK and bombardier, so the full suite runs without installing anything else:

```sh
docker build -t server-benchmarks .
docker run --rm -v ${PWD}:/data server-benchmarks
```

Reports land in the mounted directory. `${PWD}` works in both PowerShell and Unix shells.

## Results

Results are machine-generated per run and depend heavily on the hardware. A fresh official run with the current roster has not been published yet; run the suite yourself, or expand the section below for the last published numbers.

<details>
<summary>Historical results (previous roster, Jun 20, 2022)</summary>

These numbers predate the current roster (they include Martini and Buffalo, since removed, and the Kestrel app replaced by ASP.NET Core minimal APIs) and ran on Go 1.19beta1, .NET 6 and Node 18 with different request counts.

### System

|    |    |
|----|:---|
| Processor | AMD Ryzen 9 4900HS with Radeon Graphics |
| RAM | 15.42 GB |
| OS | Microsoft Windows 11 Pro |
| [Bombardier](https://github.com/codesenberg/bombardier) | v1.2.4 |
| [Go](https://go.dev) | go1.19beta1 |
| [.NET](https://dotnet.microsoft.com/) | 6.0.300 |
| [Node.js](https://nodejs.org/) | v18.2.0 |

### Test: Static

📖 Fires 1000000 requests, receives a static message as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go | 284059 | 438.34us | 49.58MB | 3.52s |
| [Chi](https://github.com/go-chi/chi) | Go | 275525 | 451.01us | 48.18MB | 3.62s |
| [Echo](https://github.com/labstack/echo) | Go | 267815 | 466.16us | 46.64MB | 3.74s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# | 263479 | 472.72us | 44.68MB | 3.80s |
| [Gin](https://github.com/gin-gonic/gin) | Go | 263399 | 472.70us | 45.98MB | 3.80s |
| [Martini](https://github.com/go-martini/martini) | Go | 233051 | 534.43us | 40.68MB | 4.29s |
| [Koa](https://github.com/koajs/koa) | Javascript | 131274 | 0.93ms | 29.24MB | 7.50s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go | 78963 | 1.58ms | 13.78MB | 12.66s |
| [Express](https://github.com/expressjs/express) | Javascript | 41078 | 3.02ms | 11.54MB | 24.22s |

### Test: Parameterized

📖 Fires 550000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go | 277099 | 449.55us | 53.07MB | 1.99s |
| [Chi](https://github.com/go-chi/chi) | Go | 272434 | 456.62us | 52.21MB | 2.02s |
| [Echo](https://github.com/labstack/echo) | Go | 261467 | 476.01us | 50.14MB | 2.10s |
| [Gin](https://github.com/gin-gonic/gin) | Go | 259308 | 480.32us | 49.70MB | 2.12s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# | 233843 | 534.73us | 43.34MB | 2.36s |
| [Martini](https://github.com/go-martini/martini) | Go | 225790 | 551.37us | 43.29MB | 2.44s |
| [Koa](https://github.com/koajs/koa) | Javascript | 114667 | 1.08ms | 27.21MB | 4.78s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go | 76747 | 1.63ms | 14.71MB | 7.17s |
| [Express](https://github.com/expressjs/express) | Javascript | 37110 | 3.32ms | 11.11MB | 14.69s |

### Test: REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go | 238954 | 521.69us | 64.15MB | 0.84s |
| [Gin](https://github.com/gin-gonic/gin) | Go | 229665 | 541.96us | 62.86MB | 0.87s |
| [Chi](https://github.com/go-chi/chi) | Go | 228072 | 545.78us | 62.61MB | 0.88s |
| [Echo](https://github.com/labstack/echo) | Go | 224491 | 553.84us | 61.70MB | 0.89s |
| [Martini](https://github.com/go-martini/martini) | Go | 198166 | 627.46us | 54.47MB | 1.01s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# | 163486 | 766.90us | 47.42MB | 1.23s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go | 102478 | 1.22ms | 28.14MB | 1.95s |
| [Koa](https://github.com/koajs/koa) | Javascript | 48425 | 2.56ms | 15.39MB | 4.14s |
| [Express](https://github.com/expressjs/express) | Javascript | 23622 | 5.25ms | 9.04MB | 8.41s |

</details>

## Terminology

**Name** is the name of the framework (or router) used under a particular test.

**Reqs/sec** is the average number of requests processed per second (the higher the better).

**Latency** is the time from when a request is made by the client until the response gets back to that client (the lower the better).

**Throughput** is the rate at which data is transferred (the higher the better; it depends on the response length, body plus headers).

**Time To Complete** is the total time the test took to complete (the lower the better).

## Contributing

Want a framework added? Read [CONTRIBUTING.md](CONTRIBUTING.md); it takes one app directory per test and one entry in `tests.yml`.

## License

This project is licensed under the [MIT License](LICENSE).
