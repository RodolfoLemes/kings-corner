FROM golang:1.19-buster as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Move to working directory /build
RUN mkdir -p /usr/src/build
WORKDIR /usr/src/build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main cmd/server/server.go

# This is the real Docker image that will be created in the end. It just carefully copies code from the intermediate.
# Note that the secrets from the intermediate will NOT be copied and published.
FROM debian:buster-slim

RUN apt-get update && \
  apt-get install -y ca-certificates libssl-dev libpq-dev

# Source code should be in the /usr/src/app folder
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY --from=builder /usr/src/build/main /usr/src/app/main

EXPOSE 50051
CMD [ "./main" ] 