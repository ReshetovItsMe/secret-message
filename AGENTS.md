# Secret Message

Use this file as the mandatory checklist for every task. Preserve the service’s purpose: share a secret through a one-time link and atomically delete its stored payload on retrieval. Keep this file in English.

## 1. Context & Scope

### Stack and source map

Check versions against `go.mod`, `web/package.json`, the lockfile, and CI; do not upgrade dependencies as part of an unrelated task.

| Component | Versions and tooling | Responsibility |
|---|---|---|
| `gateway/` | Go **1.27.0**, chi **5.3.2**, go-redis **9.22.0**, grpc-go **1.83.2**, validator **10.30.3**, lo **1.53.0** | Start at `cmd/gateway/main.go`; find business orchestration in `internal/service/`, Redis in `internal/store/`, the gRPC adapter in `internal/cryptoclient/`, and HTTP in `internal/handler/` and `internal/server/`. |
| `secret_assistant/` | Go **1.27.0**, grpc-go **1.48.0**, stdlib `crypto` | Start at `cmd/main.go`; find cryptography in `internal/encrypt/` and `internal/decrypt/`. |
| `web/` | Vue **3.5.44**, TypeScript **6.0.3**, Vite **8.3.5**, vue-router **5.4.0**, Element Plus **2.14.7**, VueUse **15.0.0**, axios **1.21.0**, vite-plugin-pwa **2.0.0** | Start at `src/main.ts`; change UI in `src/views/` and `src/components/`, and PWA registration in `src/registerServiceWorker.ts`. |
| Infrastructure | Docker Compose, nginx, `redis:7-bullseye`; Node **26** in CI, pnpm **12.16.1** | Find configuration in `docker-compose*.yml`, `**/Dockerfile`, `web/nginx.conf`, and `.github/workflows/build.yml`. |

Account for the Go version and module-path split: use `github.com/ReshetovItsMe/secret-message/gateway` for the gateway and the legacy `github.com/ReshetovItsMe/one-time-messaging-exchange-be` for the crypto service. For web, respect `engines.node`: `^22.22.2 || ^24.15.0 || >=26`; prefer Node 26 to match CI.

### Access boundaries

- Never read, modify, or publish real `.env`/`.env.*` files, private keys (`*.key`, private `*.pem`), secrets/credentials directories, Redis dumps, or real messages without separate explicit permission. Use only known-sanitized examples and synthetic test keys.
- Do not read real runtime/access logs by default: they may contain secret links and personal data. Request permission for a specific source when diagnosing an issue; output only the smallest sanitized excerpt. Analyze output from your own tests using synthetic data immediately.
- Do not hand-edit `gateway/internal/messagepb/` or `secret_assistant/proto/`; read relevant excerpts through Serena only to check the contract, and change them only through regeneration.
- Do not manually explore or edit `.git/`, `**/node_modules/`, `**/dist/`, `**/bin/`, `.serena/cache/`, or the internal `.codegraph/` index. Use Git, build tools, Serena, and CodeGraph to manage these artifacts.
- Preserve other people’s uncommitted changes; never run `git reset --hard`, `git checkout --`, broad deletions, or bulk file overwrites.

## 2. Pre-flight Check

Complete these steps in order before the first edit:

1. **Check Serena MCP.** Locate its tools, call `initial_instructions`, and verify the active project and working directory. If Serena is unavailable, ask the user to configure it and stop; do not substitute direct file reads or edits.
2. **Check a working LSP, not just connection status.** For Go, call `lsp.get_symbols_overview` on a small known file. Report initialization errors to the user; do not bypass a failure by reading code with other tools. Do not treat OpenCode’s built-in LSP as a replacement for Serena’s LSP. Obtain approval before setup/restarts; do not terminate every matching process.
3. **Restore context.** Use `mem.list_memories` to check `.serena/memories/`; read `mem:architecture` for system structure, `mem:domain_model` for terminology, `mem:data_flow` for the lifecycle, and `mem:development_guide` for commands. Read `mem:memory_maintenance` before writing memories. Verify notes against current code; investigate missing information rather than inventing it.
4. **Establish task rules and scope.** Read nested `AGENTS.md` files through Serena (for the gateway, read `gateway/AGENTS.md`). Agree on expected behavior and acceptance criteria; do not turn incidental findings into an unapproved refactor.
5. **Determine team mode.** If participating in the multi-agent team (`teamlead`, `backend`, `frontend`, `qa`, `devops`), read `AGENTS-TEAM.md` through Serena and follow its ownership, OpenSpec, kanban, and review/merge protocol. As a one-off assistant, do not assume a team role or spawn agents without an explicit request or applicable instructions.
6. **Check Git.** Record the current branch and `git status --short`; inspect the task-related diff. Work on a separate branch for code changes; do not commit directly to `main`. In team mode, use `feature/<change>` and a `wt/<task-id>` worktree under `AGENTS-TEAM.md`; agree on the base branch for solo work. Do not switch a dirty tree containing someone else’s work.
7. **Check the baseline.** Before changing code, run the affected module’s checks from section 4 and record commands, exit codes, and pre-existing failures. If clean `main` needs verification, use a separate approved worktree; never reset current changes. For documentation-only changes, check references, structure, and the diff instead of running the entire stack.
8. **Check the environment and specification.** Verify Go/Node/pnpm, linters, and dependencies; distinguish environment failures from code defects. For spec-driven changes, read the relevant `openspec/changes/<change>/` and acceptance criteria in `tasks.md`. Keep the OpenSpec plugin/skill at agent level (https://github.com/Fission-AI/OpenSpec); do not copy agent-specific commands into the repository.

## 3. Architecture & Patterns

### Topology and lifecycle

Preserve the four components and HTTP/gRPC/storage boundaries:

```text
Browser → web (Vue + nginx, :8080)
             │ /api/* → /*
             ▼
          gateway (Go + chi, :3000) ── gRPC ──▶ secret-assistant (:50052)
             │                                  RSA-OAEP + AES-GCM
             ▼
          Redis (:6379)

gRPC contract: proto/message.proto
Dev UI: Vite → /api → localhost:3000 (strip the /api prefix)
```

- Preserve the write path: `handler.Create` → `service.CreateMessage` → gRPC `Encrypt` → `store.Send` → one Redis key `msg:<uuid>` with a TTL (default 24h) → `201 { "messageId": "<uuid>" }`.
- Preserve the read path: `handler.Get` → `service.GetMessage` → `store.Consume` (`GETDEL`) → gRPC `Decrypt` → `200 { "message": "..." }`. Return `404` for a missing/expired/consumed id and `400` when `messageId` is omitted.
- Describe encryption accurately: the gateway passes plaintext to the crypto service, which uses envelope encryption — a per-message AES-256-GCM data key wrapped by a single service-held RSA-2048/OAEP(SHA-256) key. `Encrypt` returns an opaque, versioned binary blob (no key material, no JSON), and the service private key never leaves the crypto service or reaches Redis, so a Redis compromise alone reveals no plaintext. **Do not call this end-to-end encryption** (the gateway and crypto service sit in the plaintext path), and do not claim insecure gRPC provides transport encryption.
- Treat this as **at-most-once retrieval**, not guaranteed delivery: after `GETDEL`, decryption or response-delivery failure may irreversibly lose the secret. Change this model only after explicit architectural agreement.

### Module boundaries and anti-patterns

- Keep HTTP handlers thin: decode/validate requests and shape responses; put orchestration in `service`, Redis in `store`, and protobuf/gRPC in `cryptoclient`. Wire dependencies in `cmd/gateway/main.go`.
- Preserve the deep `Store` interface with `Send`/`Consume`; keep `msg:<uuid>` keys and Redis commands out of handlers/services/UI. Never replace `GETDEL` with separate `GET` + `DEL` operations or move decryption before deletion as a quick fix.
- Treat `proto/message.proto` as the source of truth. The gRPC contract carries an opaque `bytes blob` — `EncryptRequest{plaintext}` → `EncryptResponse{blob}`, `DecryptRequest{blob}` → `DecryptResponse{plaintext}` — with no JSON inside the bytes and no key field. Coordinate HTTP/gRPC contract changes; update both sides, both codegen outputs (`gateway/internal/messagepb/` and `secret_assistant/proto/`), and documentation in one change.
- Validate `message` as a required string of at most 100000 characters; reject unknown JSON fields and trailing data, and preserve the 1 MiB HTTP body limit. Do not weaken validation to make a test pass.
- Preserve `GET /healthz` → `200 ok` and the `{ "error": { "code", "message" } }` envelope (`BAD_REQUEST`, `VALIDATION_ERROR`, `NOT_FOUND`, `INTERNAL`). Return generic internal-error messages rather than Redis/gRPC error strings; do not log plaintext, keys, payloads, or secret links.
- Propagate `context.Context`, preserve the 10s HTTP timeout, wrap errors with `%w`, compare them with `errors.Is/As`, and use `log/slog`. Apply `lo.Must` only when failure is provably impossible; handle JSON/crypto errors and corrupted data explicitly.
- Preserve SPA routes `/`, `/secret/:id`, the catch-all redirect, axios `baseURL='/api'`, and matching `/api` stripping in nginx and Vite. Never cache API responses or secrets in the PWA: keep `NetworkOnly` for `/api/*`; do not enable automatic retries of the consuming GET.
- Do not expose insecure gRPC or Redis to untrusted networks; do not treat Compose networking as a substitute for TLS/auth. Agree on new dependencies and explain why the stdlib and existing set are insufficient.

### Configuration

Take defaults from `gateway/internal/config/config.go`: `HOST=0.0.0.0`, `PORT=3000`, `REDIS_HOST=redis`, `REDIS_PORT=6379`, `REDIS_DB=0`, `SECRET_ASSISTANT_URL=secret-assistant:50052`, `MESSAGE_TTL=24h`, `LOG_PRETTY=false`. Supply `REDIS_USERNAME`/`REDIS_PASSWORD` only through the environment; never persist real values. Use `-port` (50052) for the crypto service; supply its RSA private key via `RSA_PRIVATE_KEY_PATH` (a PEM file, never committed). Do not rely on Redis persistence that Compose does not configure.

## 4. Implementation & Tooling

### Exploration and targeted edits through Serena

1. Use Serena for every file investigation and edit. In the REPL, call the `s.fs`, `s.lsp`, `s.edit`, `s.mem`, and `s.shell` facades with the current session from `initial_instructions`; the names below refer to Serena APIs, not new shell commands.
2. Start with CodeGraph when an index exists, following the rule below; then use `fs.list_dir`, `fs.find_file`, and `fs.search_for_pattern` for localization and `lsp.get_symbols_overview` for structure. Do not load an entire tree or module.
3. Locate the target symbol using `lsp.find_symbol` with `relative_path`; request `include_body=True` only for needed functions. Before changing an interface, inspect `lsp.find_referencing_symbols` and, when needed, `lsp.find_implementations`/`lsp.find_declaration`.
4. To replace a whole symbol, use `edit.replace_symbol_body`; to add a neighbor, use `edit.insert_after_symbol`/`edit.insert_before_symbol`; to rename, use `lsp.rename_symbol`; to delete, use `lsp.safe_delete_symbol` after checking consumers.
5. For a few lines or textual configuration, use `edit.replace_content` with a narrow unique match; for identical edits across files, preview `edit.replace_in_files` in `dry_run` first. Do not rewrite an entire file for a small edit or disturb neighboring code.
6. For Markdown/Vue/configuration without a supported LSP, read only necessary ranges with `fs.read_file` and edit with `edit.replace_content`; create new files with `fs.create_text_file`. Manage memories through `mem.read_memory`/`mem.write_memory`/`mem.edit_memory`.
7. Do not bypass Serena with direct Read/Write/Edit tools, `cat`, `sed`, grep/find, or arbitrary file-reading/rewriting scripts. Reserve shell commands for Git, codegen, approved dependency installation, builds, tests, linters, and service startup; prefer Serena’s `shell.execute_shell_command` for terminating commands. If Serena is unavailable, stop and request configuration.
8. Check diff scope and verify behavior with tests; do not reread an entire file merely to confirm a successful symbolic edit. Do not run formatters that overwrite other people’s files without first checking their scope.

<!-- CODEGRAPH_START -->
### CodeGraph

If `.codegraph/` exists at the root, consult CodeGraph **before searching or reading code** when understanding the system or locating an implementation:

- Use MCP `codegraph_explore` when available; name a file/symbol to obtain current source and call paths, including dynamic dispatch. Query a deferred symbol separately by name.
- If the MCP tool is unavailable, run `codegraph explore "<symbol names or question>"` through Serena’s shell. Do not read the index manually; perform targeted edits through Serena.
- If `.codegraph/` is absent, skip CodeGraph entirely: leave indexing to the user. Do not treat stale index relationships as proof; verify relevant symbols through Serena.
<!-- CODEGRAPH_END -->

### Checks and local startup

Run commands in the specified directory and check their actual exit codes. Enable `pipefail` when piping output; never report successful `tail`/`echo` execution as a successful build or test run.

| Area | Run before review |
|---|---|
| `gateway/` | Format changed Go files; run `make build`, `make test` (`go test ./... -race -count=1 -cover`), `make vet`, and `make lint` (`staticcheck ./...`). Use `make tidy` only for dependency changes and inspect the diff. |
| `secret_assistant/` | Format changed Go files; run `go build ./...`, `go test ./...`, and `go vet ./...`; additionally run tests with `-race` for concurrency changes. Do not treat `[no test files]` as evidence of test coverage. |
| `web/` | Run `pnpm install --frozen-lockfile`, then `pnpm type-check`, `pnpm build`, and `pnpm lint`; account for `lint` using `--fix`. Check changed files with Prettier separately: ESLint skips formatting through `skipFormatting`. |
| Compose/CI/proxy | Run `docker compose config --quiet`; validate production configuration with `docker compose -f docker-compose.prod.yml config --quiet`. For runtime changes, verify builds and startup of affected containers; for workflow changes, validate YAML and the actual job commands. |
| Documentation only | Check Markdown structure, local references, consistency of commands with configuration, and `git diff --check`; state that runtime tests were not run when code did not change. |

Follow existing test patterns: use testify `assert`/`require`, handwritten service fakes, `httptest` with the real HTTP router, and miniredis for storage; do not require external Redis for unit tests. Cover repeat and concurrent reads, TTL expiry, dependency failures, and corrupted payloads when changing those paths. Remember that the gateway integration test uses fake crypto and does not verify the actual RSA/AES implementation. Add a regression test before fixing a bug; do not substitute builds for tests, especially in web and the crypto service, which had no tests at onboarding.

Follow Go formatting and stdlib-first conventions. For Vue, use `<script setup lang="ts">`, 4-space indentation, single quotes, and `;` according to Prettier; use pnpm and preserve `web/pnpm-lock.yaml`. Check CI against `.github/workflows/build.yml`: the `CI` workflow runs on `push` to `main` and on every `pull_request`, and includes the gateway `make lint` (staticcheck) and web `pnpm type-check` steps.

For the complete local stack, run `docker compose up --build`. For separate services, start Redis and run `go run ./cmd -port 50052` in `secret_assistant/`, then run `REDIS_HOST=localhost SECRET_ASSISTANT_URL=localhost:50052 make run` in `gateway/` and `pnpm dev` in `web/`. For proto changes, run `make proto` in the gateway and the separate source-relative command from `secret_assistant/README.md`; verify both generated outputs.

### Context mode

If the global `context-mode` MCP is available, use it to reduce context usage, not as a Serena replacement. For bulky results, use available `ctx_execute`/`ctx_execute_file` and `ctx_batch_execute` tools (I/O concurrency 1–8), and output only conclusions. Use `ctx_search(queries: [...])` for already indexed information; use `ctx_fetch_and_index` followed by search for web sources. After a break, recover decisions with `ctx_search(sort: "timeline")` when available. Check availability with `ctx stats`/`ctx doctor`; if this optional server is absent, continue through Serena rather than calling nonexistent tools.

## 5. Self-Healing & Definition of Done

### When a check fails

1. Read the complete output of your check, record the command and exit code, and isolate the first root cause rather than the final cascading symptom. Compare against the baseline: distinguish pre-existing defects, regressions, incompatible tool versions, and missing dependencies.
2. Reproduce the failure with a minimal test; inspect relevant symbols/calls through Serena. For runtime logs, follow section 1’s restrictions; request a sanitized excerpt rather than reading an entire log directory.
3. Form one testable hypothesis, make the smallest fix through Serena, and rerun the targeted test. Do not delete tests, weaken assertions/validation, or disable a linter to get a green result.
4. After the fix, rerun the complete check set for affected modules. Do not replace broken pnpm with a different package manager or blindly update the lockfile; agree on environment repairs. If a service, secret, or product decision is unavailable, stop with a precise blocker.
5. After two failed hypotheses, revisit assumptions; after three failed attempts, stop speculative edits and discuss the cause with the human. In team mode, record the blocker through the `AGENTS-TEAM.md` protocol.

### Final checklist

- [ ] Check the result against the task, OpenSpec, and acceptance criteria; preserve architectural boundaries and contracts, and state agreed deviations.
- [ ] Inspect the entire diff: remove your own debris, debug output, and accidental formatting; preserve others’ changes. Run `git diff --check` and inspect `git status --short`.
- [ ] Verify regression tests, builds, type checks, formatting, and linters with fresh commands; provide results. Explicitly list skipped/blocked checks; do not claim everything passes without evidence.
- [ ] For contract changes, verify both sides, both codegen outputs, and the API description; exclude real secrets, keys, `.env` files, and runtime logs from the diff, memories, and report.
- [ ] Update relevant `.serena/memories/` through Serena when durable architecture, flows, or startup rules change. Follow `mem:memory_maintenance`; do not record temporary environment failures as permanent project properties. Preserve `mem:...` references and validate them with `serena memories check` when the CLI is available.
- [ ] Update documentation and OpenSpec artifacts within your role. In team mode, hand changes outside your ownership to the owner instead of editing them yourself; do not mark `tasks.md` complete before the team’s prescribed stage.
- [ ] Create commits/PRs only when requested or required by the applicable team protocol. Select only task-owned files, use English imperative Conventional Commits, and do not force-push or rewrite shared history.
- [ ] Give the human a concise report: what changed, which checks actually ran, and remaining risks/blockers. In team mode, open a PR targeting `feature/<change>`, request review from `teamlead`, and stop; do not merge on behalf of another role.
