# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY cart-service/go.mod cart-service/go.sum* ./cart-service/

# Initialize isolated workspace for caching dependencies
RUN go work init ./pkg ./cart-service && go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY cart-service/ ./cart-service/

# Compile static binary
WORKDIR /app/cart-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/cart-service ./cmd/main.go

# Minimal production image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/cart-service /app/cart-service

EXPOSE 4003 50053

ENTRYPOINT ["/app/cart-service"]
