# Stage 1: Build the Go application with the latest Go version
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency caching
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go binary with static linking (no reliance on glibc)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .

# Stage 2: Create the final image with a minimal scratch image
FROM scratch

# Copy the statically compiled Go binary from the builder stage
COPY --from=builder /app/myapp /

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD ["/myapp"]
