# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY catalog-service/go.mod catalog-service/go.sum* ./catalog-service/

# Initialize isolated workspace for caching dependencies
RUN go work init ./pkg ./catalog-service && go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY catalog-service/ ./catalog-service/

# Compile static binary
WORKDIR /app/catalog-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/catalog-service ./cmd/main.go

# Minimal production image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/catalog-service /app/catalog-service

EXPOSE 4001 50050

ENTRYPOINT ["/app/catalog-service"]
