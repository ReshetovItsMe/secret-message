# Tasks: Node 26 LTS + pnpm + latest dependencies

Role ownership per `AGENTS-TEAM.md`: `frontend` = `web/**` except Dockerfile/nginx;
`devops` = `.github/**` + `web/Dockerfile`; `teamlead` = `openspec/**` + docs.

## 1. Migrate dependencies and lockfile (frontend)

- [x] 1.1 Rewrite `web/package.json` to the dependency matrix in `design.md`, add
  `engines.node` (`^22.22.2 || ^24.15.0 || >=26`) and `packageManager` (`pnpm@12.16.1`),
  remove `@vue/cli-plugin-pwa` / `register-service-worker` / `@rushstack/eslint-patch`,
  add `@element-plus/icons-vue`.
- [x] 1.2 Generate `web/pnpm-lock.yaml` and delete `web/yarn.lock`.
- [x] 1.3 Verify `cd web && pnpm install --frozen-lockfile` resolves cleanly; confirm
  `pnpm --version` is `12.16.1`.

## 2. Migrate build config, run the TypeScript spike, fix code, go green (frontend)

- [x] 2.1 Replace `.eslintrc.cjs` with `eslint.config.ts` (flat config), drop
  `@rushstack/eslint-patch`, update the `lint` script (ESLint 10 has no `--ext`).
- [x] 2.2 Point `tsconfig.node.json` at `@tsconfig/node26`; re-check the `type-check`
  script for vue-tsc 3.3.13 (drop the 1.x `--composite false` workaround).
- [x] 2.3 Review `vite.config.ts` under Vite 8 / Rolldown; fix the `SecretUrl.vue` and
  `SecretView.vue` issues and delete `web/src/registerServiceWorker.js` (see design.md).
- [x] 2.4 Run the TypeScript 7.0.2 spike; report the outcome (green → land 7.0.2, red →
  stay 6.0.3 + failure mode) in the card result.
- [x] 2.5 `cd web && pnpm build && pnpm lint && pnpm type-check` green.

## 3. Container and CI (devops)

- [x] 3.1 Update `web/Dockerfile`: `FROM node:26`, `RUN npm i -g pnpm@12.16.1`,
  `pnpm install --frozen-lockfile`, `pnpm build`.
- [x] 3.2 Update the `frontend` job in `.github/workflows/build.yml` to Node 26 + pnpm
  (see design.md).
- [x] 3.3 `docker compose build frontend` succeeds from a clean clone.

## 4. Verify acceptance criteria (qa)

- [ ] 4.1 `cd web && pnpm install --frozen-lockfile && pnpm build && pnpm lint &&
  pnpm type-check` green on Node 26 (via Docker `node:26` image).
- [ ] 4.2 `docker compose build frontend` green; builder stage is `node:26`.
- [ ] 4.3 CI frontend job green on Node 26 + pnpm.
- [ ] 4.4 No `yarn` references remain in `web/`, `.github/`, or docs.
- [ ] 4.5 No behavioural or visual regression vs. the pre-migration app (smoke: create →
  open secret once → consumed 404).

## 5. Documentation (teamlead — done in the planning PR)

- [x] 5.1 Update `AGENTS.md` (Tech Stack, Common Commands, Notes).
- [x] 5.2 Update `AGENTS-TEAM.md` section 8 frontend DoD (pnpm, Node 26).
- [x] 5.3 Record spike results in `design.md` once the frontend reports (see task 2.4).

> Note: tasks 5.1–5.2 are completed in the planning PR. Task 5.3 is completed by teamlead
> at review time, after the frontend card reports the spike outcome.
