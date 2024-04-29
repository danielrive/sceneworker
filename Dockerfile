# Use a Golang base image for the build stage
FROM golang:1.22-alpine AS build

# Copy the Go modules files and download dependencies

WORKDIR /app

COPY go.mod .

COPY go.sum .


# Copy the application source code to the container

COPY  devoteam-load-generator/internal/worker/* ./devoteam-load-generator/internal/worker/
COPY  devoteam-load-generator/cmd/load-generator.go ./devoteam-load-generator/cmd/


# Build the Go application statically

WORKDIR /app/load-generator
# Build the Go application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/load-generator load-generator.go

# Use a minimal base image for the runtime stage
FROM gcr.io/distroless/base

# Copy the binary from the build stage
COPY --from=build  /bin/load-generator /bin/load-generator

# Set the command to run the application when the container starts
CMD ["/bin/eload-generator"]
