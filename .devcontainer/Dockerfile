# Base Ubuntu
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Ferramentas básicas
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    bash-completion \
    unzip \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Install Go manual
ENV GO_VERSION=1.24.11
RUN curl -OL https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz \
    && rm go$GO_VERSION.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install mvdan.cc/gofumpt@latest

# Diretório de trabalho
WORKDIR /workspace
