# Proposal: PWA (manifest, service worker, offline shell)

## Why

The app is currently a plain SPA with no manifest, no service worker and no offline
capability, despite `AGENTS.md` having claimed PWA support. The product owner wants a
real PWA as the third of three planned web changes: after the toolchain upgrade
(`web-node26-pnpm-deps`) and the responsive UI (`web-responsive-ui`). The one-time
secret flow has a hard caching constraint (a secret must never be served from cache), so
the PWA must be designed around app-shell-only precaching.

## What Changes

- Add a web app manifest (name, icons, theme-color, display `standalone`).
- Add a service worker that precaches the app shell and uses a NetworkOnly strategy for
  every `/api/*` request — so the "read once" guarantee is never broken by the cache.
- Wire the icons already present in `web/public/img/icons/*` (maskable + apple-touch)
  into `index.html` and the manifest.
- Replace the dead `web/src/registerServiceWorker.js` with a real, Vite-aware
  registration path.
- Flag and resolve the HTTPS secure-context requirement (service workers and
  `navigator.clipboard` require a secure context) — a product/deployment decision.
- Evaluate dropping Element Plus (or switching to per-component imports) to shrink the
  bundle for mobile installs (see design).

This change adds installability and offline shell; it does **not** change layout
(that is `web-responsive-ui`) or the HTTP/gRPC contract.

## Capabilities

### New Capabilities

- `web-pwa`: the app is installable and serves a cached app shell offline while every
  `/api/*` request stays NetworkOnly so a one-time secret is never served from cache.

## Impact

- `web/**`: manifest, service worker registration, `index.html`, icons, possibly a
  dependency reduction (Element Plus).
- Deployment config only if HTTPS/TLS is enabled (flagged as a product decision).
- No change to `gateway/**`, `secret_assistant/**`, `proto/**`, or the HTTP contract.

## Sequencing

This change follows `web-node26-pnpm-deps` and `web-responsive-ui`. It is the last of the
three; it has no hard dependency on the other two beyond a green toolchain and a
responsive layout being in place.
