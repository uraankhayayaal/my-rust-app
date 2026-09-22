# syntax=docker/dockerfile:1

# Stage 1: Build the Go application
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

# Stage 2: Dev env (go, gofmt, vet, tests available)
FROM golang:latest AS dev
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
CMD ["go", "run", "."]

# Stage 3: Minimal runtime image
FROM debian:stable-slim
COPY --from=builder /app/app /usr/local/bin/app
CMD ["/usr/local/bin/app"]
