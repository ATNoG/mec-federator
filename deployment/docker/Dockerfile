# Start with the official Golang image
FROM golang:1.23.6-alpine

# Set the Current Working Directory inside the container
WORKDIR /federator

# Install make tools
RUN apk add --no-cache make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN rm .env*

# Build the Go app
RUN make build

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD  ["make", "up"]