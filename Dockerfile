FROM golang:1.24-bullseye AS builder
WORKDIR /app

COPY go.* ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app .

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    ffmpeg && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app .

ENTRYPOINT ["/app/app"]
