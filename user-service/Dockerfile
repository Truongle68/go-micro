# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY user-service/go.mod user-service/go.sum* ./user-service/

# Initialize isolated workspace for caching dependencies
RUN go work init ./pkg ./user-service && go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY user-service/ ./user-service/

# Compile static binary
WORKDIR /app/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/user-service ./cmd/main.go

# Minimal production image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/user-service /app/user-service

EXPOSE 4000

ENTRYPOINT ["/app/user-service"]
