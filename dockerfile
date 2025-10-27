# Stage 1: Build the Rust application
FROM rust:latest AS builder
WORKDIR /app
COPY . .
RUN cargo build --release

# Stage 2: Create a minimal runtime image
FROM debian:stable-slim
COPY --from=builder /app/target/release/hello_world /usr/local/bin/hello_world
CMD ["/usr/local/bin/hello_world"]