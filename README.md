## Build
```bash
docker compose build
```

## Run
```bash
docker compose run --rm app cargo run
```

## Check
```bash
docker compose run --rm app cargo check
```

## Lock cargo deps
```bash
docker compose run --rm app cargo update
```

## Prettify code
```bash
docker compose run --rm app cargo fmt
```

## Analize code
```bash
docker compose run --rm app cargo clippy
```

## Build release app
1. Uncomment publish section at `compose.yaml`
1. Build `docker compose build`
1. Run `docker compose run --rm release`

## Links
- Leaning https://doc.rust-lang.ru/book/
