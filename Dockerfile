FROM golang:1.22.4

RUN apt update ; apt install -y git make jq curl vim htop ncat iputils-ping net-tools;
RUN git config --global --add safe.directory /app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2
RUN go install github.com/go-delve/delve/cmd/dlv@latest &&\
    go install github.com/amobe/gocov-merger@latest &&\
    go install github.com/nikolaydubina/go-cover-treemap@v1.4.2 &&\
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
