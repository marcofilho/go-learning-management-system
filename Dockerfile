FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API and migration tool
RUN go build -o api ./src/cmd/api
RUN go build -o migrate ./src/cmd/migrate

FROM alpine:latest

WORKDIR /app

# Copy binaries and migrations
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]
