# Credential Grouping Design

## Boundary And Configuration

Add `api-key-groups` as a top-level `map[string]string` beside the compatible
`api-keys: []string` list. Each configured upstream API-key credential type
gets an optional `group` field. File-backed OAuth credentials use optional
`group` in their existing JSON metadata.

Groups are trimmed and case-sensitive. A missing or empty group denotes the
ungrouped bucket. Configuration validation trims group values, rejects mappings
whose downstream key is absent from `api-keys`, and preserves legacy behavior
when `api-key-groups` is empty.

## Runtime Data Flow

1. The config access provider authenticates the request and resolves the
   authenticated downstream API key to its configured group.
2. The access result carries that group as metadata. `AuthMiddleware` puts the
   metadata in the Gin context.
3. Handler execution metadata copies the group constraint to executor options.
4. The auth manager filters candidate runtime `Auth` entries by group before
   availability checks, selector scheduling, retries, and failover. The same
   predicate is used in legacy, built-in scheduler, mixed-provider, and
   model-listing paths.
5. If strict isolation is active and no matching credential exists, normal
   no-auth-available behavior is returned. No fallback may escape the group.

`Auth` gains a group field so configured and file-backed credentials expose one
uniform value to the selection layer.

## Model Discovery

Authenticated model-discovery requests use the same group constraint to filter
models by credentials that can route them. Unauthenticated deployments retain
the legacy full catalog. This prevents discovery of models that a caller cannot
execute.

## Management And TUI

Expose a dedicated Management resource for `api-key-groups`, with validated
read/write operations that persist the config and trigger the existing reload
path. Extend configured credential payloads with `group`; extend auth-file
patch synchronization so OAuth JSON group edits update the runtime auth.

The TUI fetches and displays API-key group assignments beside masked downstream
keys, provides editing controls for them, and adds `Group` to OAuth credential
details/edit fields. Configured credential group fields are visible through the
existing management-backed key lists.

## Compatibility And Rollback

The absence of `api-key-groups`, or an empty mapping, is compatibility mode.
Operators can assign groups to credentials first, then enable isolation by
adding at least one downstream-key mapping. Removing all mappings rolls back to
legacy routing without deleting group annotations.
