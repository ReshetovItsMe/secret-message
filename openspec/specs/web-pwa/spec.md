# web-pwa Specification

## Purpose
Makes the secret-message web app installable and able to serve a cached app shell
offline, while guaranteeing that one-time secrets are never served from cache.

## Requirements

### Requirement: Installable web app
The application SHALL be installable via a web app manifest that declares a name, icons
(including maskable), a theme colour and `display: standalone`.

#### Scenario: Install prompt is available
- **WHEN** the app is served over HTTPS and opened in a supporting browser
- **THEN** the browser recognises the manifest and the app can be added to the home screen

### Requirement: App shell precache with NetworkOnly API
The service worker SHALL precache the app shell and SHALL use a NetworkOnly strategy for
every request under `/api/`, so no secret payload is ever served from the cache.

#### Scenario: Secret is never served from cache
- **WHEN** the user opens a one-time secret URL while offline
- **THEN** the `/api/*` request is not answered from the cache, and the one-time secret
  is never replayed from a cached response

#### Scenario: App shell is available offline
- **WHEN** the app shell assets have been precached and the device goes offline
- **THEN** the app shell still loads so the user reaches the UI (the API call itself fails
  or returns the consumed/missing state, never a cached secret)

### Requirement: Secure context for service worker and clipboard
The service worker and `navigator.clipboard` SHALL only operate in a secure context, and
the deployment SHALL document that HTTPS is required for PWA features on real devices.

#### Scenario: Non-secure context does not break the app
- **WHEN** the app is served over plain HTTP (no secure context)
- **THEN** the service worker and clipboard are not registered/used, and the rest of the
  app continues to function without PWA features
