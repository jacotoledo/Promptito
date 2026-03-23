# Multi-stage Dockerfile optimized for Fly.io and general deployment

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o promptito ./cmd/promptito

RUN adduser -D -u 1000 appuser

# Runtime stage - Alpine for health check support
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 appuser

COPY --from=builder /build/promptito .
COPY --from=builder /build/public ./public

RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./promptito"]
CMD ["-prompts", "./public/prompts", "-static", "./public"]
