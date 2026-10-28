# Tasks: PWA (manifest, service worker, offline shell)

Role ownership: `frontend` = `web/**` except Dockerfile/nginx; `devops` = Dockerfile /
CI / deployment config; `qa` = verification; `teamlead` = decisions + docs.

## 1. Manifest and icons (frontend)

- [x] 1.1 Add `web/public/manifest.webmanifest` (name, theme-color, display standalone,
  start_url) referencing the existing `web/public/img/icons/*`.
- [x] 1.2 Wire manifest, icons and `theme-color` into `web/index.html`.

## 2. Service worker (frontend)

- [x] 2.1 Add a service worker that precaches the app shell and routes `/api/*` as
  NetworkOnly (never cached).
- [x] 2.2 Replace the deleted `registerServiceWorker.js` with a real Vite-aware
  registration guarded by a secure-context check.
- [x] 2.3 `pnpm build && pnpm lint && pnpm type-check` green.

## 3. HTTPS and deployment (devops)

- [x] 3.1 Resolve/record the HTTPS secure-context decision (flag as a product/deployment
  decision); ensure the app degrades gracefully on plain HTTP.
  RESOLVED (product decision, recorded in design.md): TLS/HTTPS are planned — the
  service WILL be served over HTTPS. There is no deployment server yet, which is why the
  `deploy` step was removed from CI (commit 62e60b8). Deployment and the
  installability-over-HTTPS verification are deferred until a server exists. Graceful
  degradation on plain HTTP is implemented (`registerServiceWorker.ts` guards on
  `window.isSecureContext`).

## 4. Verify (qa)

- [x] 4.1 Assert a one-time secret is never served from cache (open once, then confirm a
  second/offline attempt does not replay the secret).
  Verified: `/api/*` is NetworkOnly in `vite.config.ts` and in compiled `dist/sw.js`
  (`registerRoute(/\/api\//, new NetworkOnly, "GET")`); precache = 22 app-shell entries,
  zero `/api` entries; runtime cache inspection after an `/api` fetch shows no
  write-through. Literal end-to-end replay run deferred — no Docker/backend on the host.
- [x] 4.2 Verify installability + offline app shell over HTTPS, and graceful degradation
  over plain HTTP.
  Graceful degradation over plain HTTP VERIFIED (headless browser: `isSecureContext`
  false, no SW registration, app renders as a normal SPA). Installability + offline
  shell over HTTPS DEFERRED — no secure context/server available (see 3.1).

## 5. Documentation (teamlead)

- [x] 5.1 Record the Element Plus drop/per-component decision and the HTTPS decision in
  `design.md` at review time.
