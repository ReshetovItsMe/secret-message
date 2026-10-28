# Tasks: PWA (manifest, service worker, offline shell)

Role ownership: `frontend` = `web/**` except Dockerfile/nginx; `devops` = Dockerfile /
CI / deployment config; `qa` = verification; `teamlead` = decisions + docs.

## 1. Manifest and icons (frontend)

- [ ] 1.1 Add `web/public/manifest.webmanifest` (name, theme-color, display standalone,
  start_url) referencing the existing `web/public/img/icons/*`.
- [ ] 1.2 Wire manifest, icons and `theme-color` into `web/index.html`.

## 2. Service worker (frontend)

- [ ] 2.1 Add a service worker that precaches the app shell and routes `/api/*` as
  NetworkOnly (never cached).
- [ ] 2.2 Replace the deleted `registerServiceWorker.js` with a real Vite-aware
  registration guarded by a secure-context check.
- [ ] 2.3 `pnpm build && pnpm lint && pnpm type-check` green.

## 3. HTTPS and deployment (devops)

- [ ] 3.1 Resolve/record the HTTPS secure-context decision (flag as a product/deployment
  decision); ensure the app degrades gracefully on plain HTTP.

## 4. Verify (qa)

- [ ] 4.1 Assert a one-time secret is never served from cache (open once, then confirm a
  second/offline attempt does not replay the secret).
- [ ] 4.2 Verify installability + offline app shell over HTTPS, and graceful degradation
  over plain HTTP.

## 5. Documentation (teamlead)

- [ ] 5.1 Record the Element Plus drop/per-component decision and the HTTPS decision in
  `design.md` at review time.
