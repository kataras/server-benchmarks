# Contributing

First of all read our [Code of Conduct](CODE_OF_CONDUCT.md).

## Found a bug?

Open a new [issue](https://github.com/kataras/server-benchmarks/issues/new).
 * Write the Operating System and the version of your machine.
 * Describe your problem, what did you expect to see and what you see instead.
 * If it's a feature request, describe your idea as better as you can.

## Adding a framework

Only HTTP/2-featured web frameworks are acceptable here as this is the production line and the future of a secure environment.

1. Fork the [repository](https://github.com/kataras/server-benchmarks).
2. Make your changes.
    * Add the source code of the stress test at `./_code/%FRAMEWORK%/%TEST%` directory.
    * Edit the [./tests.yml](./tests.yml) configuration file accordingly.
        * See the available tests (e.g. "Parameterized") and add the framework on its "Envs" list.
            * If not available test for your stress-case please add a new, the configuration file is fully customizable.
        * Order does not matter, the fastest will be shown first and e.t.c.
3. Compare & Push the PR from [here](https://github.com/kataras/server-benchmarks/compare).
