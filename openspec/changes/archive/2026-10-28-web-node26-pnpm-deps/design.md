# Design: Node 26 LTS + pnpm + latest dependencies

## Context

`web/` is a Vue 3 SPA (Vite, vue-router, element-plus, @vueuse/core, axios) built with a
Yarn-1 toolchain on Node 18. Everything here is build/runtime infra; no user-visible
behaviour changes. This document records the decisions the implementation cards depend
on, so workers do not re-derive or re-open them.

## Decisions (fixed by the product owner — do not re-open)

| Decision | Value |
|---|---|
| Node runtime | **26.12.1** (Active LTS). `node:18` in `web/Dockerfile` and CI is EOL. |
| Package manager | **pnpm 12.16.1** (moving off Yarn 1; `yarn` is not installed on this host). |
| TypeScript | **6.0.3** now; **spike 7.0.2** against `vue-tsc@3.3.13` (see Spike section). |
| Vue Router | **5.4.0** (no breaking changes for projects without `unplugin-vue-router`). |
| Vue | **3.5.44** (do NOT use `3.6.0-rc.13`). |

## Dependency matrix

Baseline (`web/package.json`, verified 2026-10-28) → target (npm registry, 2026-10-28):

| Package | From | To |
|---|---|---|
| vue | ^3.3.4 | **3.5.44** |
| vue-router | 4 | **5.4.0** |
| element-plus | ^2.4.0 | **2.14.7** |
| @vueuse/core | ^10.7.2 | **15.0.0** |
| axios | ^1.5.1 | **1.21.0** |
| @element-plus/icons-vue | *(missing, hoisted)* | **2.3.2** |
| vite | ^4.4.9 | **8.3.5** (Rolldown-based: `rolldown ~1.2.14` + `lightningcss`; engines `^20.19 \|\| >=22.12`) |
| @vitejs/plugin-vue | ^4.3.4 | **6.0.10** |
| typescript | ~5.2.0 | **6.0.3** (spike 7.0.2) |
| vue-tsc | ^1.8.11 | **3.3.13** |
| eslint | ^8.49.0 | **10.13.0** (flat config only) |
| eslint-plugin-vue | ^9.17.0 | **10.12.0** |
| @vue/eslint-config-typescript | ^12.0.0 | **14.9.0** |
| @vue/eslint-config-prettier | ^8.0.0 | **10.2.0** |
| @vue/tsconfig | ^0.4.0 | **0.9.1** |
| @tsconfig/node18 | ^18.2.2 | **@tsconfig/node26 26.0.1** |
| @types/node | ^18.17.17 | **26.6.5** |
| npm-run-all2 | ^6.0.6 | **9.0.4** (engines `^22.22.2 \|\| ^24.15.0 \|\| >=26`) |
| prettier | ^3.0.3 | **3.9.10** |

Removed: `@vue/cli-plugin-pwa`, `register-service-worker`, `@rushstack/eslint-patch`.
Added: `@element-plus/icons-vue`.

## Environment: `engines` and pnpm availability (teamlead decision)

The build host runs Node **24.20.0** (`~/.local/bin/node`, no nvm/fnm) and has no `yarn`.
Node 26 is the *runtime target* for CI and the container; it is not installed locally.

Decision — **option (c), relax `engines`** rather than block on a host upgrade:

```json
"engines": { "node": "^22.22.2 || ^24.15.0 || >=26" },
"packageManager": "pnpm@12.16.1"
```

Rationale:

- `engines` declares a *compatibility range*. `^24.15.0` keeps the current dev host
  (24.20.0) able to run `pnpm install/build/lint/type-check`, while CI pins
  `node-version: 26` and the Docker builder is `node:26` — so the product decision
  "runtime = 26.12.1" is honoured exactly where it matters (CI + production image).
- The range is the union of the strictest dependency engine constraints
  (`npm-run-all2@9` = `^22.22.2 || ^24.15.0 || >=26`; `vite@8` = `^20.19 || >=22.12`),
  so no dependency is violated on either Node 24 or Node 26.

pnpm is installed on the host via corepack (`corepack prepare pnpm@12.16.1 --activate`) or
`npm i -g pnpm@12.16.1`; the implementation cards must confirm `pnpm --version` is
`12.16.1` before running. `node:26` images no longer bundle corepack/yarn, so the
Dockerfile installs pnpm explicitly (`RUN npm i -g pnpm@12.16.1`).

## Lockfile migration

- Delete `web/yarn.lock` (Yarn v1, 5263 lines).
- Generate `web/pnpm-lock.yaml` by running `pnpm install` against the new `package.json`
  (then `pnpm install --frozen-lockfile` for reproducibility).
- CI and Docker must reference `pnpm-lock.yaml`, not `yarn.lock`.

## ESLint flat-config migration

- Delete `.eslintrc.cjs` and the `@rushstack/eslint-patch` dependency.
- Create `eslint.config.ts` using `defineConfigWithVueTs` (from
  `@vue/eslint-config-typescript`) + `vueTsConfigs.recommended`, plus
  `@vue/eslint-config-prettier/skip-formatting`, keeping the existing
  `vue/multi-word-component-names: off` rule and the 4-space / prettier formatting.
- ESLint 10 removed `--ext`, so the `lint` script must drop `--ext ...` and rely on the
  flat-config `files` globs (e.g. `eslint . --fix --ignore-path .gitignore`).
- Known friction point to verify explicitly: ESLint 10 flat config + pnpm's strict
  `node_modules` plugin resolution. If a config preset is not resolvable, declare the
  affected plugin as a direct devDependency (pnpm does not hoist transitive plugins).

## TypeScript / vue-tsc and scripts

- `type-check` currently is `vue-tsc --noEmit -p tsconfig.app.json --composite false`.
  `--composite false` is a vue-tsc 1.x workaround; re-check under 3.3.13 and use the
  idiomatic invocation (`vue-tsc --noEmit -p tsconfig.app.json`, or `vue-tsc --build`
  where the project-reference setup warrants it).
- `tsconfig.node.json` extends `@tsconfig/node18/tsconfig.json` → `@tsconfig/node26/tsconfig.json`.

## Spike: TypeScript 7.0.2

Procedure: with the rest of the migration at 6.0.3 green, switch `typescript` to `7.0.2`
and run `pnpm type-check && pnpm build` against `vue-tsc@3.3.13`.

- **Green** → land 7.0.2 in this change.
- **Red** → stay on 6.0.3 and record the failure mode here.

> Result: **RED — stay on 6.0.3.** `vue-tsc@3.3.13` resolves `typescript/lib/tsc`, which
> TypeScript 7.0.2 no longer exports (`ERR_PACKAGE_PATH_NOT_EXPORTED`); independently
> `@typescript-eslint@8.72.1` declares a peer range `typescript >=4.8.4 <6.1.0` that
> rejects 7.x. Landed on 6.0.3.

## Spike: Vite 8 (Rolldown) compatibility

`vite.config.ts` currently only registers `@vitejs/plugin-vue` and a `@` alias — the
minimal config. Verify it loads under Vite 8 / Rolldown with `@vitejs/plugin-vue@6.0.10`
and no further changes; record any plugin/alias incompatibility discovered.

> Result: **GREEN.** The minimal `vite.config.ts` (plugin-vue + `@` alias) loads and
> builds cleanly under Vite 8 / Rolldown with `@vitejs/plugin-vue@6.0.10`; no plugin/alias
> incompatibility found.

## Required code fixes (no behaviour change)

- `web/src/components/SecretUrl.vue`:
  - `@element-plus/icons-vue` is now a declared dependency (see matrix).
  - `:34` `v-model` on `<el-text>` is unsupported → render the URL as text content (no
    `v-model`); the clipboard source is already `textToCopy`.
  - `:39` `type="''"` is an invalid prop value → remove it (default type) or use a valid
    string literal.
- `web/src/views/SecretView.vue`: remove the `async` top-level `await` + the nested
  `<Suspense>` inside the view (App.vue already wraps `<RouterView>` in `<Suspense>`);
  move the fetch into a `ref`/`computed` state and distinguish a 404 (consumed/missing
  secret) from a network/server error so the user sees the right message. No visual
  redesign — keep the current `.secret-container` look.
- Delete `web/src/registerServiceWorker.js` (dead code; reads `process.env.NODE_ENV` /
  `BASE_URL` that do not exist under Vite). PWA is change 3, out of scope here.

## Dockerfile and CI (devops-owned)

- `web/Dockerfile`: `FROM node:26` (builder), `RUN npm i -g pnpm@12.16.1`, `COPY
  package.json pnpm-lock.yaml ./`, `RUN pnpm install --frozen-lockfile`, `RUN pnpm build`.
  Production stage unchanged (nginx).
- `.github/workflows/build.yml` `frontend` job: `actions/setup-node@v4` with
  `node-version: 26`, install pnpm (`pnpm/action-setup@v4` with `version: 12.16.1`),
  `pnpm install --frozen-lockfile`, `pnpm build`, `pnpm lint`.

## Out of scope

Layout/responsive changes (change 2 `web-responsive-ui`), PWA/service-worker work
(change 3 `web-pwa`), `gateway/**`, `secret_assistant/**`, `proto/**`, HTTP contract.

## Risks

- ESLint 10 flat config + pnpm resolution (mitigated above).
- Vite 8 / Rolldown behavioural differences (minimal config → low risk; spike verifies).
- TS 7.0.2 vs vue-tsc 3.3.13 compatibility (spike decides 6.0.3 vs 7.0.2).
