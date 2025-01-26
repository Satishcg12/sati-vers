FROM golang:1.23.5 AS builder
# Set destination for COPY
WORKDIR /usr/src/sso-mono

RUN go install github.com/air-verse/air@latest

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify
 
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

COPY . .


# Run the binary
CMD ["air", "-c", ".air.toml"]