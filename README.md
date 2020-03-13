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

> Last updated: Mar 11, 2020 at 9:52pm (UTC)

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
| [Iris](https://github.com/kataras/iris) | Go |230706.37 |540.64us |40.24MB |4.34s |
| [Gin](https://github.com/gin-gonic/gin) | Go |204165.70 |611.03us |35.60MB |4.90s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | CSharp |200807.10 |620.73us |34.07MB |4.98s |
| [Echo](https://github.com/labstack/echo) | Go |200720.61 |621.68us |34.99MB |4.99s |
| [Chi](https://github.com/pressly/chi) | Go |197347.35 |632.23us |34.41MB |5.07s |
| [Martini](https://github.com/go-martini/martini) | Go |165797.23 |751.13us |28.95MB |6.03s |

### Test:Parameterized

📖 Fires 250000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |210667.30 |592.87us |40.24MB |1.19s |
| [Chi](https://github.com/pressly/chi) | Go |201200.54 |621.27us |38.40MB |1.25s |
| [Echo](https://github.com/labstack/echo) | Go |195142.05 |639.42us |37.32MB |1.28s |
| [Gin](https://github.com/gin-gonic/gin) | Go |194082.20 |642.02us |37.21MB |1.29s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | CSharp |178642.40 |700.65us |33.02MB |1.41s |
| [Martini](https://github.com/go-martini/martini) | Go |160786.34 |776.21us |30.76MB |1.56s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |167687.80 |742.40us |45.89MB |1.19s |
| [Gin](https://github.com/gin-gonic/gin) | Go |154708.25 |812.40us |42.06MB |1.31s |
| [Chi](https://github.com/pressly/chi) | Go |152704.46 |817.39us |40.93MB |1.31s |
| [Echo](https://github.com/labstack/echo) | Go |151601.11 |822.82us |41.58MB |1.32s |
| [Martini](https://github.com/go-martini/martini) | Go |138567.52 |0.90ms |37.15MB |1.45s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | CSharp |127258.53 |0.98ms |39.09MB |1.58s |

## License

This software is licensed under [MIT License](LICENSE).