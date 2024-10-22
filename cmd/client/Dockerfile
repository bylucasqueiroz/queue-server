# cmd/client/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the client binary
RUN go build -o client ./cmd/client/main.go

# Command to run the client
CMD ["/app/client"]
