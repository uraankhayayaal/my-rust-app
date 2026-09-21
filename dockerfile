# syntax=docker/dockerfile:1

# Stage 1: Build the Rust application
FROM rust:latest AS builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 2: Dev env (cargo, rustfmt, clippy, tests available)
FROM rust:latest AS dev
WORKDIR /app
COPY . .
RUN rustup component add rustfmt clippy
CMD ["cargo", "run"]

# Stage 3: Minimal runtime image
FROM debian:stable-slim
COPY --from=builder /app/target/release/App /usr/local/bin/App
CMD ["/usr/local/bin/App"]
