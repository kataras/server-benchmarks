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
| [Go](https://golang.org) | go1.14.2 |
| [.Net Core](https://dotnet.microsoft.com/) | 3.1.102 |
| [Node.js](https://nodejs.org/) | v13.10.1 |

> Last updated: Apr 12, 2020 at 2:41am (UTC)

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
| [Iris](https://github.com/kataras/iris) | Go |205603 |606.83us |35.86MB |4.87s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |194821 |639.91us |33.06MB |5.13s |
| [Chi](https://github.com/pressly/chi) | Go |192223 |648.71us |33.54MB |5.20s |
| [Gin](https://github.com/gin-gonic/gin) | Go |187246 |666.13us |32.67MB |5.34s |
| [Echo](https://github.com/labstack/echo) | Go |185804 |671.59us |32.40MB |5.39s |
| [Martini](https://github.com/go-martini/martini) | Go |160296 |777.36us |27.98MB |6.24s |
| [Koa](https://github.com/koajs/koa) | Javascript |107892 |1.14ms |21.53MB |9.17s |
| [Express](https://github.com/expressjs/express) | Javascript |91181 |1.44ms |22.20MB |11.60s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |39010 |3.20ms |6.81MB |25.63s |

### Test:Parameterized

📖 Fires 250000 requests with a dynamic parameter of string, receives a hello text based on the parameter as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |192088 |649.16us |36.78MB |1.30s |
| [Chi](https://github.com/pressly/chi) | Go |187973 |663.14us |36.00MB |1.33s |
| [Echo](https://github.com/labstack/echo) | Go |177471 |703.68us |33.96MB |1.41s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |172411 |728.64us |31.78MB |1.46s |
| [Gin](https://github.com/gin-gonic/gin) | Go |169312 |740.39us |32.27MB |1.48s |
| [Martini](https://github.com/go-martini/martini) | Go |145707 |0.85ms |27.94MB |1.71s |
| [Koa](https://github.com/koajs/koa) | Javascript |74756 |1.67ms |15.99MB |3.35s |
| [Express](https://github.com/expressjs/express) | Javascript |65553 |2.20ms |15.57MB |4.41s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |38577 |3.24ms |7.40MB |6.48s |

### Test:Sessions

📖 Fires 250000 requests, sets a session and displays its value.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |102580 |1.22ms |35.22MB |2.44s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |74056 |1.69ms |34.07MB |3.39s |
| [Echo](https://github.com/labstack/echo) | Go |73521 |1.70ms |38.14MB |3.40s |
| [Chi](https://github.com/pressly/chi) | Go |67068 |1.86ms |34.77MB |3.73s |
| [Martini](https://github.com/go-martini/martini) | Go |66955 |1.86ms |34.81MB |3.73s |
| [Gin](https://github.com/gin-gonic/gin) | Go |60140 |2.10ms |24.95MB |4.20s |
| [Koa](https://github.com/koajs/koa) | Javascript |52796 |2.75ms |20.47MB |5.51s |
| [Express](https://github.com/expressjs/express) | Javascript |31982 |4.27ms |7.82MB |8.56s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |16548 |7.55ms |23.74MB |15.11s |

### Test:REST

📖 Fires 200000 requests with a dynamic parameter of int, sends JSON as request body and receives JSON as response.

| Name | Language | Reqs/sec | Latency | Throughput | Time To Complete |
|------|:---------|:---------|:--------|:-----------|:-----------------|
| [Iris](https://github.com/kataras/iris) | Go |147754 |843.88us |40.43MB |1.35s |
| [Chi](https://github.com/pressly/chi) | Go |141918 |0.88ms |38.15MB |1.41s |
| [Kestrel](https://github.com/dotnet/aspnetcore) | C# |136747 |0.92ms |39.68MB |1.47s |
| [Gin](https://github.com/gin-gonic/gin) | Go |136480 |0.91ms |37.32MB |1.47s |
| [Echo](https://github.com/labstack/echo) | Go |134209 |0.93ms |36.84MB |1.49s |
| [Martini](https://github.com/go-martini/martini) | Go |123638 |1.01ms |33.22MB |1.62s |
| [Buffalo](https://github.com/gobuffalo/buffalo) | Go |56722 |2.20ms |15.56MB |3.53s |
| [Koa](https://github.com/koajs/koa) | Javascript |47089 |2.66ms |13.91MB |4.26s |
| [Express](https://github.com/expressjs/express) | Javascript |41018 |3.29ms |13.59MB |5.28s |


## License

This project is licensed under the [MIT License](LICENSE).
