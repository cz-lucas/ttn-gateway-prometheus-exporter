# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o exporter .

# Final stage
FROM busybox

WORKDIR /app

COPY --from=builder /app/exporter .

EXPOSE 9000

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD wget --spider http://localhost:9000/health

ENTRYPOINT ["./exporter"]