# Design: Add golangci-lint

## Context

The gateway's `make lint` runs only `staticcheck`; `secret_assistant` has no lint step.
`golangci-lint` is installed on this machine but fails while downloading modules in this
dev environment (a network limitation), whereas GitHub Actions downloads modules fine.
See proposal.md — Why.

## Goals / Non-Goals

**Goals:**
- One shared, conservative linter configuration for both Go modules.
- `make lint` runs the richer set; the bare `staticcheck` check stays available.

**Non-Goals:**
- No CI workflow edit (that is `ci-pipeline`, which invokes `make lint`).
- No formatting linters (gofmt/goimports stay in the existing `make fmt`/`gofmt` path).
- No production-code behaviour change.

## Decisions

- **Single root `.golangci.yml`.** golangci-lint walks up to the nearest config, so one
  file at the repo root serves both `gateway` and `secret_assistant`. The config lives in
  `backend`'s zone because it governs Go code quality.
- **Conservative linter set** — `staticcheck`, `govet`, `errcheck`, `ineffassign`,
  `unused`, `misspell` — to surface real defects without a noisy first run. Formatting is
  deliberately excluded (Prettier/gofmt handle it elsewhere).
- **`make lint` becomes `golangci-lint run`**, and a new `lint-staticcheck` target keeps
  the old `staticcheck ./...` available for parity and any per-tool triage.
- **`secret_assistant` gets a lint entry point** (a small `Makefile` `lint` target or a
  documented `golangci-lint run` command in its README).
- **Findings are fixed in this change**, or a justified `//nolint` is added with a reason;
  no disabling of the linter as a whole to get a green result.

## Risks / Trade-offs

- [golangci-lint cannot run locally here due to module download] → mitigated: it runs in
  CI (network available) and the version is pinned there by `ci-pipeline`; local runs work
  once module download succeeds.
- [First run may surface many findings] → mitigated by a conservative linter set and fixing
  in one change rather than accumulating exceptions.
