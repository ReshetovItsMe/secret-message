# Tasks: CI pipeline (build, test, static analysis)

Role ownership: `devops` = `.github/**`, `docker-compose*.yml`, root `Makefile`.

## 1. Rewrite the CI workflow (devops)

- [x] 1.1 Rewrite `.github/workflows/build.yml` as a CI-only workflow named `CI`,
  triggered on `push` to `main` and on `pull_request`.
- [x] 1.2 Remove the `push_to_docker_hub` job (Docker Hub login, QEMU/Buildx, build-push
  steps) and drop any reference to `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`.
- [x] 1.3 `frontend` job (Node 26, pnpm 12.16.1): `pnpm install --frozen-lockfile`,
  `pnpm type-check`, `pnpm build`, `pnpm lint`.
- [x] 1.4 `gateway` job (Go 1.27.x): `go mod download`, `go mod verify`,
  `go build ./...`, `go vet ./...`, `go test ./... -race -count=1 -cover`, `make lint`
  (staticcheck).
- [x] 1.5 `secret_assistant` job (Go >=1.18): `go mod download`, `go mod verify`,
  `go build ./...`, `go test ./...`, `go vet ./...`.
- [x] 1.6 Add a Compose-validation job: `docker compose config --quiet` and
  `docker compose -f docker-compose.prod.yml config --quiet`.
- [x] 1.7 Refresh action versions to `actions/checkout@v4`, `actions/setup-go@v5`,
  `actions/setup-node@v4`, `pnpm/action-setup@v4`.

## 2. Verification (devops)

- [x] 2.1 Validate the workflow YAML and `git diff --check` clean.
- [x] 2.2 Confirm a CI run is green on a PR and on a push to `main`, with every job
  listed above passing and no Docker Hub/registry step present.
