# cmd/server/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the server binary
RUN go build -o server ./cmd/server/main.go

# Command to run the server
CMD ["/app/server"]
