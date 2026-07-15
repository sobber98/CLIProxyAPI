# Quality Guidelines

> Code quality standards for backend development.

---

## Overview

<!--
Document your project's quality standards here.

Questions to answer:
- What patterns are forbidden?
- What linting rules do you enforce?
- What are your testing requirements?
- What code review standards apply?
-->

(To be filled by the team)

---

## Forbidden Patterns

<!-- Patterns that should never be used and why -->

(To be filled by the team)

---

## Required Patterns

<!-- Patterns that must always be used -->

(To be filled by the team)

---

## Testing Requirements

<!-- What level of testing is expected -->

(To be filled by the team)

---

## Code Review Checklist

<!-- What reviewers should check -->

(To be filled by the team)

---

## Scenario: Manual Multi-Platform GHCR Publication

### 1. Scope / Trigger

The repository publishes its server image through `.github/workflows/docker-publish.yml`. This CI integration is a deployment contract because it controls the registry credential, public image tag, image architectures, and binary build metadata.

### 2. Signatures

- Trigger: GitHub Actions `workflow_dispatch` only.
- Image: `ghcr.io/sobber98/cliproxyapi:latest`.
- Platforms: `linux/amd64,linux/arm64`.
- Docker build arguments: `VERSION=latest`, `COMMIT=<short Git SHA>`, and `BUILD_DATE=<UTC RFC 3339 timestamp>`.

### 3. Contracts

- Workflow permissions are `contents: read` and `packages: write`; authenticate to GHCR with `${{ secrets.GITHUB_TOKEN }}`.
- QEMU must be initialized before Buildx so an amd64 GitHub-hosted runner can produce the arm64 image variant.
- The build action pushes a combined multi-platform manifest and publishes no tag other than `latest`.

### 4. Validation & Error Matrix

| Condition | Required behavior |
| --- | --- |
| Workflow is not manually dispatched | Do not build or publish an image. |
| GHCR login or build fails | Fail the workflow; the previously published `latest` image remains unchanged. |
| Multi-platform build misses an architecture | Do not publish a replacement manifest. |
| Bad image must be rolled back | Manually dispatch the workflow for a known-good Git ref. |

### 5. Good/Base/Bad Cases

- Good: A maintainer manually runs the workflow and receives a `latest` manifest for both required architectures with the selected ref's metadata.
- Base: Re-running a known-good ref intentionally replaces `latest` with an equivalent supported manifest.
- Bad: A push event automatically overwrites `latest`, or a build publishes only the runner architecture.

### 6. Tests Required

- Parse the workflow YAML and assert `workflow_dispatch` is its only trigger.
- Assert the permissions, GHCR tag, platform list, Buildx/QEMU setup, push flag, and three Docker build arguments.
- Run `go test ./...` and `go build -o test-output ./cmd/server && rm test-output` for the source image build contract.

### 7. Wrong vs Correct

#### Wrong

```yaml
on:
  push:
    branches: [main]
```

#### Correct

```yaml
on:
  workflow_dispatch:
```

---

## Scenario: Credential Group Isolation

### 1. Scope / Trigger

Credential grouping is a cross-layer authorization contract: a downstream proxy
API key selects one upstream credential group. It applies to configured API-key
credentials, file-backed OAuth credentials, request failover, and model
discovery.

### 2. Signatures

- Config: `api-key-groups: map[string]string` maps a value in `api-keys` to one
  group. Configured upstream credentials and OAuth auth JSON accept `group`.
- Management: `GET` and `PUT /v0/management/api-key-groups` use the wrapper
  payload `{"api-key-groups":{"client-key":"team-a"}}`.
- Execution metadata: `credential_group` and `credential_group_strict` carry
  the authenticated group's constraint from access middleware to auth selection.

### 3. Contracts

- Group names are trimmed and case-sensitive. Empty or absent means ungrouped.
- Downstream API keys remain opaque exact values; only group names are trimmed.
- Strict isolation is active only when `api-key-groups` is non-empty. An
  unassigned downstream key then matches only ungrouped upstream credentials.
- Every credential candidate source, including scheduler, retry, model pools,
  and Antigravity credits fallback, must use the same group predicate.
- Model discovery must filter its catalog with the same route-model compatibility
  rules used for selection, including aliases and Codex client-version catalogs.

### 4. Validation & Error Matrix

| Condition | Required behavior |
| --- | --- |
| `api-key-groups` contains an unknown key | Reject config load, hot reload, YAML update, and Management mutation. |
| API key rename has an assignment | Move its assignment atomically. |
| API key deletion has an assignment | Remove its assignment atomically. |
| Strict request has no matching credential | Return normal no-auth-available behavior; never fall back across groups. |
| Empty mapping | Preserve legacy unrestricted selection and model catalog behavior. |

### 5. Good/Base/Bad Cases

- Good: `team-a-key` maps to `team-a`; a `team-a` OAuth or configured credential
  is selected and only its routable models are listed.
- Base: Credentials have groups but `api-key-groups` is empty; legacy routing
  continues unchanged until an operator enables the mapping.
- Bad: A request for `team-a` falls through to an ungrouped or `team-b`
  credential after a failure.

### 6. Tests Required

- Config validation must assert trimming and rejection of an unknown mapped key.
- Access-provider tests must assert group and strict metadata propagation.
- Auth-manager tests must assert group filtering for normal selection and every
  fallback path, including credits fallback.
- Model-listing tests must assert grouped callers cannot discover other groups'
  models and aliases remain discoverable when their credential can route them.
- Management and TUI tests must assert group updates persist, invalid mutations
  leave the current mapping intact, and API-key rename/delete preserves mapping
  integrity.

### 7. Wrong vs Correct

#### Wrong

```go
// A fallback candidate can bypass the group constraint.
for _, candidate := range fallbackCandidates {
    return candidate
}
```

#### Correct

```go
for _, candidate := range fallbackCandidates {
    if !authMatchesCredentialGroup(candidate, opts.Metadata) {
        continue
    }
    return candidate
}
```
