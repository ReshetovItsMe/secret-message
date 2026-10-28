## Purpose

Defines the responsive behaviour of the secret-message single-page application so the
create-share-read-once flow is usable on mobile, tablet and desktop viewports.

## ADDED Requirements

### Requirement: Responsive layout across viewports
The application SHALL render a usable layout at mobile (< 768px), tablet (768–1100px) and
desktop (> 1100px) widths, with no horizontal scrolling at any width the app supports.

#### Scenario: Mobile layout has no horizontal scroll
- **WHEN** the app is viewed at a 360px-wide viewport
- **THEN** no element overflows horizontally and the page does not scroll on the x-axis

#### Scenario: Desktop layout is unchanged in structure
- **WHEN** the app is viewed at a width greater than 1100px
- **THEN** the existing header/main/footer arrangement is preserved and remains usable

### Requirement: Usable secret editor on touch devices
The secret input and its action button SHALL have touch targets large enough for reliable
operation on a touch screen and remain reachable on narrow screens.

#### Scenario: Share button is tappable on mobile
- **WHEN** the app is viewed on a mobile viewport
- **THEN** the Share button is at least the platform-recommended touch-target size and
  reacts to a tap to submit the secret

### Requirement: Secret flow works on mobile
The one-time secret flow SHALL behave identically on mobile as on desktop: creating a
secret yields a shareable URL, opening it reveals the message exactly once, and a second
open returns the consumed-secret (404) state.

#### Scenario: Create and open once on mobile
- **WHEN** a user creates a secret and opens its URL once on a mobile device
- **THEN** the message is displayed, and a subsequent open shows the consumed-secret state

#### Scenario: Consumed secret is distinguishable from network failure
- **WHEN** opening an already-consumed or unknown secret fails
- **THEN** the user is shown a distinct consumed/missing message, not a generic error
