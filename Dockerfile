FROM golang:1.24 AS builder

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
# Turn off CGO to ensure static binaries
RUN CGO_ENABLED=0 go build -v -o bot cmd/bot/main.go

FROM alpine AS production

# Move to working directory /app
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/bot ./

# Start the application
CMD ["/app/bot"]