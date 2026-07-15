# Implementation Plan: Manual GHCR Publication

1. Create `.github/workflows/docker-publish.yml` with a `workflow_dispatch` trigger and least-privilege `contents: read` and `packages: write` permissions.
2. Configure Buildx and login to `ghcr.io` using `GITHUB_TOKEN`.
3. Build and push the existing Dockerfile for `linux/amd64,linux/arm64`, publishing only `ghcr.io/sobber98/cliproxyapi:latest`.
4. Derive `VERSION=latest`, the selected revision's short SHA, and UTC build timestamp, then pass all three Docker build arguments.
5. Validate the workflow YAML and inspect the generated diff for trigger, permissions, tag, platform, metadata, and push behavior.

## Validation

- Parse the workflow YAML with an available YAML parser or action linter.
- Confirm the workflow contains only `workflow_dispatch` as its event trigger.
- Confirm the build command targets both configured Linux architectures and exactly the approved image tag.

## Rollback

Delete `.github/workflows/docker-publish.yml` to remove the capability. For a published-image rollback, manually dispatch the workflow against a known-good revision.
