## Build
```bash
docker compose build
```

## Run
```bash
docker compose run rust-app
```

## Lock cargo deps
1. Build builder
    ```bash
    docker build --target=builder -t rust-app-builder .
    ```
1. Update Cargo
    ```bash
    docker run -v ./:/app rust-app-builder cargo update
    ```