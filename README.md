# Server Benchmarks

A benchmark suite which, **transparently**, stress-tests web servers and generates a report in markdown. It measures the requests per second, data transferred and time between requests and responses.

![Benchmarks: Jul 18, 2020 at 10:46am (UTC)](https://iris-go.com/images/benchmarks.svg)

## Why YABS (Yet Another Benchmark Suite)

It's true, there already enough of benchmark suites to play around. However, most of them don't even contain real-life test applications to benchmark, therefore the results are not always accurate e.g. a route handler executes SQL queries or reads and sends JSON. This benchmark suite is a fresh start, it can contain any type of tests as the tests are running as self-executables and the measuring is done by a popular and trusted 3rd-party software which acts as a real HTTP Client (one more reason of transparency). [Contributions](CONTRIBUTING.md) and improvements are always welcomed here.

## Use case

Measure the performance of application(s) between different versions or implementations (or web frameworks).

This suite can be further customized, through its [tests.yml](tests.yml) file, in order to test personal or internal web applications before their public releases.

## Installation

The only requirement for the benchmark tool is the [Go Programming Language](https://golang.org/dl).

```sh
$ go get -u github.com/kataras/server-benchmarks
$ go get -u github.com/codesenberg/bombardier
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

The following generated README contains benchmark results from builtin tests between popular **HTTP/2 web frameworks as of 2020**.

_Note:_ it's possible that the contents of this file will be updated regularly to accept even more tests cases and frameworks.

## System

|    |    |
|----|:---|
| Processor | Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz |
| RAM | 15.85 GB |
| OS | Microsoft Windows 10 Pro for Workstations |
| [Bombardier](https://github.com/codesenberg/bombardier) | v1.2.4 |
| [Go](https://golang.org) | go1.14.6 |
| [.Net Core](https://dotnet.microsoft.com/) | 3.1.102 |
| [Node.js](https://nodejs.org/) | v14.4.0 |

> Last updated: Jul 18, 2020 at 10:46am (UTC)

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
| [Iris](https://github.com/kataras/iris) | Go |206880 |602.75us |36.09MB |4.84s |
| [Gin](https://github.com/gin-gonic/gin) | Go |194837 |639.72us |34.00MB |5.13s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |193333 |644.46us |32.81MB |5.17s |
| [Chi](https://github.com/pressly/chi) | Go |187499 |664.76us |32.72MB |5.33s |
| [Echo](https://github.com/labstack/echo) | Go |185269 |673.04us |32.33MB |5.40s |
| [Martini](https://github.com/go-martini/martini) | Go |151410 |823.06us |26.43MB |6.60s |
| [Koa](https://github.com/koajs/koa) | Javascript |106631 |1.16ms |21.26MB |9.28s |
| [Express](https://github.com/expressjs/express) | Javascript |83514 |1.49ms |21.59MB |11.93s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |38466 |3.25ms |6.71MB |26.00s |

### Test:Parameterized

📖 Fires 550000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |193774 |643.09us |37.14MB |2.84s |
| [Chi](https://github.com/pressly/chi) | Go |184545 |676.70us |35.32MB |2.99s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |182190 |684.19us |33.87MB |3.02s |
| [Echo](https://github.com/labstack/echo) | Go |176197 |708.66us |33.73MB |3.13s |
| [Gin](https://github.com/gin-gonic/gin) | Go |174488 |714.74us |33.45MB |3.15s |
| [Martini](https://github.com/go-martini/martini) | Go |145396 |0.86ms |27.88MB |3.78s |
| [Koa](https://github.com/koajs/koa) | Javascript |88820 |1.37ms |19.46MB |6.06s |
| [Express](https://github.com/expressjs/express) | Javascript |73546 |1.85ms |18.48MB |8.17s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |37609 |3.32ms |7.21MB |14.62s |

### Test:Sessions

📖 Fires 250000 requests, sets a session and displays its value.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |106836 |1.17ms |33.41MB |2.34s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |78444 |1.60ms |36.11MB |3.20s |
| [Echo](https://github.com/labstack/echo) | Go |73867 |1.69ms |38.28MB |3.39s |
| [Chi](https://github.com/pressly/chi) | Go |68098 |1.83ms |35.31MB |3.67s |
| [Martini](https://github.com/go-martini/martini) | Go |67507 |1.85ms |35.03MB |3.70s |
| [Gin](https://github.com/gin-gonic/gin) | Go |57493 |2.18ms |24.07MB |4.36s |
| [Koa](https://github.com/koajs/koa) | Javascript |47820 |2.79ms |20.15MB |5.60s |
| [Express](https://github.com/expressjs/express) | Javascript |27617 |4.38ms |7.64MB |8.77s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |16810 |7.45ms |24.08MB |14.90s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |150430 |826.05us |41.25MB |1.33s |
| [Chi](https://github.com/pressly/chi) | Go |146274 |0.85ms |39.32MB |1.37s |
| [Gin](https://github.com/gin-gonic/gin) | Go |141664 |0.88ms |38.74MB |1.41s |
| [Echo](https://github.com/labstack/echo) | Go |138915 |0.90ms |38.15MB |1.44s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |136935 |0.91ms |39.79MB |1.47s |
| [Martini](https://github.com/go-martini/martini) | Go |128590 |0.97ms |34.57MB |1.56s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |58954 |2.12ms |16.18MB |3.40s |
| [Koa](https://github.com/koajs/koa) | Javascript |50948 |2.61ms |14.15MB |4.19s |
| [Express](https://github.com/expressjs/express) | Javascript |38451 |3.24ms |13.77MB |5.21s |


## License

This project is licensed under the [MIT License](LICENSE).
