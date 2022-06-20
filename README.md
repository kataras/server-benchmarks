# Server Benchmarks

A benchmark suite which, **transparently**, stress-tests web servers and generates a report in markdown. It measures the requests per second, data transferred and time between requests and responses.

![Benchmarks: Jun 20, 2020 at 8:17pm (UTC)](https://iris-go.com/images/benchmarks.svg)

## Why YABS (Yet Another Benchmark Suite)

It's true, there already enough of benchmark suites to play around. However, most of them don't even contain real-life test applications to benchmark, therefore the results are not always accurate e.g. a route handler executes SQL queries or reads and sends JSON. This benchmark suite is a fresh start, it can contain any type of tests as the tests are running as self-executables and the measuring is done by a popular and trusted 3rd-party software which acts as a real HTTP Client (one more reason of transparency). [Contributions](CONTRIBUTING.md) and improvements are always welcomed here.

## Use case

Measure the performance of application(s) between different versions or implementations (or web frameworks).

This suite can be further customized, through its [tests.yml](tests.yml) file, in order to test personal or internal web applications before their public releases.

## Installation

The only requirement for the benchmark tool is the [Go Programming Language](https://go.dev/dl/).

```sh
$ go get github.com/kataras/server-benchmarks@latest
$ go install github.com/codesenberg/bombardier@latest
```

Depending on your test cases you may want to install [Node.js](https://nodejs.org/en/download/current/) and [.NET Core](https://dotnet.microsoft.com/download) too.

## How to run

1. Navigate to your tests directory, the one which includes a  **tests.yml** file
1. Open a terminal and execute: `server-benchmarks`
2. Wait for the executable _server-benchmarks_ (or _server-benchmarks.exe_ for windows) to finish
3. That's all, now open the the results from the generated **RESULTS.md** file.

### Advanced usage

- Read the tests from the _./tests.dev.yml_ file
- Wait 3 seconds between tests
- Output the results to the _./dev_ directory
- Write the results to a remote [google spreadsheet](https://www.google.com/sheets/about/) table, which you can convert to a graph later on (as shown above).

```sh
$ server-benchmarks --wait-run=3s -i ./tests.dev.yml -o ./dev -g-spreadsheet $GoogleSpreadsheetID -g-secret client_secret.json
```

### Run using Docker

The only requirement is [Docker](https://docs.docker.com/).

```sh
$ docker run -v ${PWD}:/data kataras/server-benchmarks
```

## Benchmarks

The following generated README contains benchmark results from builtin tests between popular **HTTP/2 web frameworks as of 2022**.

_Note:_ it's possible that the contents of this file will be updated regularly to accept even more tests cases and frameworks.

## System

|    |    |
|----|:---|
| Processor | AMD Ryzen 9 4900HS with Radeon Graphics          |
| RAM | 15.42 GB |
| OS | Microsoft Windows 11 Pro |
| [Bombardier](https://github.com/codesenberg/bombardier) | v1.2.4 |
| [Go](https://golang.org) | go1.19beta1 |
| [.Net Core](https://dotnet.microsoft.com/) | 6.0.300 |

| [Node.js](https://nodejs.org/) | v18.2.0 |

> Last updated: Jun 20, 2022 at 8:17pm (UTC)

## Terminology

**Name** is the name of the framework(or router) used under a particular test.

**Reqs/sec** is the avg number of total requests could be processed per second (the higher the better).

**Latency** is the amount of time it takes from when a request is made by the client to the time it takes for the response to get back to that client (the smaller the better).

**Throughput** is the rate of production or the rate at which data are transferred (the higher the better, it depends from response length (body + headers).

**Time To Complete** is the total time (in seconds) the test completed (the smaller the better).

## Results

### Test:Static

📖 Fires 1000000 requests, receives a static message as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |284059 |438.34us |49.58MB |3.52s |
| [Chi](https://github.com/go-chi/chi) | Go |275525 |451.01us |48.18MB |3.62s |
| [Echo](https://github.com/labstack/echo) | Go |267815 |466.16us |46.64MB |3.74s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |263479 |472.72us |44.68MB |3.80s |
| [Gin](https://github.com/gin-gonic/gin) | Go |263399 |472.70us |45.98MB |3.80s |
| [Martini](https://github.com/go-martini/martini) | Go |233051 |534.43us |40.68MB |4.29s |
| [Koa](https://github.com/koajs/koa) | Javascript |131274 |0.93ms |29.24MB |7.50s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |78963 |1.58ms |13.78MB |12.66s |
| [Express](https://github.com/expressjs/express) | Javascript |41078 |3.02ms |11.54MB |24.22s |

### Test:Parameterized

📖 Fires 550000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |277099 |449.55us |53.07MB |1.99s |
| [Chi](https://github.com/go-chi/chi) | Go |272434 |456.62us |52.21MB |2.02s |
| [Echo](https://github.com/labstack/echo) | Go |261467 |476.01us |50.14MB |2.10s |
| [Gin](https://github.com/gin-gonic/gin) | Go |259308 |480.32us |49.70MB |2.12s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |233843 |534.73us |43.34MB |2.36s |
| [Martini](https://github.com/go-martini/martini) | Go |225790 |551.37us |43.29MB |2.44s |
| [Koa](https://github.com/koajs/koa) | Javascript |114667 |1.08ms |27.21MB |4.78s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |76747 |1.63ms |14.71MB |7.17s |
| [Express](https://github.com/expressjs/express) | Javascript |37110 |3.32ms |11.11MB |14.69s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |238954 |521.69us |64.15MB |0.84s |
| [Gin](https://github.com/gin-gonic/gin) | Go |229665 |541.96us |62.86MB |0.87s |
| [Chi](https://github.com/go-chi/chi) | Go |228072 |545.78us |62.61MB |0.88s |
| [Echo](https://github.com/labstack/echo) | Go |224491 |553.84us |61.70MB |0.89s |
| [Martini](https://github.com/go-martini/martini) | Go |198166 |627.46us |54.47MB |1.01s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |163486 |766.90us |47.42MB |1.23s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |102478 |1.22ms |28.14MB |1.95s |
| [Koa](https://github.com/koajs/koa) | Javascript |48425 |2.56ms |15.39MB |4.14s |
| [Express](https://github.com/expressjs/express) | Javascript |23622 |5.25ms |9.04MB |8.41s |

## License

This project is licensed under the [MIT License](LICENSE).
