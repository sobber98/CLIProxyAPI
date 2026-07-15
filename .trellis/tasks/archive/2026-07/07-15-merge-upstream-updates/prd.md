# Merge upstream updates

## Goal

Merge the current upstream `router-for-me/CLIProxyAPI:main` updates into the
fork's `main` branch without regressing the credential grouping feature.

## Confirmed Facts

- The fork remote is `origin` (`sobber98/CLIProxyAPI`); the canonical upstream
  is `https://github.com/router-for-me/CLIProxyAPI.git`.
- Local `main` contains the credential-grouping feature commit `d530e7a3` plus
  local Trellis archive and journal commits.
- The fetched upstream `main` contains many commits not present locally,
  including changes in auth selection, executor behavior, translator workflows,
  request logging, UI, usage, and model catalog areas.
- The worktree contains only untracked Trellis runtime, task, spec-template, and
  workspace files. They must not be overwritten, staged, or included in the
  upstream merge.

## Requirements

- Integrate the fetched upstream mainline history into local `main`.
- Use a merge commit; do not rebase or rewrite the already-pushed local history.
- Preserve the credential-grouping behavior: group-constrained credential
  selection and failover, model discovery filtering, configuration validation,
  Management API, and TUI editing.
- Resolve conflicts in favor of the current feature contract where upstream does
  not provide an equivalent mechanism.
- Do not modify or include untracked `.trellis/` runtime/task files.
- Keep the completed merge local; do not push it to `origin/main`.
- Validate the resulting merge with focused credential-grouping tests and the
  required server build; run broader tests when practical and report unrelated
  pre-existing failures separately.

## Acceptance Criteria

- [ ] Local `main` contains the upstream mainline history and the
  credential-grouping commit remains reachable.
- [ ] Group constraints still apply across normal selection, scheduler/failover,
  and model-discovery paths after conflict resolution.
- [ ] Focused credential-grouping tests and server build pass.
- [ ] Untracked `.trellis/` runtime/task files remain unmodified and uncommitted.

## Out Of Scope

- Rewriting published history through rebase.
- Pushing the resulting merge to `origin/main`.
