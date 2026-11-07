# Proposal: Add golangci-lint

## Why

Issue #16. The gateway's `make lint` runs only `staticcheck`; `secret_assistant` has no
linter wired at all. `golangci-lint` gives a single, configurable runner over a larger
set of linters (including staticcheck, errcheck, govet, ineffassign, unused) and is the
standard Go CI linter. It is already installed on this machine but currently fails while
downloading modules in this dev environment — an environment/network limitation, not a
code defect; GitHub Actions has network access and will run it cleanly.

## What Changes

- Add a root `.golangci.yml` with a conservative linter set (staticcheck, govet, errcheck,
  ineffassign, unused, misspell and a few others; gofmt/goimports left to the existing
  `gofmt`/`make fmt` path).
- Change `gateway/Makefile` `make lint` to run `golangci-lint run`, keeping `staticcheck`
  as a separate `lint-staticcheck` target so the old check remains available.
- Add a `golangci-lint` invocation for `secret_assistant` (a small `Makefile` target or a
  documented command in `secret_assistant/README.md`), since it has no lint step today.
- Fix any findings surfaced in `gateway/**` and `secret_assistant/**`, or document a
  deliberate `//nolint` with a reason where appropriate.
- CI integration is intentionally **not** part of this change: `ci-pipeline` invokes
  `make lint`, which will pick up golangci-lint automatically.

## Capabilities

### New Capabilities

_None._ This is developer tooling; there is no product-behaviour change, so the change
declares `skip_specs: true`.

## Impact

- `.golangci.yml` (new, repo root), `gateway/Makefile`, `secret_assistant` lint wiring,
  and any Go fixes the linter requires.

## Sequencing

Lands before `ci-pipeline` so that the gateway's CI static-analysis step (`make lint`)
runs golangci-lint rather than only staticcheck. Independent of `secure-key-storage` and
`crypto-tests`.
