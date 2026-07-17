# Merge upstream updates

## Goal

Bring the current fork's `main` branch up to date with the upstream CLIProxyAPI
repository, preserving fork-specific commits and existing uncommitted workspace
files.

## Confirmed Facts

- The current branch is `main`.
- `origin` is `https://github.com/sobber98/CLIProxyAPI.git`.
- `upstream` is `https://github.com/router-for-me/CLIProxyAPI.git`.
- The worktree contains untracked `.trellis/` files. They must not be deleted,
  staged, or included in the upstream synchronization commit.

## Requirements

- Fetch the current upstream branch and determine divergence from local `main`.
- Integrate upstream changes with Git merge semantics, resolving conflicts
  without discarding fork-specific changes.
- Build and test the resulting tree.
- Push the completed `main` branch to `origin` only after local verification.

## Acceptance Criteria

- [ ] Local `main` contains the fetched upstream branch history.
- [ ] Existing local commits and untracked `.trellis/` files are preserved.
- [ ] `go build -o test-output ./cmd/server && rm test-output` succeeds.
- [ ] Relevant test suite succeeds, or any pre-existing failure is documented.
- [ ] The verified local `main` is pushed to `origin/main`.

## Out Of Scope

- Changing application behavior unrelated to resolving merge conflicts.
- Adding `.trellis/` files to Git.
