# Use the golang:1.19 image as the build stage
FROM golang:1.23 AS build

# Copy the source code
COPY . /app

# Set the working directory
WORKDIR /app

# Download and verify module dependencies
RUN go mod download && go mod verify

# Build the sesame-gateway executable
RUN cd /app/cmd && go build -o /app/sesame-gateway
# Expose the necessary port(s)
EXPOSE 8080

# Set the command to run when the container starts
CMD ["/app/sesame-gateway"]
