# Upstream Merge Design

## Integration Strategy

Merge `upstream/main` into local `main` with a merge commit. This preserves the
already-pushed credential-grouping commit and local Trellis bookkeeping history.
The completed merge remains local and is not pushed to `origin/main`.

## Conflict Resolution

Preflight identifies conflicts in `internal/config/config.go`,
`internal/watcher/synthesizer/config.go`, and
`sdk/api/handlers/openai/codex_client_models.go`.

For each conflict, retain upstream's current behavior and incorporate the
credential-grouping extension:

- Configuration types and normalizers retain `APIKeyGroups` and upstream
  configuration fields/normalization.
- Config credential synthesis retains the upstream model and provider behavior
  while setting runtime `Auth.Group` from configured credentials.
- Codex client-version model catalogs retain upstream catalog generation and
  filter results through the group-aware discovery helper.

## Safety And Rollback

Untracked `.trellis/` files remain outside Git operations. If conflict
resolution or verification exposes a regression, abort the merge to restore the
pre-merge local mainline, leaving the feature commit intact.
