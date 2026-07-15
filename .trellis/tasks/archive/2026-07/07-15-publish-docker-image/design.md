# Design: Manual GHCR Publication

## Architecture

Add one GitHub Actions workflow under `.github/workflows/`. It runs only through `workflow_dispatch` and uses Docker Buildx to build the repository's existing `Dockerfile` for `linux/amd64` and `linux/arm64`.

The workflow authenticates to GitHub Container Registry with the repository-scoped `GITHUB_TOKEN` and `packages: write` permission. It publishes the combined multi-platform manifest as `ghcr.io/sobber98/cliproxyapi:latest`.

## Build Metadata

The workflow derives the selected source revision's short commit SHA and the current UTC timestamp at runtime. It passes them as `COMMIT` and `BUILD_DATE`. `VERSION` is set to `latest`, matching the sole published image tag. The existing Dockerfile embeds these values in the Go binary.

## Boundaries and Compatibility

- The Dockerfile and application code remain unchanged.
- No Docker Hub or repository secrets are required.
- Manual publication lets an authorized repository user select the Git ref in the GitHub Actions UI; no push or pull-request event publishes an image.
- Re-running the workflow intentionally replaces `latest` with the selected revision.

## Failure and Rollback

Build or registry failures leave the prior published `latest` image intact. To roll back a bad image, manually rerun the workflow against the previously known-good Git ref.
