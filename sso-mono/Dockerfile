FROM golang:1.23.5 AS builder
# Set destination for COPY
WORKDIR /usr/src/sso-mono

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify
 
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

COPY . .

# Build
RUN CGO_ENABLED=0 go build -v -o ./bin/app .


FROM alpine:latest

# Set destination for COPY
WORKDIR /usr/src/app

# Copy the binary from the builder stage
COPY --from=builder /usr/src/sso-mono/bin .

# Run the binary
CMD ["./app"]