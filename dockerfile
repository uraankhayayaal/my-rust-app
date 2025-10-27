# Stage 1: Build the Rust application
FROM rust:latest AS builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 2: Dev env
FROM rust:latest AS dev
WORKDIR /app
COPY . .
RUN rustup component add rustfmt clippy

# Stage 3: Create a minimal runtime image
FROM debian:stable-slim
COPY --from=builder /app/target/release/App /usr/local/bin/App
CMD ["/usr/local/bin/App"]