# Tasks: Add golangci-lint

Role ownership: `backend` = `gateway/**`, `secret_assistant/**`, plus the root
`.golangci.yml` (Go tooling configuration). `devops` owns `.github/**`.
`qa` owns the `gateway/internal/integration/**` findings in task 2.2; `backend`
owns production fixes and Makefile wiring. Complete task 2.2 only after both
QA and backend PRs are merged and the full gateway lint run passes.

All work targets feature/golangci-lint refreshed from current main. Use pinned
golangci-lint v2.15.0 with Go 1.27; the old v1.64.8 is not an acceptable runner.
Do not apply the historical PrivateKey/json.Unmarshal finding without reproduction.

## 1. Configuration (backend)

- [x] 1.1 Add a root `.golangci.yml` enabling a conservative set (at least `staticcheck`,
  `govet`, `errcheck`, `ineffassign`, `unused`, `misspell`); leave formatting to the
  existing `gofmt`/`make fmt` path. Validate v2 config and prove both module
  analyzers initialize with pinned v2.15.0 before wiring (report real findings
  separately from environment failures).

## 1a. CI prerequisite (devops)

- [ ] 1.2 Install pinned golangci-lint v2.15.0 in both Go CI jobs using Go 1.27.x,
  preserve staticcheck availability. Verify CI before
  gateway make lint switches; do not redesign unrelated jobs.

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
