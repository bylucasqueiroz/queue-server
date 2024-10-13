# cmd/producer/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the producer binary
RUN go build -o producer ./cmd/producer/main.go

# Command to run the producer
CMD ["/app/producer"]
