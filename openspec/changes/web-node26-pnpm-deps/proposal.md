# Proposal: Migrate `web/` to Node 26 LTS + pnpm + latest dependencies

## Why

`web/` is pinned to a toolchain that is end-of-life: Node 18 (EOL 2025-04-30), Yarn 1
(classic), Vite 4, ESLint 8 (`.eslintrc.cjs`), and Vue 3.3 / TypeScript 5.2. The host that
builds and CI already cannot run `yarn` at all, and the current lockfile is a Yarn v1
file. Two follow-up workstreams (responsive UI, PWA) depend on a modern, buildable
foundation, so the toolchain upgrade is the first of three planned changes.

## What Changes

- Bump every runtime and dev dependency in `web/` to the latest stable versions (Vue
  3.5.44, Vue Router 5.4.0, Vite 8.3.5 / Rolldown, TypeScript 6.0.3, ESLint 10 flat
  config, vue-tsc 3.3.13, element-plus 2.14.7, etc.).
- Switch the package manager from Yarn 1 to **pnpm 12.16.1**; replace `web/yarn.lock` with
  a generated `web/pnpm-lock.yaml`.
- Declare `engines.node` and `packageManager` in `web/package.json`.
- Migrate ESLint from `.eslintrc.cjs` + `@rushstack/eslint-patch` to `eslint.config.ts`
  (flat config; ESLint 10 dropped eslintrc).
- Point `tsconfig.node.json` at `@tsconfig/node26`.
- Remove dead/misleading dependencies and files: `@vue/cli-plugin-pwa`,
  `register-service-worker`, `@rushstack/eslint-patch`, `web/src/registerServiceWorker.js`.
- Declare `@element-plus/icons-vue` explicitly (currently resolved only by Yarn 1 hoisting).
- Fix component code that breaks under pnpm's strict `node_modules` and the new Vue/TS:
  `SecretUrl.vue` (`v-model` on `<el-text>`, `type="''"` prop) and `SecretView.vue`
  (top-level `await` + nested `<Suspense>` + indistinguishable 404/network error).
- Update `web/Dockerfile` (builder `node:26`, install pnpm) and the frontend job in
  `.github/workflows/build.yml` (Node 26 + pnpm).

**No behavioural or visual change** to the running app is in scope — that is change 2
(`web-responsive-ui`).

## Capabilities

This change is a pure tooling / dependency / build-config migration. It changes no
spec-level user-visible behaviour, adds no capability, and modifies no existing
capability. Specs are therefore skipped (`skip_specs: true`).

## Impact

- `web/**` (package.json, lockfile, build & lint & tsconfig, `src/**` code fixes,
  `src/registerServiceWorker.js` deletion).
- `web/Dockerfile` and the `frontend` job of `.github/workflows/build.yml` (devops-owned).
- `AGENTS.md` and `AGENTS-TEAM.md` (documented toolchain and DoD) — teamlead-owned, done
  in the same planning PR.
- No change to `gateway/**`, `secret_assistant/**`, `proto/**`, or the HTTP/gRPC contract.
