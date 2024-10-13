# cmd/consumer/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the consumer binary
RUN go build -o consumer ./cmd/consumer/main.go

# Command to run the consumer
CMD ["/app/consumer"]
