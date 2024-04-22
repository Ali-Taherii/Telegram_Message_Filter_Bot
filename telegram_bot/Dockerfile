# Use an official Go runtime as a parent image
FROM golang:1.22.0 AS builder

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Download Go dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o tele_bot .

# Use a minimal Alpine Linux image as a base
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /go/src/app/tele_bot .

# Set environment variables for MySQL password and bot token
ENV MYSQL_PASSWORD your_mysql_password
ENV BOT_TOKEN your_bot_token

# Expose any necessary ports
# EXPOSE 8080

# Command to run the Go application
CMD ["./tele_bot"]
