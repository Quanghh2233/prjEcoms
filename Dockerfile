FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the whole project
COPY . .

# Download dependencies
RUN go mod download

RUN go build -o api ./cmd/api

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api ./

# Set environment variables
ENV INIT_DB=true

# Expose the API port
EXPOSE 8080

CMD ["./api"]
