# Use official Golang image to build the binary
FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use a minimal alpine image for the final container
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Install ca-certificates just in case
RUN apk --no-cache add ca-certificates

# Expose port (Cloud Run uses 8080 by default)
EXPOSE 8080

# Run the web service
CMD ["./server"]
