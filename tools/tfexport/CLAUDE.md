# tfexport

Generates a Terraform configuration for an existing Descope project: resource blocks for the values a human
would have written, import blocks so the configuration adopts the live entities, and large payloads in their
own files. It only reads and generates, and never changes customer data.

## Pipeline

`discover -> read -> prune -> emit`, wired together in `internal/export`.

**discover** builds the work list. The project snapshot is an index of what exists, not a source of values.
Three things about it are easy to get wrong:

- Not everything is in it. Inbound apps have their own load-all endpoint.
- It rewrites role and permission ids for portability, so reading one back by its snapshot id fails and the
  entity is silently lost. Both are name keyed, so live ids are resolved by name instead. App scoped entries
  keep their application ids.
- System OAuth providers always exist in it with built-in values, so they are read like any other provider
  but only exported when pruning leaves something besides their id and a redirect_url equal to the one most
  system providers share.

**read** loads each entity through the provider's own export hooks, the same path a `terraform import`
refresh takes, so the tool never grows a second opinion about what an attribute means.

**prune** cuts the values down to what a person would have typed. **emit** renders HCL with `hclwrite`.

## Pruning

Dropped: computed-only attributes, unset ones, and anything equal to its schema default. Three comparisons
are not the obvious one:

- **Custom value types.** The models' types have a type-strict `Equal`, so identical content compares as
  different against a plain basetypes default. Compare at the `tftypes` level.
- **Durations.** `"1 weeks"` and `"1 week"` are the same length of time. `durationattr.Set` deliberately
  keeps whichever spelling the state holds, so compare seconds, not text.
- **Server zero values.** The backend returns `""` or `[]` rather than omitting a field. On an
  Optional+Computed attribute with no declared default, dropping one keeps the stored value instead of
  planning a change. Required attributes are excluded: an empty value there is the configuration.

An object whose fields all match their nested defaults but which differs from its own declared default is
emitted verbatim, otherwise omitting it would plan a change. Elements of a nested collection are never
dropped, only their fields, so an all-default element still appears as an empty block.

## Secrets

Secret values never reach a generated file. A required secret becomes a variable, and so does an optional
one with a declared null default, since omitting that would erase the stored value on the next apply. One
with no declared default cannot be expressed at all, so it is dropped and the run warns.

An import read records a stored connector secret as the backend's placeholder, which is how prune learns it
exists. A secret map becomes one variable per key, since variables are strings. The placeholder itself must
never reach a generated file, and the round-trip integrity check fails if it does.

When pruning leaves a configuration that fails the model's cross-field validation, `ensureValidConfig`
promotes the smallest set of omitted secrets that makes it valid. Paths are dotted, and a secret inside a
collection carries its index (`headers[1].value`) so a promotion lands on one element and each element gets
its own variable.

## Labels and references

Labels come from entity names, fall back to the entity id, and are made unique per resource type with a
numeric suffix. The project and every singleton are labelled after the project itself, so two exports can be
pasted into one configuration without colliding.

Server-assigned ids become references to the resource managing them rather than opaque literals. Short
human-chosen ids are excluded, since a flow id or attribute name would match far too eagerly. Permissions
are additionally mapped by name, because roles grant them by name and terraform needs the ordering.

Import block id formats have to match what `baseResource.ImportState` expects.

## Deliberately not exported

`descope_access_key`, `descope_management_key` and `descope_descoper` are out of scope by design, and
`descope_engine` is deferred. They are listed in `generativeSkippedResources` with their reasons. Removing
one from that list without adding discovery makes `TestGenerativeCoverage` fail with a discovery gap, which
is the intended guard: a registered resource with no discovery path is silently absent from every export.

## Test suites

Everything here needs `TF_ACC=1`, and the round-trip battery additionally needs
`DESCOPE_TFEXPORT_ROUNDTRIP=1` because it creates and destroys real projects. A green offline suite proves
very little: most of the fixture surface only executes under those flags.

- **`TestRoundTrip`** seeds a project from `testdata/scenarios`, exports it, applies the export, and
  requires the next plan to be a no-op. Each step also checks export integrity, formatting, determinism,
  and that every seeded value survived into the export.
- **`TestGenerativeCoverage`** generates one of every registered resource type and reports attribute
  coverage. It is what catches a resource that has no discovery path.
- **`TestExportReplication`** applies one project's export into a second existing project. The project block
  is stripped and its references retargeted, since the target already exists.
- **`TestExportDayTwo`** re-exports into a live workspace after the project evolves.

Two retries in the harness are deliberate: destroys are retried because parallel connector deletion fails
spuriously and each pass deletes more, and the exported apply is retried once because concurrent settings
writes fail on backend version conflicts.

`componentsVersion` is stamped by the rendering service, not by the exporter: orchestrationservice merges
its hardcoded default theme version into every read, so two reads of an unchanged project can legitimately
disagree. It is stripped before determinism comparison, at the top level of a styles payload and in the
metadata of a flow or widget, and nowhere else.

`attributesExemptFromConsistency` holds values whose seed-state and import-read forms legitimately differ,
each key scoping its reason to one resource type, attribute and JSON path. `knownInconsistencies` is
different: those are confirmed open bugs, logged rather than failed, and each entry must go when it is
fixed.
