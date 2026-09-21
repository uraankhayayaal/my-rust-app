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

## CI/CD Architecture (ARCH-03)

### Pipeline: GitHub Actions (.github/workflows/ci.yml)

```
push/PR → jobs.test (cargo test/check/fmt/clippy)
           → jobs.qa  (cargo test -- --nocapture, needs: test)
           → jobs.deploy (kubectl apply -f k8s/, needs: [test, qa], env: production)
```

- **Runner**: ubuntu-latest
- **Toolchain**: dtolnay/rust-toolchain@stable (rustfmt, clippy)
- **QA gate**: job `qa` blocks merge if any unit test fails
- **Deploy gate**: job `deploy` requires both `test` and `qa` green
- **Secrets**: K8S_CONFIG (base64 kubeconfig), K8S_NAMESPACE (default: my-rust-app)

### Secondary pipeline: GitLab CI (.gitlab-ci.yml)

Остаётся как вторичный пайплайн (build → test → lint через docker compose). Не конфликтует с GitHub Actions.

### Design decisions

- CI выполняет cargo-команды напрямую на раннере (без docker compose run) — проще, быстрее, KISS.
- Docker Compose остаётся локальным инструментом разработки, не CI-рантаймом.
- Автотесты QA = unit-тесты в src/main.rs (mod tests), прогон через `cargo test`.
- Деплой в прод только через Kubernetes (k8s/ манифесты), не через Docker Compose.
