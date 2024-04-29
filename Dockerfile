# Use a Golang base image for the build stage
FROM golang:1.22-alpine AS build

# Copy the Go modules files and download dependencies

WORKDIR /app

COPY go.mod .

COPY go.sum .


# Copy the application source code to the container

COPY  /common/* ./common/
COPY  /utils/* ./utils/
COPY  /internal/worker/* ./internal/worker/
COPY  /cmd/load_generator.go ./load_generator.go


# Build the Go application statically

WORKDIR /app
# Build the Go application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/load_generator load_generator.go

# Use a minimal base image for the runtime stage
FROM gcr.io/distroless/base

# Copy the binary from the build stage
COPY --from=build  /bin/load_generator /bin/load_generator

# Set the command to run the application when the container starts
CMD ["/bin/load_generator"]
