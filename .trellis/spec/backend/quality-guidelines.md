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

---

## Scenario: Management Credential Group Editing

### 1. Scope / Trigger

The Management Center groups configured credentials, OAuth auth files, and
downstream API keys. Group mutations cross the panel, Management API, config or
auth-file persistence, and subsequent Management list responses.

### 2. Signatures

- `GET /v0/management/auth-files` returns `files[].group` as the trimmed
  persisted group in both runtime-auth and disk-fallback responses.
- `PATCH /v0/management/auth-files/fields` accepts `{"name":"...","group":"team-a"}`;
  `group: ""` clears the assignment.
- Configured API-key patch resources, including `PATCH /v0/management/xai-api-key`,
  accept `{"match":"api-key","value":{"group":"team-a"}}`.
- `GET`/`PUT /v0/management/api-key-groups` retain their wrapper contract
  `{"api-key-groups":{"client-key":"team-a"}}`; remove a mapping to clear it.

### 3. Contracts

- Group values are trimmed and case-sensitive. Missing or empty values mean
  ungrouped.
- A Management list response must round-trip every editable group field. A
  successful patch whose list response omits `group` is a broken UI contract.
- Client bulk updates must select configured credentials by stable identity:
  `match` for API-key resources and `name` for OpenAI-compatible providers.
  Do not send a stale list index alongside an identity selector.
- OpenAI-compatible `api-key-entries` need both their original index and API-key
  value. Re-fetch the provider before its full-list PATCH, verify the selected
  index still has the selected API key, and fail rather than update a moved item.

### 4. Validation & Error Matrix

| Condition | Required behavior |
| --- | --- |
| OAuth group is patched then listed | Runtime and disk fallback each return the trimmed group. |
| xAI group patch assigns or clears a group | Persist the trimmed value, including `""` for ungrouped. |
| Provider list is concurrently reordered | Identity-only patch selects the matching record; no stale index is used. |
| OpenAI API-key entry moved or removed | Reject that selected entry; never mutate another entry with the same key. |
| `api-key-groups` endpoint unavailable | Keep upstream grouping usable and disable/report only downstream grouping. |

### 5. Good/Base/Bad Cases

- Good: an OAuth group saved through Management appears after panel refresh and
  restricts selection once a downstream key maps to it.
- Base: clearing a credential group persists `""`; clearing a downstream key
  removes its `api-key-groups` entry.
- Bad: the panel displays all OAuth credentials as ungrouped after refresh, or
  a stale index changes a different credential after an operator reorders the
  configuration.

### 6. Tests Required

- Management tests cover OAuth group patch followed by runtime listing and
  disk-fallback listing.
- Management tests cover xAI group assignment, trimming, and clearing.
- Panel tests cover group API payloads, downstream mapping removal, partial
  failures, OpenAI-compatible duplicate API keys, and reordered entry refusal.
- Panel tests assert configured credential and provider mutations omit stale
  indexes when a stable identity is available.

### 7. Wrong vs Correct

#### Wrong

```ts
await apiClient.patch('/gemini-api-key', {
  index: loadedIndex,
  match: apiKey,
  value: { group },
});
```

#### Correct

```ts
await apiClient.patch('/gemini-api-key', {
  match: apiKey,
  value: { group },
});
```
