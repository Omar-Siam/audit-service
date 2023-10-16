# Use the official Golang image as a base image
FROM golang:1.21

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code into the container
COPY . .

# Build the Go app
RUN go build -o auditlog ./cmd/auditlog

# Command to run the application
CMD ["./auditlog"]
