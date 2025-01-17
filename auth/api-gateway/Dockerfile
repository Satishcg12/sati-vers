FROM golang:1.23.4 AS build

# Set destination for COPY
WORKDIR /usr/src/app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY ./api-gateway ./api-gateway

# Build
RUN CGO_ENABLED=0 go build -v -o ./bin/api-gateway ./api-gateway

FROM alpine:latest

# Set destination for COPY
WORKDIR /usr/local/bin

# Copy the binary from the build stage
COPY --from=build /usr/src/app/bin/api-gateway ./api-gateway

# Run
CMD ["api-gateway"]