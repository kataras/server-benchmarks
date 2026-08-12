# The image carries every runtime the benchmark apps need (Go, Node.js,
# .NET SDK) plus the bombardier load generator — the apps are started with
# `go run` / `node` / `dotnet run` at benchmark time, so a multi-stage
# build would gain nothing here.
FROM ubuntu:24.04

ARG TARGETARCH=amd64
ARG GO_VERSION=1.26.5
ARG NODE_MAJOR=24
ARG DOTNET_CHANNEL=10.0
ARG BOMBARDIER_VERSION=v2.0.2

LABEL org.opencontainers.image.title="server-benchmarks" \
      org.opencontainers.image.description="Benchmarks between HTTP web frameworks" \
      org.opencontainers.image.source="https://github.com/kataras/server-benchmarks" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.authors="Gerasimos (Makis) Maropoulos <contact@hellenic.dev>"

ENV DEBIAN_FRONTEND=noninteractive \
    DOTNET_CLI_TELEMETRY_OPTOUT=1 \
    DOTNET_NOLOGO=1

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates curl git \
    && rm -rf /var/lib/apt/lists/*

# Node.js LTS (NodeSource).
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_MAJOR}.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/* \
    && node --version

# .NET SDK via the official install script (the Microsoft apt feed is
# being phased out and lags behind on Ubuntu).
RUN curl -fsSL https://dot.net/v1/dotnet-install.sh | bash -s -- \
        --channel ${DOTNET_CHANNEL} --install-dir /usr/share/dotnet \
    && ln -s /usr/share/dotnet/dotnet /usr/local/bin/dotnet \
    && dotnet --version

# Go, from the official tarball (pinned; distro packages lag).
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Bombardier, the HTTP load generator doing the actual measuring.
RUN curl -fsSL -o /usr/local/bin/bombardier \
        https://github.com/codesenberg/bombardier/releases/download/${BOMBARDIER_VERSION}/bombardier-linux-${TARGETARCH} \
    && chmod +x /usr/local/bin/bombardier \
    && bombardier --version

WORKDIR /app

# Warm every dependency cache before copying the full sources, so code
# changes don't invalidate these slow layers.
COPY go.mod go.sum ./
RUN go mod download

COPY _code/ _code/
RUN set -e; for m in $(find _code -name go.mod); do \
        (cd "$(dirname "$m")" && go mod download); \
    done
RUN set -e; for d in _code/*/express _code/*/koa _code/*/fastify; do \
        (cd "$d" && npm ci --omit=dev --no-fund --no-audit); \
    done
RUN set -e; for p in $(find _code -name '*.csproj'); do \
        dotnet restore "$p"; \
    done

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /usr/local/bin/server-benchmarks .

VOLUME ["/data"]
ENTRYPOINT ["server-benchmarks", "-o", "/data", "-wait-run", "6s"]

# Build: docker build -t server-benchmarks .
# Run:   docker run --rm -v ${PWD}:/data server-benchmarks
