# Design: PWA (manifest, service worker, offline shell)

## Context

`web/` becomes a PWA on the toolchain from `web-node26-pnpm-deps` and the layout from
`web-responsive-ui`. The hard constraint is the one-time-secret guarantee: **the service
worker must NEVER cache `/api/*`**.

## Non-negotiable constraints (product decisions)

- Service worker: **app-shell precache only; NetworkOnly for `/api/***`. Serving a
  one-time secret from cache would break the "read once" guarantee.
- `navigator.clipboard` and service workers require a **secure context**. Production
  today is plain HTTP (`docker-compose.prod.yml` maps `8080:80`, no TLS anywhere). PWA on
  real devices therefore requires **HTTPS** — this is a product/deployment decision to be
  made before PWA features are released; flag it, do not silently assume it.

## Approach

- Manifest: `web/public/manifest.webmanifest` referencing the existing
  `web/public/img/icons/*` set (maskable 192/512, apple-touch, favicon sizes); add
  `theme-color`, `display: standalone`, `start_url`.
- Service worker: a small precache of the built app shell (JS/CSS/index.html), with a
  runtime handler routing `/api/*` → NetworkOnly (no cache fallback). Use a maintained
  Vite PWA plugin (e.g. `vite-plugin-pwa`) or a hand-rolled worker — decide at
  implementation; the `/api` NetworkOnly rule is mandatory either way.
- Registration: replace the deleted `web/src/registerServiceWorker.js` with a real,
  Vite-aware registration guarded by a secure-context check.
- `index.html`: reference manifest, icons and `theme-color`.

## Element Plus / bundle size

The product owner leans toward dropping Element Plus (app uses only `el-button`,
`el-text`, `el-image`, `el-container`, `ElMessage`; ~353 KB CSS). This change is the
natural point to evaluate it for a smaller mobile bundle. Options, in order of
preference:

1. Replace the five usages with small native/Vue components and drop `element-plus`.
2. Switch to per-component imports (`unplugin-vue-components`) to tree-shake the kit.

Final call at this change's implementation decomposition; the responsive change
(`web-responsive-ui`) must not have already committed to a redesign that depends on the
kit.

**Decision (recorded at `web-pwa` review):** keep `element-plus` as-is (2.14.7). Dropping
it or switching to per-component imports is deferred to a separate change
(`web-bundle-optimization`); the ~353 KB CSS / 1.08 MB JS chunk is a known, accepted cost
for this change and is out of scope here.

## HTTPS deployment

Recorded as an open product/deployment decision: PWA features need HTTPS on real
devices. Until TLS is available, the app degrades gracefully (no SW, no clipboard API) on
plain HTTP and remains fully functional as a normal SPA.

## Out of scope

Layout/responsive work (`web-responsive-ui`), HTTP/gRPC contract changes.

## Risks

- Accidentally caching `/api/*` (guard with an explicit NetworkOnly rule + a qa test that
  asserts a secret is never replayed from cache).
- Releasing PWA features on plain HTTP and confusing users — mitigated by the graceful
  degradation above and the flagged HTTPS decision.
