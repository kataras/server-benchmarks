# Server Benchmarks

A benchmark suite which, **transparently**, stress-tests web servers and generates a report in markdown. It measures the requests per second, data transferred and time between requests and responses.

## Why YABS (Yet Another Benchmark Suite)

It's true, there already enough of benchmark suites to play around. However, most of them don't even contain real-life test applications to benchmark, therefore the results are not always accurate e.g. a route handler executes SQL queries or reads and sends JSON. This benchmark suite is a fresh start, it can contain any type of tests as the tests are running as self-executables and the measuring is done by a popular and trusted 3rd-party software which acts as a real HTTP Client (one more reason of transparency). [Contributions](CONTRIBUTING.md) and improvements are always welcomed here.

## Use case

Measure the performance of application(s) between different versions or implementations (or web frameworks).

This suite can be further customized, through its [tests.yml](tests.yml) file, in order to test personal or internal web applications before their public releases.

## How to run

1. Install [Go](https://golang.org/dl), [Bombardier](https://github.com/codesenberg/bombardier/releases/tag/v1.2.4), [Node.js](https://nodejs.org/en/download/current/) and [.NET Core](https://dotnet.microsoft.com/download)
2. Clone the repository
3. Stress-tests are described inside [tests.yml](tests.yml) file, it can be customized to fit your needs
4. Execute: `go build -o server-benchmarks`
5. Run and wait for the executable _server-benchmarks_ (or _server-benchmarks.exe_ for windows) to finish
6. Read the results from the generated _README.md_ file.

### Docker

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
| [Go](https://golang.org) | go1.14 |
| [.Net Core](https://dotnet.microsoft.com/) | 3.1.102 |
| [Node.js](https://nodejs.org/) | v13.10.1 |

> Last updated: Mar 18, 2020 at 2:03am (UTC)

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
| [Iris](https://github.com/kataras/iris) | Go |193538 |644.73us |33.75MB |5.17s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |185700 |671.60us |31.49MB |5.39s |
| [Gin](https://github.com/gin-gonic/gin) | Go |180037 |692.86us |31.41MB |5.56s |
| [Chi](https://github.com/pressly/chi) | Go |177093 |704.56us |30.89MB |5.65s |
| [Echo](https://github.com/labstack/echo) | Go |174810 |713.51us |30.49MB |5.72s |
| [Martini](https://github.com/go-martini/martini) | Go |147709 |843.08us |25.80MB |6.76s |
| [Koa](https://github.com/koajs/koa) | Javascript |104831 |1.20ms |20.43MB |9.66s |
| [Express](https://github.com/expressjs/express) | Javascript |81845 |1.51ms |21.16MB |12.17s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |35507 |3.52ms |6.20MB |28.17s |

### Test:Parameterized

📖 Fires 250000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |185449 |671.86us |35.53MB |1.35s |
| [Chi](https://github.com/pressly/chi) | Go |176367 |708.03us |33.70MB |1.42s |
| [Echo](https://github.com/labstack/echo) | Go |172300 |723.40us |33.01MB |1.45s |
| [Gin](https://github.com/gin-gonic/gin) | Go |164878 |756.54us |31.59MB |1.52s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |164297 |758.24us |30.51MB |1.52s |
| [Martini](https://github.com/go-martini/martini) | Go |143134 |0.87ms |27.42MB |1.75s |
| [Koa](https://github.com/koajs/koa) | Javascript |66423 |1.68ms |15.91MB |3.37s |
| [Express](https://github.com/expressjs/express) | Javascript |58290 |2.26ms |15.13MB |4.54s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |36316 |3.44ms |6.96MB |6.89s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |144702 |0.86ms |39.44MB |1.39s |
| [Chi](https://github.com/pressly/chi) | Go |140038 |0.89ms |37.61MB |1.43s |
| [Gin](https://github.com/gin-gonic/gin) | Go |134783 |0.93ms |36.92MB |1.49s |
| [Echo](https://github.com/labstack/echo) | Go |131910 |0.94ms |36.14MB |1.52s |
| [Martini](https://github.com/go-martini/martini) | Go |123382 |1.01ms |33.10MB |1.62s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |119648 |1.04ms |36.71MB |1.68s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |56020 |2.23ms |15.37MB |3.57s |
| [Koa](https://github.com/koajs/koa) | Javascript |46229 |2.70ms |13.67MB |4.34s |
| [Express](https://github.com/expressjs/express) | Javascript |36306 |3.20ms |13.94MB |5.14s |


## License

This project is licensed under the [MIT License](LICENSE).
