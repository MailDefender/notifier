# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY --exclude=go.* . .

# Build
WORKDIR /app/cmd/notifier
RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier



FROM alpine:latest

# Install CA certificates for SMTP TLS connections
RUN apk add --no-cache ca-certificates

COPY --from=build /notifier /notifier
WORKDIR /

EXPOSE 8080

# Run the notifier service
CMD ["/notifier"]