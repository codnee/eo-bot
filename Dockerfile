# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application (optimized release build)
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -trimpath -o /app/bin/eo-bot ./src

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy the binary from the builder
COPY --from=builder /app/bin/eo-bot /app/eo-bot

# Run the application
CMD ["/app/eo-bot"]
