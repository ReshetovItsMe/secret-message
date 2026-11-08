# Tasks: Issue templates

Role ownership: `devops` = `.github/**`.

## 1. Add issue templates (devops)

- [x] 1.1 Add `.github/ISSUE_TEMPLATE/bug_report.yml` (version, environment, reproduction
  steps, expected vs actual behaviour, logs).
- [x] 1.2 Add `.github/ISSUE_TEMPLATE/feature_request.yml` (problem statement, proposed
  solution, alternatives, scope).
- [x] 1.3 Add `.github/ISSUE_TEMPLATE/config.yml` with `blank_issues_enabled: false` and a
  link to the contributing notes.

## 2. Verification (devops)

- [x] 2.1 Confirm all three YAML files parse and pass local structural checks.
- [ ] 2.2 After integration to the default branch, confirm both templates render in the
  GitHub "New issue" chooser, both forms open, the contributing link works, and
  the blank-issue path is disabled. Track evidence in kanban t_5597a968.

## Integration exception

Owner-authorized sequencing resolution (t_2c03fe9d): integrate after local validation
and green PR CI, before live default-branch verification. Archive retains task 2.2
as pending; archiving is not a claim of final acceptance. This external verification
does not gate the independent golangci-lint change. See verification.md.
