FROM golang:1.26

RUN apt update ; apt install -y git make jq curl vim htop ncat iputils-ping net-tools;
RUN FZF_VERSION="0.72.0" && \
    wget "https://github.com/junegunn/fzf/releases/download/v${FZF_VERSION}/fzf-${FZF_VERSION}-linux_amd64.tar.gz" && \
    tar -xzf fzf-${FZF_VERSION}-linux_amd64.tar.gz -C /tmp && rm fzf-${FZF_VERSION}-linux_amd64.tar.gz && \
    mv /tmp/fzf /usr/local/bin/

RUN git config --global --add safe.directory /app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2
RUN go install github.com/go-delve/delve/cmd/dlv@v1.26.0 &&\
    go install github.com/amobe/gocov-merger@v1.0.0 &&\
    go install github.com/nikolaydubina/go-cover-treemap@v1.4.2 &&\
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod tidy
