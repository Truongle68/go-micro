# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY order-service/go.mod order-service/go.sum* ./order-service/

# Initialize isolated workspace for caching dependencies
RUN go work init ./pkg ./order-service && go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY order-service/ ./order-service/

# Compile static binary
WORKDIR /app/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/order-service ./cmd/main.go

# Minimal production image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/order-service /app/order-service

EXPOSE 4004

ENTRYPOINT ["/app/order-service"]
