# Tasks: Mobile/tablet responsive UI

Role ownership: `frontend` = `web/**` except Dockerfile/nginx; `qa` = verification.

## 1. Responsive layout (frontend)

- [x] 1.1 Implement mobile/tablet/desktop breakpoints and responsive styles across
  `web/src/**` (views, components, `index.css`).
- [x] 1.2 Ensure no horizontal scroll at 360px and that touch targets meet minimum size.
- [x] 1.3 Keep desktop layout structurally unchanged; `pnpm build && pnpm lint &&
  pnpm type-check` green.

## 2. Verify (qa)

- [ ] 2.1 Verify the create → share → open-once → consumed-404 flow at each breakpoint.
- [ ] 2.2 Verify no horizontal scrolling and readable content on mobile/tablet/desktop;
  capture evidence (screenshots) at each breakpoint.

## 3. Documentation (teamlead)

- [x] 3.1 Record the final breakpoint/decision choices in `design.md` at review time.
