## Build
```bash
go build -o app .
```

## Run
```bash
./app [easy|hard|normal]
```

## Check
```bash
go vet ./...
```

## Test
```bash
go test ./...
```

## Prettify code
```bash
gofmt -w .
```

## Links
- Go: https://go.dev/doc/

## CI/CD Architecture (ARCH-03)

### Pipeline: GitHub Actions (.github/workflows/ci.yml)

```
push/PR → jobs.test (go test/vet/gofmt)
           → jobs.qa  (go test -v, needs: test)
           → jobs.deploy (kubectl apply -f k8s/, needs: [test, qa], env: production)
```

- **Runner**: ubuntu-latest
- **Toolchain**: actions/setup-go@v5 (Go 1.21)
- **QA gate**: job `qa` blocks merge if any unit test fails
- **Deploy gate**: job `deploy` requires both `test` and `qa` green
- **Secrets**: K8S_CONFIG (base64 kubeconfig), K8S_NAMESPACE (default: my-rust-app)

### Secondary pipeline: GitLab CI (.gitlab-ci.yml)

Остаётся как вторичный пайплайн (build → test → lint через docker compose). Не конфликтует с GitHub Actions.

### Design decisions

- CI выполняет go-команды напрямую на раннере (без docker compose run) — проще, быстрее, KISS.
- Docker Compose остаётся локальным инструментом разработки, не CI-рантаймом.
- Автотесты QA = unit-тесты в main_test.go, прогон через `go test`.
- Деплой в прод только через Kubernetes (k8s/ манифесты), не через Docker Compose.

## План работ

- [x] GO-01: Переписать src/main.rs на Go (main.go + go.mod) — Senior Go Developer
- [x] GO-02: Обновить Dockerfile и CI/CD под Go — Senior Go Developer
