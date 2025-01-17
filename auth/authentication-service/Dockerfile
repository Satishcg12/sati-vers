FROM golang:1.23.4 AS builder
# Set destination for COPY
WORKDIR /usr/src/app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
 
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY ./authentication-service ./authentication-service

# Build
RUN CGO_ENABLED=0 go build -v -o ./bin/authentication-service ./authentication-service


FROM alpine:latest

# Set destination for COPY
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /usr/src/app/bin/authentication-service ./authentication-service

# Run

CMD ["./authentication-service"]