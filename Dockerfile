FROM ubuntu:18.04
# Use ubuntu because of .NET Core, alpine version for latest .NET Core fails to install.
LABEL MAINTAINER="Gerasimos (Makis) Maropoulos <kataras2006@hotmail.com>"
RUN apt-get update && \
    apt-get install -y curl wget

# Install nodejs
RUN curl -sL https://deb.nodesource.com/setup_13.x | bash && \
    apt-get install -y nodejs

# Install .NET Core
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*

RUN curl https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
    apt-add-repository https://packages.microsoft.com/ubuntu/18.04/prod && \
    apt-get update && \
    apt-get install -y dotnet-sdk-3.1 && \
    dotnet --version

# Install Go
RUN add-apt-repository ppa:longsleep/golang-backports && \
    apt update && \
    apt install -y golang-go

ENV GOPATH="/go"
ENV PATH="/go/bin:${PATH}"

RUN mkdir -p $GOPATH/src/server-benchmarks
RUN mkdir $GOPATH/bin

# Install Bombardier
RUN wget -O /go/bin/bombardier https://github.com/codesenberg/bombardier/releases/download/v1.2.4/bombardier-linux-amd64
RUN chmod +x /go/bin/bombardier

WORKDIR /go/src/server-benchmarks
# Cache node modules
#
# static test
COPY ./_code/static/express/package.json ./_code/static/express/package.json
RUN cd ./_code/static/express && npm install
COPY ./_code/static/koa/package.json ./_code/static/koa/package.json
RUN cd ./_code/static/koa && npm install
# parameterized test
COPY ./_code/parameterized/express/package.json ./_code/parameterized/express/package.json
RUN cd ./_code/parameterized/express && npm install
COPY ./_code/parameterized/koa/package.json ./_code/parameterized/koa/package.json
RUN cd ./_code/parameterized/koa && npm install
# rest test
COPY ./_code/rest/express/package.json ./_code/rest/express/package.json
RUN cd ./_code/rest/express && npm install
COPY ./_code/rest/koa/package.json ./_code/rest/koa/package.json
RUN cd ./_code/rest/koa && npm install

# Cache go modules, build and execute the binary
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY go.mod .
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN go install

VOLUME [ "/data" ]

ENTRYPOINT ["server-benchmarks", "-o", "/data"]

# docker build -t server-benchmarks:latest .
# docker run -v ${PWD}:/data server-benchmarks