# Contributing

First of all, read our [Code of Conduct](CODE_OF_CONDUCT.md).

## Found a bug?

Open a new [issue](https://github.com/kataras/server-benchmarks/issues/new) with:

* your operating system and machine specs,
* what you expected to see and what you saw instead,
* for feature requests, a description of the idea.

## Adding a framework

Actively maintained frameworks with real-world usage are welcome.

1. Fork the [repository](https://github.com/kataras/server-benchmarks).
2. Add one self-contained app per test at `./_code/<test>/<framework>/`:
   * Go: `go.mod`, `go.sum` and `main.go`
   * Javascript: `package.json`, `package-lock.json` and `main.js`
   * C#: a `.csproj` and `Program.cs`
3. Add the framework to each test's `Envs` list in [tests.yml](tests.yml). Order does not matter; reports sort by measured speed.
4. Verify your pairing locally, e.g. `server-benchmarks run -t static.myframework`.
5. Open a [pull request](https://github.com/kataras/server-benchmarks/compare).

### The app contract

Every app listens on `localhost:5000` and answers:

| Test | Request | Expected response |
|------|:--------|:------------------|
| static | `GET /` | 200, body `Index` |
| parameterized | `GET /hello/{name}` | 200, body `Hello {name}` |
| rest | `POST /{id}` with a JSON array of `{name, language, id, bio, version}` | 200, JSON `{"id": <int>, "count": <array length>, "first_id": <first element's id>}` |

For the rest test: limit the request body to 2MB, answer 404 for a non-integer `{id}` and 400 for a malformed body. Use the framework's production settings (release mode, logging off), the way its own documentation recommends. Existing apps under [`_code/rest/`](_code/rest) are the reference.

Benchmarking an unreleased or private framework needs no public repository at all; see "Test a private or unreleased framework" in the [README](README.md).
