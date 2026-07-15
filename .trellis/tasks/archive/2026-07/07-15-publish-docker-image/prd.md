# Publish Docker image

## Goal

Automatically build a Docker image for `sobber98/CLIProxyAPI` and publish it for project contributors.

## Confirmed Facts

- The repository has a multi-stage `Dockerfile` that builds `./cmd/server` and accepts `VERSION`, `COMMIT`, and `BUILD_DATE` build arguments.
- There is no existing `.github/workflows/` directory or CI publishing workflow.
- The local `origin` remote is `https://github.com/sobber98/CLIProxyAPI.git`; this public fork currently has no published GitHub packages or releases.
- The image will publish to GitHub Container Registry as `ghcr.io/sobber98/cliproxyapi`, authenticated with the GitHub Actions `GITHUB_TOKEN`.
- The repository has no Git tags. The local Docker build scripts derive version, commit, and UTC build date from Git.
- Publishing must be initiated manually through GitHub Actions `workflow_dispatch`; pushes and pull requests must not publish images.
- Published images must support `linux/amd64` and `linux/arm64`.

## Requirements

- Add a GitHub Actions workflow that builds the existing Docker image and publishes it to the intended Contributors container repository.
- Pass release metadata to the existing Dockerfile build arguments.
- On every manual run, publish only the `latest` tag.
- Do not alter application runtime behavior unless required to support image publication.

## Acceptance Criteria

- [ ] A GitHub Actions workflow builds the repository's existing Docker image.
- [ ] The workflow publishes image tags to `ghcr.io/sobber98/cliproxyapi`.
- [ ] The workflow can only publish when manually dispatched from GitHub Actions.
- [ ] Each published tag resolves to an image manifest supporting `linux/amd64` and `linux/arm64`.
- [ ] A successful manual run updates only `ghcr.io/sobber98/cliproxyapi:latest`.
- [ ] Published builds include source version, commit, and build-date metadata.
- [ ] The workflow configuration declares its manual-only trigger, tag policy, platforms, and minimum GitHub token permissions.
