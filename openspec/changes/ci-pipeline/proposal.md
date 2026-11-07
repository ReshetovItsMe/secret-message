# Proposal: CI pipeline (build, test, static analysis)

## Why

Issue #25. The only workflow, `.github/workflows/build.yml`, is named "Build & Deploy" and
ends with a `push_to_docker_hub` job that logs into Docker Hub and pushes three images —
a deployment step for a project that has **no deployment target**. It also runs only on
push to `main`/`chore/build-ci` (not on pull requests), and its jobs skip the web
type-check and the gateway's `staticcheck` lint that `AGENTS.md` mandates.

The product decision is explicit: **CD is out of scope — no server, no registry push, no
Kubernetes, no ArgoCD.** We keep only CI: build, test and static analysis for every
service, on push and on pull requests.

## What Changes

- Rewrite `.github/workflows/build.yml` as a CI-only pipeline:
  - drop the `push_to_docker_hub` job (registry login + push) and the Docker Hub secrets
    it relied on;
  - rename the workflow to `CI` and trigger on `push` to `main` **and** on `pull_request`;
  - frontend: `pnpm install --frozen-lockfile`, `pnpm type-check`, `pnpm build`,
    `pnpm lint`;
  - gateway: `go mod download`, `go mod verify`, `go build ./...`, `go vet ./...`,
    `go test ./... -race`, `make lint` (staticcheck);
  - secret_assistant: `go mod download`, `go mod verify`, `go build ./...`, `go test ./...`,
    `go vet ./...`;
  - validate the Compose files with `docker compose config --quiet`;
  - refresh pinned action versions to current majors (`checkout@v4`, `setup-go@v5`,
    `setup-node@v4`, `pnpm/action-setup@v4`).
- No production code change, no contract change.

## Capabilities

### New Capabilities

_None._ This is CI/tooling only; there is no product-behaviour change, so the change
declares `skip_specs: true`.

## Impact

- `.github/workflows/build.yml` (rewritten). No `docker-compose*.yml` content change —
  only the new Compose-validation step references them.

## Sequencing

Coordinates with `golangci-lint`: the gateway's static-analysis step runs `make lint`
(staticcheck today, golangci-lint once that change lands). Runs independently of
`secure-key-storage` and `crypto-tests`.
