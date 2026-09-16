# --- Stage 1: Build the application ---
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o data-worker ./cmd/worker

# --- Stage 2: Lightweight runtime image ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/data-worker /app/data-worker
COPY internal/config/config.yaml /app/internal/config/config.yaml

RUN chmod +x /app/data-worker

EXPOSE 8083

ENTRYPOINT ["/app/data-worker"]