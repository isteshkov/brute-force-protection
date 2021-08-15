# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.16-buster as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the command inside the container.
RUN make deps
RUN make build

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN apt-get update && apt-get install -yq ca-certificates
# # Copy the binary to the production image from the builder stage.
COPY --from=builder /app/brute_force_protection /
# # Run the service on container startup.
CMD ["./brute_force_protection"]
