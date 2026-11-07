# Design: CI pipeline (build, test, static analysis)

## Context

`.github/workflows/build.yml` is the sole workflow and does two jobs beyond CI: it logs
into Docker Hub and pushes three images (`push_to_docker_hub`). With no deployment
target, that job is dead weight that also drags two secrets (`DOCKERHUB_USERNAME`,
`DOCKERHUB_TOKEN`) into the repo's CI surface. The workflow only runs on push to `main`
and `chore/build-ci`, so a PR can merge green checks that the actual PR never exercised.
Finally, the job matrix under-runs the checks `AGENTS.md` §4 prescribes: no web
`type-check` step and no gateway `staticcheck`.

## Goals / Non-Goals

**Goals:**
- CI runs on every push to `main` and every pull request.
- Each service is built, tested and statically analysed with the commands `AGENTS.md`
  already mandates.
- No deployment step remains: no Docker Hub login, no registry push, no image tagging.

**Non-Goals:**
- No deployment of any kind (no server, K8s, ArgoCD, registry push) — explicit product
  decision.
- No Docker image build job (building images is fine later if a target appears; it is not
  CI's job today).
- No code changes, no `docker-compose*.yml` edits.

## Decisions

- **Remove `push_to_docker_hub` entirely**, including the QEMU/Buildx/login steps and the
  two Docker Hub secrets it depends on. Nothing else references these secrets.
- **Rename the workflow to `CI`** and set triggers to `push: branches: [main]` plus a
  bare `pull_request:` (all PR targets). This is the single highest-value change: it makes
  PR checks actually run.
- **Per-service commands** mirror `AGENTS.md` §4:
  - `frontend` (Node 26, pnpm 12.16.1): `pnpm install --frozen-lockfile`,
    `pnpm type-check`, `pnpm build`, `pnpm lint`. Note `pnpm build` already runs
    `type-check` (`run-p type-check build-only`), but an explicit `type-check` step fails
    faster and is clearer in the log.
  - `gateway` (Go 1.27.x): add `make lint` (staticcheck) to the existing
    build/vet/test-race steps; align the test step with `make test` semantics
    (`-race -count=1 -cover`).
  - `secret_assistant` (Go >=1.18): change `go build` (in `cmd/`) to `go build ./...`
    from the module root so the whole module compiles; keep `go test ./...` and
    `go vet ./...`.
- **Add a Compose validation job** running `docker compose config --quiet` (and the prod
  file) — `AGENTS.md` lists it under the Compose/CI/proxy checks.
- **Refresh action versions** to current majors while we are already rewriting the file
  (`actions/checkout@v4`, `actions/setup-go@v5`, `actions/setup-node@v4`,
  `pnpm/action-setup@v4`). No other tooling is added.
- **Static-analysis hand-off**: the gateway job's `make lint` runs `staticcheck` today and
  `golangci-lint` after the `golangci-lint` change lands. The CI file does not pin the
  tool, it invokes `make lint`, so the two changes do not conflict.

## Risks / Trade-offs

- Bumping action majors carries a small risk of a new default (e.g. setup-go cache
  behaviour); mitigated by keeping the same `cache`/`cache-dependency-path` settings and
  verifying on the first green run.
- Removing the Docker Hub job means no images are published anywhere; that is exactly the
  intended scope and is reversible later by adding a deploy change.
- The workflow file is YAML and is exercised only in CI; local verification is
  `docker compose config --quiet` plus careful YAML review (`git diff --check`).
