# Design: Issue templates

## Context

The repository has no issue templates; issues arrive as free-form text. See proposal.md —
Why.

## Goals / Non-Goals

**Goals:**
- Structured, consistently-labelled bug and feature requests.

**Non-Goals:**
- No code, no contract change, no automation beyond the templates themselves.

## Decisions

- **GitHub issue forms (`.yml`) rather than Markdown templates.** Forms enforce required
  fields (version, environment, reproduction) and validate input, which Markdown
  templates cannot; they also render a cleaner "New issue" chooser.
- **Two forms plus a `config.yml`.** `bug_report.yml` and `feature_request.yml` cover the
  two recurring kinds of issue; `config.yml` sets `blank_issues_enabled: false` so
  reports are routed through a template without blocking ad-hoc discussion.
- **Ownership is `devops`** because the files live under `.github/**`.

## Risks / Trade-offs

- [Overly rigid templates deter quick reports] → mitigated by keeping fields short and
  marking only the version/description as required.
