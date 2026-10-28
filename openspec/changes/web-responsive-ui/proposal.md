# Proposal: Mobile/tablet responsive UI

## Why

The app is usable on desktop but is not designed for mobile or tablet: fixed-width
containers, desktop-oriented spacing, and touch targets that are too small. The product
owner wants a responsive UI as the second of three planned web changes (after the
toolchain upgrade in `web-node26-pnpm-deps`, before the PWA in `web-pwa`).

## What Changes

- Make the whole SPA responsive across mobile (< 768px), tablet (768–1100px) and desktop
  (> 1100px) viewports: layout, typography, spacing, and input/button sizing.
- Ensure the secret editor, the generated secret URL, and the secret-viewer screen are
  usable and readable on narrow screens with no horizontal scrolling.
- Enlarge touch targets and improve keyboard/pointer ergonomics on touch devices.
- Ensure the one-time secret flow (create → share → open once → consumed 404) works
  identically on mobile.

This change owns layout and presentation only. It does **not** add a service worker,
manifest, offline support, or installability (that is `web-pwa`).

## Capabilities

### New Capabilities

- `web-responsive-layout`: the SPA renders a usable, readable, non-horizontally-scrolling
  layout across mobile, tablet and desktop viewports.

## Impact

- `web/src/**` (views, components, global CSS): responsive layout and styling.
- `web/index.html` (viewport meta already present; theme-color/meta may be added if
  needed for mobile rendering).
- No change to `gateway/**`, `secret_assistant/**`, `proto/**`, or the HTTP contract.

## Sequencing

This change follows `web-node26-pnpm-deps` (the toolchain must be green first) and
precedes `web-pwa` (a PWA should be responsive before it is offline-capable). The three
changes are delivered in that order; `web-responsive-ui` has no hard dependency on
`web-pwa`.
