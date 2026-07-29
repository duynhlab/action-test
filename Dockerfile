# Multi-stage build for the action-test service.
#
# The runtime base is deliberately pinned to an OLD alpine release so that Trivy
# has real CVEs to find. This repository exists to measure the behaviour of the
# shared CI workflows, so a clean image would tell us nothing about the gate.
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Copy the module files first so dependency resolution is cached separately from
# the source. There are no external dependencies yet, but keep the layout right.
COPY go.mod ./
RUN go mod download

COPY main.go ./

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/action-test .

# ── Runtime ──────────────────────────────────────────────────────────────────
FROM alpine:3.17.0

RUN adduser -D -u 10001 app

COPY --from=builder /out/action-test /usr/local/bin/action-test

USER 10001
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/action-test"]
