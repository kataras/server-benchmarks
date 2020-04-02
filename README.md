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
| [Go](https://golang.org) | go1.14.1 |
| [.Net Core](https://dotnet.microsoft.com/) | 3.1.102 |
| [Node.js](https://nodejs.org/) | v13.10.1 |

> Last updated: Apr 2, 2020 at 12:13pm (UTC)

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
| [Iris](https://github.com/kataras/iris) | Go |196926 |633.38us |34.35MB |5.08s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |183808 |678.28us |31.19MB |5.44s |
| [Gin](https://github.com/gin-gonic/gin) | Go |179364 |695.59us |31.28MB |5.58s |
| [Chi](https://github.com/pressly/chi) | Go |179312 |695.57us |31.29MB |5.58s |
| [Echo](https://github.com/labstack/echo) | Go |177640 |702.52us |30.98MB |5.63s |
| [Martini](https://github.com/go-martini/martini) | Go |149662 |832.78us |26.12MB |6.68s |
| [Koa](https://github.com/koajs/koa) | Javascript |103589 |1.19ms |20.65MB |9.56s |
| [Express](https://github.com/expressjs/express) | Javascript |86151 |1.52ms |21.16MB |12.17s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |35823 |3.49ms |6.25MB |27.92s |

### Test:Parameterized

📖 Fires 250000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |192633 |648.16us |36.81MB |1.30s |
| [Chi](https://github.com/pressly/chi) | Go |182947 |682.43us |35.01MB |1.37s |
| [Echo](https://github.com/labstack/echo) | Go |176633 |706.68us |33.77MB |1.42s |
| [Gin](https://github.com/gin-gonic/gin) | Go |173763 |720.73us |33.12MB |1.45s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |156261 |798.38us |29.02MB |1.60s |
| [Martini](https://github.com/go-martini/martini) | Go |144401 |0.86ms |27.75MB |1.73s |
| [Koa](https://github.com/koajs/koa) | Javascript |76815 |1.63ms |16.41MB |3.27s |
| [Express](https://github.com/expressjs/express) | Javascript |54752 |2.20ms |15.59MB |4.41s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |36779 |3.40ms |7.04MB |6.81s |

### Test:Sessions

📖 Fires 250000 requests, sets a session and displays its value.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |98570 |1.27ms |33.75MB |2.54s |
| [Echo](https://github.com/labstack/echo) | Go |72485 |1.72ms |37.58MB |3.45s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |70600 |1.77ms |32.63MB |3.54s |
| [Martini](https://github.com/go-martini/martini) | Go |65812 |1.90ms |34.11MB |3.80s |
| [Chi](https://github.com/pressly/chi) | Go |65760 |1.90ms |34.10MB |3.80s |
| [Gin](https://github.com/gin-gonic/gin) | Go |57697 |2.17ms |24.19MB |4.34s |
| [Koa](https://github.com/koajs/koa) | Javascript |41894 |2.81ms |19.94MB |5.65s |
| [Express](https://github.com/expressjs/express) | Javascript |29781 |4.35ms |7.68MB |8.72s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |17857 |7.00ms |25.62MB |14.00s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |143094 |0.87ms |39.10MB |1.40s |
| [Chi](https://github.com/pressly/chi) | Go |137219 |0.91ms |36.87MB |1.46s |
| [Gin](https://github.com/gin-gonic/gin) | Go |133730 |0.93ms |36.59MB |1.50s |
| [Echo](https://github.com/labstack/echo) | Go |130751 |0.95ms |35.90MB |1.53s |
| [Martini](https://github.com/go-martini/martini) | Go |119487 |1.04ms |32.11MB |1.67s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |115142 |1.08ms |35.37MB |1.74s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |55439 |2.25ms |15.21MB |3.61s |
| [Express](https://github.com/expressjs/express) | Javascript |43582 |3.30ms |13.57MB |5.29s |
| [Koa](https://github.com/koajs/koa) | Javascript |41394 |2.75ms |13.48MB |4.40s |


## License

This project is licensed under the [MIT License](LICENSE).
