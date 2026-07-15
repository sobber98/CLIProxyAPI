# Credential Grouping Implementation Plan

1. Add group fields and validation/normalization to configuration and runtime
   auth representations; synthesize the field from every configured credential
   type and from OAuth-file metadata.
2. Extend the config access provider to resolve the authenticated downstream
   key's group and carry it through access metadata into execution options.
3. Add one candidate-group predicate in the auth manager and apply it to legacy,
   scheduler, mixed-provider, retry/failover, and model-discovery paths.
4. Add Management API endpoints and persistence validation for `api-key-groups`.
   Extend configured credential and auth-file update paths to set and synchronize
   credential groups.
5. Extend the TUI client and keys/auth tabs to display and edit downstream and
   OAuth credential group assignments.
6. Update `config.example.yaml` with migration-safe configuration examples.
7. Add focused unit and handler tests for validation, access metadata,
   credential synthesis, selection isolation/failover, model listing,
   Management endpoints, and TUI-facing client behavior.

## Validation

- `gofmt -w` on changed Go files.
- Focused `go test` commands for touched config, access, auth, API management,
  API handler, watcher synthesizer, and TUI packages.
- `go test ./...`.
- `go build -o test-output ./cmd/server && rm test-output`.

## Review Gates And Rollback

- Verify every auth-selection entry point applies the same group predicate;
  specifically inspect legacy, built-in scheduler, mixed-provider, retry, and
  model-discovery routes.
- Verify strict mode is solely controlled by non-empty `api-key-groups` and that
  the group predicate does not alter legacy configurations.
- To roll back production behavior, remove all `api-key-groups` entries; group
  annotations can remain because they are inert in compatibility mode.
