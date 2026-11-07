# Tasks: Add golangci-lint

Role ownership: `backend` = `gateway/**`, `secret_assistant/**`, plus the root
`.golangci.yml` (Go tooling configuration).

## 1. Configuration (backend)

- [ ] 1.1 Add a root `.golangci.yml` enabling a conservative set (at least `staticcheck`,
  `govet`, `errcheck`, `ineffassign`, `unused`, `misspell`); leave formatting to the
  existing `gofmt`/`make fmt` path.

## 2. Wire into the gateway (backend)

- [ ] 2.1 Change `gateway/Makefile` `make lint` to run `golangci-lint run`, and add a
  `lint-staticcheck` target that still runs `staticcheck ./...`.
- [ ] 2.2 Run `cd gateway && make lint` and fix any findings (or add a justified
  `//nolint`).

## 3. Wire into secret_assistant (backend)

- [ ] 3.1 Add a `golangci-lint` invocation for `secret_assistant` (a `Makefile` `lint`
  target or a documented command in `secret_assistant/README.md`).
- [ ] 3.2 Run it and fix any findings.

## 4. Verification (backend)

- [ ] 4.1 `cd gateway && make lint && make lint-staticcheck && make test` green.
- [ ] 4.2 `cd secret_assistant && golangci-lint run && go build ./... && go test ./...`
  green.
