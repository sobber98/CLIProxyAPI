# Management credential quick grouping

## Goal

Let operators quickly assign credential groups to upstream credentials and
downstream API keys from the Management control panel, rather than requiring
individual API calls or manual configuration edits.

## Confirmed Facts

- `management.html` is an externally maintained control-panel asset. This
  repository downloads it from
  `router-for-me/Cli-Proxy-API-Management-Center` to the local static directory
  (`internal/managementasset/updater.go:28-36,189-191`), so it is not a
  versioned frontend source file in this repository.
- The Management Center implementation target is the user's fork at
  `https://github.com/sobber98/Cli-Proxy-API-Management-Center`.
- The Management API already provides validated replacement of the complete
  downstream API-key-to-group mapping through
  `GET`/`PUT /v0/management/api-key-groups`
  (`internal/api/server.go:862-868`,
  `internal/api/handlers/management/config_lists.go:226-250`).
- Configured upstream API-key Management resources accept a `group` field in
  their existing patch payloads, while OAuth auth-file patches synchronize a
  `group` field to the runtime credential
  (`internal/api/handlers/management/config_lists.go:299`,
  `internal/api/handlers/management/auth_files.go:1647-1650`).
- The terminal UI can edit one downstream API-key group at a time, but the
  requested management control-panel workflow is currently unavailable
  (`internal/tui/keys_tab.go:151-163`).

## Requirements

- The control panel must provide an efficient grouping workflow for both
  upstream credentials and downstream API keys.
- Add a dedicated grouping view that shows all upstream credentials and
  downstream API keys organized by their current group.
- Operators can select multiple entries and assign them to one named group or
  clear their assignment in one operation.
- The first release covers every currently supported upstream credential type,
  including OAuth-file credentials and every configured upstream API-key
  resource, plus downstream API keys.
- OpenAI-compatible providers with `api-key-entries` expose each API-key entry
  as a separately selectable credential. Before saving their group changes, the
  control panel must refresh the provider and update only the selected entries'
  `group` fields to avoid overwriting unrelated changes.
- The view must save each selected entry through its applicable existing
  Management API resource, and show failures without treating failed entries as
  successfully updated.
- Existing group semantics remain unchanged: group names are trimmed and
  case-sensitive; an empty group denotes the ungrouped bucket.

## Acceptance Criteria

- [ ] The control panel has a dedicated grouping view that lists upstream
  credentials and downstream API keys with their current group assignments.
- [ ] An operator can select multiple upstream credentials, apply one group or
  clear their groups, and see the resulting assignments after refresh.
- [ ] An operator can select multiple downstream API keys, apply one group or
  clear their groups, and see the resulting assignments after refresh.
- [ ] A failed update is reported and is not presented as a successful change.
- [ ] The change uses the existing Management API contracts and does not alter
  credential-routing or group-isolation semantics.

## Out Of Scope

- Changing credential-routing, group-isolation, or Management API validation
  semantics.
- Adding a new credential type or changing the upstream Management Center
  project outside the user's fork.

## Notes

- The implementation may require a coordinated change in the separately hosted
  Management Center repository.
