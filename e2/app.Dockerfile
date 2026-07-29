# App image for E2, built FROM the controlled mutable base tag.
# Build context is the repo root so the Go sources are visible.
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/action-test .

# The whole point of E2: a mutable tag that gets rewritten between phases.
FROM ghcr.io/duynhlab/action-test/e2base:latest

COPY --from=builder /out/action-test /usr/local/bin/action-test

USER 10001
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/action-test"]
