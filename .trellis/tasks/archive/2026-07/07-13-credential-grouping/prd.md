# Credential grouping

## Goal

Allow operators to partition upstream credentials into groups and bind each
downstream proxy API key to one group, so a request authenticated with that key
can only select credentials from its assigned group and only discovers models
available through that group.

## Confirmed Facts

- Downstream proxy authentication is currently configured as the top-level
  `api-keys: []string` field in `internal/config/sdk_config.go`.
- Credential selection is centralized in `sdk/cliproxy/auth.Manager`; both the
  legacy and scheduler paths select from the manager's runtime `Auth` records.
- Runtime credentials already have a `Prefix`, `Attributes`, and `Metadata`, but
  no grouping field. Configured API-key credentials use per-provider config
  entries, while OAuth/file-backed credentials are loaded into the same runtime
  representation. File-backed metadata already supports fields such as
  `prefix`, `proxy_url`, `disabled`, and `priority`.
- `force-model-prefix` controls model-prefix routing only. A model prefix is not
  an authorization boundary for a downstream client API key.
- The authentication middleware stores the authenticated downstream API key and
  access metadata in Gin context. Request execution metadata is constructed
  from that context before the auth manager selects a credential.
- Management provides list-oriented API-key endpoints and a full YAML editor.
  The TUI supports inline editing of API keys and selected OAuth-file metadata;
  configured upstream API-key lists are managed through dedicated endpoints.

## Requirements

- A downstream API key must constrain all credential selection and retry/failover
  for its request to the credential group assigned to that key.
- Credentials outside the assigned group must never be selected for that request.
- Existing ungrouped `api-keys` and credentials must retain their current
  behavior unless the operator opts into grouping.
- The configuration example and management/config behavior must document and
  preserve the supported grouping configuration.
- Credential grouping applies to configured API-key credentials and file-backed
  OAuth credentials.
- Once grouping is configured, an ungrouped downstream API key can only select
  ungrouped credentials. It must not access any named group.
- Downstream key assignment retains the existing `api-keys: []string` format and
  uses a separate `api-key-groups: map[string]string` mapping from API key to a
  single credential group.
- Configuration loading and hot reload must reject an `api-key-groups` entry
  whose API key is not present in `api-keys`.
- Provide dedicated Management API endpoints and TUI controls to view and edit
  downstream API key group assignments, plus group fields for configured and
  file-backed credentials.
- Group names are trimmed, case-sensitive identifiers. An empty group value
  means ungrouped.
- Strict isolation activates only when `api-key-groups` is non-empty. Before
  that, grouping credentials alone must not change legacy routing behavior.
- When strict isolation is active, model-discovery endpoints must return only
  models callable through credentials in the authenticated API key's group.

## Acceptance Criteria

- [ ] A request using an API key assigned to a group only selects credentials in
  that group, including after retry or failover.
- [ ] A request cannot select a credential belonging to another group.
- [ ] Existing configurations without groups behave as before.
- [ ] When `api-key-groups` is empty, credentials with a group retain legacy
  selection behavior.
- [ ] When `api-key-groups` is non-empty, an unassigned API key can select only
  ungrouped credentials.
- [ ] Configuration loading and hot reload reject an API key-group mapping for
  a key not present in `api-keys`.
- [ ] Model discovery for a grouped API key exposes only its group's models.
- [ ] Dedicated Management API and TUI controls can view and edit downstream
  key groups and credential groups for configured and file-backed credentials.
- [ ] Automated tests cover group isolation and the no-match behavior.
- [ ] `config.example.yaml` documents the new configuration.

## Out Of Scope

- Multiple credential groups per downstream API key.
- Case-insensitive group matching or automatic group-name canonicalization.
- Changing the existing `api-keys` YAML representation.
