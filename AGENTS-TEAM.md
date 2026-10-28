# Multi-agent team protocol

This file is the contract for **multi-agent work on this repository**. Any agent that
participates in the team (Hermes profiles `teamlead`, `backend`, `frontend`, `qa`,
`devops`) MUST read this file together with `AGENTS.md` before touching code.
Agents that are not part of the team (single-shot assistants, editors) may ignore it.

`AGENTS.md` describes the system. This file describes **who may change what, and how
work travels from idea to `main`**.

## 1. Roles

| Profile | Role | Owns |
|---|---|---|
| `teamlead` | team lead, architect, spec author, reviewer, integrator | `AGENTS-TEAM.md`, `openspec/**`, `docs/**`, kanban board |
| `backend` | Go services, API and gRPC contracts | `gateway/**`, `secret_assistant/**`, `proto/**` |
| `frontend` | Vue 3 SPA | `web/**` except `web/Dockerfile` and `web/nginx.conf` |
| `qa` | test strategy, e2e, regression, bug reports | `qa/**`, `web/tests/**`, `gateway/internal/integration/**` |
| `devops` | build, CI, containers, proxy, observability | `.github/**`, `docker-compose*.yml`, `**/Dockerfile`, `web/nginx.conf`, root `Makefile` |

`teamlead` does not write product code. It specifies, decomposes, reviews and merges.
The remaining four roles do not merge and do not review each other's work.

Unit tests live next to the code and belong to the code's owner. Integration and e2e
tests belong to `qa`.

## 2. Path ownership is hard

- Edit files only inside your own zone. "It's a one-character fix" is not an exception.
- If your task requires a change outside your zone, create a kanban card for the owning
  profile with `parents=[<your card>]` and state exactly what is needed. Do not edit it
  yourself and do not ask another agent in chat.
- If you discover a file that several open cards are touching, leave a
  `hotspot: <path> — <reason>` comment on your card instead of adding more changes to it.
- Never edit generated code by hand. `gateway/internal/messagepb/**` and
  `secret_assistant/proto/**` are produced by `make proto`.

## 3. Contract ownership

- `proto/message.proto`, the HTTP endpoints (`POST /message`, `GET /message?messageId=`,
  `/healthz`), the error envelope and the cross-service JSON keys (`privateKey`,
  `encryptedKey`, `data`) are owned by `backend`.
- A contract change is only allowed when a `teamlead` decision is recorded in the change
  spec and in the card body. `frontend` and `qa` consume the contract and never change it.
- Changing a contract requires updating both sides, the generated code and the
  **API Contract** section of `AGENTS.md` in the same change.
- No silent field renames, no extra response fields invented on the client side.

## 4. Specs (OpenSpec)

`openspec/` is the spec source of truth and is **written only by `teamlead`**. It uses
schema `spec-driven`: every change lives in `openspec/changes/<change>/` with `proposal.md`,
`design.md` (optional), spec deltas and `tasks.md`.

Work enters the team in one of two ways:

1. **Pre-authored change.** The product owner has already created a change. `teamlead`
   runs `openspec list --json`, `openspec status --change <name> --json` and
   `openspec instructions apply --change <name> --json`, checks the artifacts, fills small
   gaps (or blocks the card for product decisions), then decomposes `tasks.md` into cards.
2. **Global idea.** The product owner describes a goal. `teamlead` researches the
   codebase, then creates one or more changes (`openspec new change <kebab-name>` plus the
   schema's artifacts), validates them, commits them in a spec pull request, and only then
   decomposes them into cards. Independent workstreams get their own change; tightly
   coupled work stays in one change with role-scoped task groups.

Rules:

- Every implementation card names its change (`OpenSpec change: <name>`) and the exact
  `tasks.md` lines it closes.
- Implementers read `openspec/changes/<change>/**`; they never edit it. Spec and code
  disagreeing is a `kanban_block`, not a local spec edit.
- `teamlead` keeps `tasks.md` in sync: after a PR is merged it ticks the completed items
  (`- [ ]` → `- [x]`). A tick is a claim that the work is merged and tested.
- A change is closed with `openspec archive` (plus spec sync) in its own
  `docs(specs): archive <change>` pull request, and only when every task is merged and the
  test suite is green.

## 5. Work flow

```
product owner idea  ──┐
pre-authored change ──┤
                      ▼
teamlead: research → OpenSpec change → spec PR → epic card + child cards
                                                    (assignee + parents)
                      │
                      ▼
dispatcher spawns a worker per card, in an isolated git worktree
                      │
                      ▼
worker: implement → test → commit → push → PR to main → request review
                      │
                      ▼
teamlead: review → request changes (loop) or squash-merge → tick tasks.md → complete card
                      │
                      ▼
qa card (when the epic has one): verify acceptance criteria → bug cards
                      │
                      ▼
teamlead: verify change → archive change → report to product owner
```

The kanban board is the **only** coordination channel. There are no agent-to-agent
chats, no decisions in terminal output, no "I'll tell the other agent later". Every
decision that a downstream card depends on must be written into that card's body —
workers cannot see sibling cards.

## 6. Branches, commits, pull requests

- Base branch is always `main`. Work happens in a linked worktree on branch
  `wt/<task-id>`, created by the dispatcher from the board's project clone.
- One card = one branch = one pull request. Do not bundle unrelated cleanup.
- Commit messages: Conventional Commits, English, imperative mood
  (`feat(gateway): reject oversize payloads`, not `fixed stuff`).
- Pull requests target `main`, title is `<type>(<scope>): <summary>`, body states what
  changed, why, how to verify it, and the kanban card id.
- Never commit to `main` directly, never force-push a shared branch, never rewrite
  published history. Undo mistakes with a revert PR.
- No secrets, tokens, credentials, `.env` values or debug logging in commits.

## 7. Review and merge

- The implementer finishes with `kanban_request_review(reviewer="teamlead")` and stops.
- The `teamlead` reviews the diff against the change spec and the card's acceptance
  criteria: scope, ownership boundaries, contract stability, tests, commit hygiene, no
  secrets.
- Approved → `gh pr merge <N> --squash --subject "<type>(<scope>): <summary>" --delete-branch`,
  then `tasks.md` is ticked and the kanban card is completed.
- Rejected → `kanban_request_changes` plus a PR review comment listing concrete
  `file:line → what to change`. Vague feedback ("improve quality") is not a review.
- Two review cycles without progress, or a product-level disagreement → the card is
  blocked and escalated to the product owner.
- `main` must stay green: the merge is rejected if the relevant test suite fails.

## 8. Definition of done

- `backend`: code + unit tests, `cd gateway && make test` and `make lint` green (or the
  equivalent for `secret_assistant`), generated code up to date, contract documented.
- `frontend`: `pnpm type-check`, `pnpm build` and `pnpm lint` green, states for loading /
  error / consumed-secret covered, no leftover logs or dead code.
- `qa`: every acceptance criterion from the change marked verified with the command or
  scenario used; failures have a reproducing card assigned to the owning profile.
- `devops`: `docker compose config` valid, stack boots from a clean clone with one
  command, CI run linked and green, no secrets in the repo.
- `teamlead`: `tasks.md` fully ticked and matching reality, change archived and specs
  synced, PRs merged, `main` green, product owner has a report listing what shipped and
  what is still open.

## 9. Escalation

- Missing credentials, environment, or a product/UX decision → `kanban_block(reason=...)`
  with the context already written into a comment. Only the product owner unblocks.
- Headless workers cannot ask questions interactively; a blocked card is the only
  correct way to surface a decision that cannot be inferred.
