# Management Credential Quick Grouping Implementation Plan

## Implementation

1. In CLIProxyAPI, expose OAuth auth-file groups consistently in list responses
   and support xAI `group` PATCH assignment and clearing, with focused
   management regressions. In the Management Center fork, extend API types, normalizers, and API
   services so credential groups are retained from existing responses and each
   existing Management resource can receive a minimal group update. Add API-key
   group mapping read/write support.
2. Add a reusable grouping data adapter that loads configured credential lists,
   auth files, and downstream API keys; normalizes their identities and labels;
   and performs the resource-specific batch updates. Preserve OpenAI-compatible
   API-key entry changes by re-fetching before its whole-entry-list PATCH.
3. Implement the responsive credential-grouping page with independent upstream
   and downstream sections, filtering, multi-selection, group assignment,
   ungrouping, batch progress, partial-failure reporting, and refresh after
   submission.
4. Register the page route and sidebar item, add locale strings to all current
   locale files, and update documentation to describe the new grouping page and
   the fork configuration required to distribute it.
5. Add focused tests for normalization and updates, including downstream map
   replacement, clearing assignments, mixed batch success/failure, and
   OpenAI-compatible entry refresh-before-update behavior.

## Validation

Run in `/opt/Cli-Proxy-API-Management-Center`:

```bash
bun run test
bun run lint
bun run type-check
bun run build
```

Manually verify the built single-file panel against a CLIProxyAPI instance with
at least one credential from every supported configured provider type, OAuth
auth files, OpenAI-compatible API-key entries, and multiple downstream API
keys. Verify assignment, clearing, reload persistence, partial error display,
desktop layout, and narrow-screen layout.

## Risk And Rollback

- The risky update is OpenAI-compatible API-key entry replacement because its
  PATCH accepts a full entry list. Re-fetch immediately before update and merge
  only selected entries.
- If a release causes a UI regression, restore the previous `management.html`
  release or point `panel-github-repository` back to the upstream repository.
- No CLIProxyAPI source or persisted configuration schema changes are planned.
