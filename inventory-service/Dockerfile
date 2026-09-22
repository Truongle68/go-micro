# syntax=docker/dockerfile:1
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY inventory-service/go.mod inventory-service/go.sum* ./inventory-service/

# Initialize isolated workspace for caching dependencies
RUN go work init ./pkg ./inventory-service && go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY inventory-service/ ./inventory-service/

# Compile static binaries for both the API service and outbox publisher
WORKDIR /app/inventory-service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/inventory-service ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/event-consumer ./cmd/event-consumer/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/outbox-publisher ./cmd/outbox-publisher/main.go

# Minimal production image
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bin/inventory-service /app/inventory-service
COPY --from=builder /app/bin/event-consumer /app/event-consumer
COPY --from=builder /app/bin/outbox-publisher /app/outbox-publisher

EXPOSE 4002 50052

CMD ["/app/inventory-service"]
