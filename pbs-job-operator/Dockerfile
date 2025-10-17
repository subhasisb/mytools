# Build the manager binary
FROM golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source (relies on .dockerignore to filter)
COPY . .

# Build
# the GOARCH has no default value to allow the binary to be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use Ubuntu 24.04 as base image to provide PBS tools and utilities
FROM ubuntu:24.04
WORKDIR /

# Install basic utilities and PBS client tools
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update \
    && apt-get -y upgrade \
    && apt-get -y install \
        ca-certificates \
        lsb-release \
        software-properties-common \
        wget \
        iputils-ping \
        jq \
        wget python3 openssl sudo openssh-client coreutils \
    && apt-get -y autoremove \
    && apt-get -y autoclean \
    && apt-get -y update

COPY --from=builder /workspace/manager .

ENV PATH=/opt/pbs/bin:$PATH

ENTRYPOINT ["/manager"]
