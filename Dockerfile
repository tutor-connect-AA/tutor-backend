# Use the official Golang image as the base image
FROM golang:1.22.3

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o myapp .

# Expose the application port
EXPOSE 8080

# Start the application
CMD ["./myapp"]
