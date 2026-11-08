# Proposal: Issue templates

## Why

Issue #18. The repository has no issue templates, so bug reports and feature requests
arrive as free-form text that omits the version, reproduction steps and expected
behaviour the team needs to triage efficiently. GitHub issue forms give structured,
consistently-labelled issues without blocking ad-hoc reports.

## What Changes

- Add `.github/ISSUE_TEMPLATE/bug_report.yml` (version, environment, steps to reproduce,
  expected vs actual behaviour, logs).
- Add `.github/ISSUE_TEMPLATE/feature_request.yml` (problem statement, proposed solution,
  alternatives, scope).
- Add `.github/ISSUE_TEMPLATE/config.yml` to set `blank_issues_enabled: false` and link
  to the contributing notes.

## Capabilities

### New Capabilities

_None._ This is repository process metadata; there is no product-behaviour change, so the
change declares `skip_specs: true`.

## Impact

- `.github/ISSUE_TEMPLATE/*.yml` (new files). No code, no contract change.

## Sequencing

Independent of every other change.
