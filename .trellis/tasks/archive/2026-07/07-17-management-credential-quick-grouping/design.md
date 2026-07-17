# Management Credential Quick Grouping Design

## Boundary

Implement the Management Center change in the user's fork at
`/opt/Cli-Proxy-API-Management-Center`. The Management API supplies the needed
routes, with two compatibility remediations in CLIProxyAPI: auth-file listings
expose persisted OAuth groups in both runtime and disk-fallback responses, and
xAI API-key patches accept the existing `group` field contract. No
credential-routing or credential-selection change is required.

Add a route and sidebar entry for a dedicated credential-grouping page. The
page is the single grouping surface but has two independent batch sections:

1. **Upstream credentials** list all configured provider API-key credentials
   and OAuth auth files. It also expands OpenAI-compatible providers that have
   `api-key-entries` into one selectable item per entry; a provider without
   entries remains one selectable credential.
2. **Downstream API keys** list proxy API keys and their `api-key-groups`
   assignments.

Each section provides filters, per-row selection, select-all for its visible
rows, a group-name input, and actions to assign that name or clear selected
assignments. The sections deliberately submit independently: upstream types
use different Management resources, while downstream keys require one complete
mapping replacement.

## Data And Update Flow

1. The page concurrently loads auth files, configured credential resources,
   downstream API keys, and API-key groups. It normalizes them into UI records
   containing a stable selection identity, display type, masked label, group,
   and an update descriptor.
2. A downstream bulk assignment starts from a freshly fetched
   `api-key-groups` mapping, changes only selected key entries, and replaces it
   with `PUT /api-key-groups`. Clearing removes the selected mapping rather
   than persisting an empty group value.
3. OAuth-file updates use one `PATCH /auth-files/fields` request per selected
   filename with `{ group }`, including `group: ""` to clear.
4. Configured Gemini, Interactions, Claude, Codex, xAI, and Vertex credentials
   use their existing resource-specific PATCH endpoint, locating the entry by
   existing identity and sending only `{ group }` in `value`.
5. An OpenAI-compatible provider with no `api-key-entries` uses its existing
   provider PATCH with `{ group }`. For providers with entries, the page first
   re-fetches the provider list, updates group values only on selected matching
   `api-key-entries`, then PATCHes that provider with the current complete
   entry list. This preserves unrelated entry edits that occurred after page
   load.
6. All selected updates use `Promise.allSettled`. The page reports failed
   records by label and refreshes its data after the batch so the displayed
   state reflects persisted values, including partial success.

## Compatibility And UX

Group values are sent unchanged except for trimming at the existing server
boundary. Clearing uses an empty group as the existing API contract requires.
The control panel remains usable against older servers: unavailable grouping
endpoints or unsupported `group` patches surface as a visible batch failure;
the task does not add API-version probing or alter server fallback behavior.

All user-visible text is added to each existing locale file. Reuse the
project's established layout, selection controls, notification components, and
responsive SCSS patterns. The page must remain practical on narrow screens by
keeping each section independently scrollable and its batch action controls
accessible.

## Delivery And Rollback

Build the fork into its single `dist/index.html` output and release it as
`management.html`. Configure deployments that should use the fork with
`remote-management.panel-github-repository` pointing to
`https://github.com/sobber98/Cli-Proxy-API-Management-Center`. Rollback is a
panel release rollback or restoring the default upstream panel repository; no
server data migration is involved.
