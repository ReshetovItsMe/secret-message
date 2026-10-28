# Design: Mobile/tablet responsive UI

## Context

`web/` becomes a responsive SPA on the toolchain delivered by `web-node26-pnpm-deps`
(Node 26 + pnpm + Vite 8 + Vue 3.5). This change is presentation-only; it does not add a
service worker or offline support.

## Approach

- Define three breakpoints and drive the layout from them: mobile `< 768px`, tablet
  `768–1100px`, desktop `> 1100px`. The current `MainView.vue` already has a 1100px split
  for the editor width — keep that boundary as the desktop/tablet seam and add the mobile
  seam below it.
- Work in `web/src/**` (views, components, `index.css`). Prefer CSS media queries and the
  existing scoped styles over a new framework; no new runtime dependency for layout.
- Preserve the existing header/main/footer structure on desktop.

## Element Plus decision (deferred to `web-pwa`)

The product owner leans toward dropping Element Plus (the app uses only `el-button`,
`el-text`, `el-image`, `el-container`, `ElMessage`; the kit ships ~353 KB of CSS). This
change **keeps** the current Element Plus usage (already upgraded to 2.14.7 in change 1)
so that responsive layout is isolated from any component-kit migration. The drop /
per-component-import decision is evaluated in `web-pwa`, where bundle size matters most
for a mobile offline app.

## Finalized decisions (recorded at review)

- Breakpoints: mobile `< 768px`, tablet `768–1099px`, desktop `≥ 1100px`. The desktop
  seam (`min-width: 1100px`) is the existing `MainView.vue` editor-width split, unchanged;
  the new mobile seam is `max-width: 767px`; tablet is the band between them.
- Footer icons are **not** reordered on narrow screens: the footer is centered
  (`justify-content: center`) below `768px` instead.
- `web/index.html` gains **no** `theme-color` or mobile meta in this change — deferred to
  `web-pwa`, which wires the manifest, icons and `theme-color`.

## Out of scope

Service worker, manifest, offline caching, installability (`web-pwa`). HTTP/gRPC contract,
`gateway/**`, `secret_assistant/**`.

## Risks

- A visual regression on desktop while fixing mobile — mitigate with qa screenshots at
  each breakpoint.
- Over-scoping into a full redesign; the acceptance target is "usable, readable, no
  horizontal scroll", not a pixel-perfect new design.
