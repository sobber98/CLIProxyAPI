# Update Docker image reference

## Goal

Make Docker Compose pull the published CLIProxyAPI image from the requested GitHub Container Registry location by default.

## Confirmed Facts

- `docker-compose.yml:3` and `docker-compose.cluster.yml:3` both set the default `CLI_PROXY_IMAGE` value to `eceasy/cli-proxy-api:latest`.
- Both Compose services use `pull_policy: always`.

## Requirements

- Change the default image in both Compose files to `ghcr.io/sobber98/cliproxyapi:latest`.
- Preserve the `CLI_PROXY_IMAGE` environment override and all unrelated Compose configuration.

## Acceptance Criteria

- [ ] `docker-compose.yml` defaults `CLI_PROXY_IMAGE` to `ghcr.io/sobber98/cliproxyapi:latest`.
- [ ] `docker-compose.cluster.yml` defaults `CLI_PROXY_IMAGE` to `ghcr.io/sobber98/cliproxyapi:latest`.
- [ ] Both Compose files remain valid after interpolation.

## Out of Scope

- Changing Docker build settings, container runtime configuration, or image publishing workflows.
