# Design: Add golangci-lint

## Context

The gateway's `make lint` runs only `staticcheck`; `secret_assistant` has no lint step.
The installed golangci-lint v1.64.8 cannot read Go 1.27 export data.
The current CI installs staticcheck only, so changing make lint without a prior
CI tool-install step would fail. See proposal.md — Why.

## Goals / Non-Goals

**Goals:**
- One shared, conservative linter configuration for both Go modules.
- `make lint` runs the richer set; the bare `staticcheck` check stays available.

**Non-Goals:**
- No unrelated CI redesign; only pinned tool installation as a prerequisite.
- No formatting linters (gofmt/goimports stay in the existing `make fmt`/`gofmt` path).
- No production-code behaviour change.

## Decisions

- **Pinned golangci-lint v2.15.0, v2 config (`version: "2"`).** Release and module
  dependencies were inspected upstream. Backend must verify version output, config
  validation, and successful analyzer initialization on both modules with Go 1.27.
  An actual finding is not an initialization failure. If incompatible, stop and
  report exact errors rather than assuming CI fixes it. Do not overwrite the old
  global binary; use an isolated tool install and explicit PATH.
- **Fresh baseline.** Current-main gateway build/race tests/vet/staticcheck and
  crypto build/tests/vet/staticcheck pass. Do not reproduce the obsolete
  PrivateKey/json.Unmarshal patch from the old epic comment.
- **Minimal CI prerequisite owned by devops.** Install exactly v2.15.0 in both Go
  jobs with a compatible Go toolchain (1.27.x) and keep staticcheck available.
  Merge after config and before gateway make lint changes. Crypto lint is
  verified locally; introducing a new crypto CI lint gate is out of scope.
  No environment/tool failure may be accepted as green verification.

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

- [Analyzer/toolchain incompatibility] → require a version-pinned smoke run on both
  current modules before wiring. A download or compatibility failure blocks actual
  golangci-lint acceptance; staticcheck/vet are baseline evidence, not substitutes.
- [First run may surface many findings] → mitigated by a conservative linter set and fixing
  in one change rather than accumulating exceptions.
