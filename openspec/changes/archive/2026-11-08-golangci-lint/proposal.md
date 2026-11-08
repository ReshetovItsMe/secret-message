# Proposal: Add golangci-lint

## Why

Issue #16. The gateway's `make lint` runs only `staticcheck`; `secret_assistant` has no
linter wired at all. `golangci-lint` gives a single, configurable runner over a larger
set of linters (including staticcheck, errcheck, govet, ineffassign, unused) and is the
standard Go CI linter. The installed v1.64.8 has a Go 1.27 export-data
incompatibility; network access alone does not resolve analyzer incompatibility.
Use pinned golangci-lint v2.15.0 and prove it can analyze both current modules
before wiring make lint or claiming CI readiness.

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
- Minimal CI prerequisite is part of this change: devops installs pinned v2.15.0
  before gateway make lint changes and preserves staticcheck for lint-staticcheck.
  Crypto lint is verified locally; no unrelated workflow redesign.

## Capabilities

### New Capabilities

_None._ This is developer tooling; there is no product-behaviour change, so the change
declares `skip_specs: true`.

## Impact

- `.golangci.yml` (new, repo root), `gateway/Makefile`, `secret_assistant` lint wiring,
  and any Go fixes the linter requires.

## Sequencing

Starts from current main after crypto-tests, secure-key-storage and ci-pipeline.
The historical PrivateKey/json.Unmarshal SA4006 finding is obsolete after the
envelope-encryption rewrite; rerun baseline instead of applying that old fix.
Serialize existing cards: configuration/tool compatibility → CI prerequisite →
gateway → crypto service → verification. Issue-template live acceptance is independent.
