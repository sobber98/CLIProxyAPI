# Upstream Merge Implementation Plan

1. Fetch and retain `upstream/main` as a tracked remote reference.
2. Start a non-fast-forward merge into local `main`; resolve each conflict by
   combining upstream behavior with credential grouping.
3. Inspect the merge result to ensure group selection, fallback, management,
   and discovery changes remain present.
4. Run formatting, focused grouping tests, and the required server build. Run
   broader tests when practical and distinguish existing unrelated failures.
5. Commit the merge locally. Do not push to `origin/main`.

## Validation

- `gofmt -w` on manually resolved Go files.
- Focused tests for configuration, access, auth manager, handler/model catalog,
  management, watcher synthesis, and TUI packages.
- `go build -o test-output ./cmd/server && rm test-output`.
- `git diff --check` and review of merge parents/history.

## Rollback

- Before the merge commit, use `git merge --abort` if resolution is unsafe.
- After committing, use a separate revert commit if rollback is required; do not
  rewrite the previously published credential-grouping commit.
