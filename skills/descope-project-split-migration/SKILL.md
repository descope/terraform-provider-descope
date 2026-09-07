---
name: descope-project-split-migration
description: Human-checkpointed workflow for safely splitting legacy descope_project state into standalone resources. Use when a developer says "migrate descope_project", "resource split migration", "split the project resource", "adopt standalone descope resources", or "descope_project state migration".
---

# Descope Project Split Migration

Terraform cannot fan one resource's state into many destinations through `UpgradeState`, `MoveState`, or `moved` blocks. This migration therefore uses native adoption: capture legacy state -> build an ownership and ID manifest -> detach the legacy `descope_project` from state without deleting anything remotely by using a `removed` block with `lifecycle { destroy = false }` -> adopt every destination with native `import` blocks. The adoption plan must then pass a non-destructive reconciliation check. The `tfmigrate` utility is the single source of migration logic and generated artifacts.

## Install

Install this skill from the public provider repository into any supported agent:

```bash
npx skills add descope/terraform-provider-descope --skill descope-project-split-migration
```

The Skills CLI copies this directory into the agent's skill store. A checkout of the provider repository is not required: the workflow runs the migration utility from the exact tagged target-provider release.

For one-time use without installation:

```bash
npx skills use descope/terraform-provider-descope@descope-project-split-migration
```

## Scope

Use this skill only to migrate ownership from a legacy `descope_project` instance to the provider's standalone Descope resources. The workflow inventories existing state, generates Terraform configuration and review material, separates detach from adoption, and validates the saved adoption plan.

Hard constraints:

- Terraform must be `>= 1.8.0`. Terraform 1.7 introduced `removed` blocks and `for_each` in import blocks; 1.8 adds the plan `complete` signal that `verify-plan` requires to reject targeted or deferred plans.
- Keep the source provider version pinned throughout the detach stage.
- Keep the target provider version pinned throughout the adopt stage.
- Run `tfmigrate` from the exact Git tag matching the target provider version. Never use `@latest`, a branch, or an unpinned local checkout for a customer migration.
- Both exact provider versions are recorded in `manifest.json` under `source_provider.version` and `target_provider.version`.
- The utility owns all mapping, destination addressing, ID derivation, artifact generation, and plan analysis. Invoke it; never reimplement any of that logic in shell, HCL, or agent reasoning.
- Resource instance addresses with root-module `count` or `for_each` are supported when the declaration has one live instance. An indexed module address fails closed because child-module HCL applies to every module instance and cannot safely be generated from one instance's state.

## When to Use

- A developer asks to "migrate descope_project" or perform a "descope_project state migration".
- A team needs a "resource split migration" or wants to "split the project resource".
- Existing nested project configuration must "adopt standalone descope resources" without recreating remote objects.

## Prepare the Utility

Run the utility from the released provider module, not from a provider source checkout. Set the exact versions once. Provider version values omit the `v` prefix; Go module tags include it:

```bash
SOURCE_PROVIDER_VERSION="<X.Y.Z>"
TARGET_PROVIDER_VERSION="<A.B.C>"
TFMIGRATE_PACKAGE="github.com/descope/terraform-provider-descope/tools/tfmigrate@v${TARGET_PROVIDER_VERSION}"
```

The target release must contain `tools/tfmigrate`. `TFMIGRATE_PACKAGE` must remain unchanged for the entire migration. Before reading state, confirm the pinned utility is available:

```bash
go run "$TFMIGRATE_PACKAGE" help
```

Do not substitute `@latest`, a branch name, or a locally modified checkout.

## Command Reference

| Command | Purpose |
|---------|---------|
| `go run "$TFMIGRATE_PACKAGE" plan -state <state.json> -address <legacy address> -source-provider-version "$SOURCE_PROVIDER_VERSION" -target-provider-version "$TARGET_PROVIDER_VERSION" -out <dir>` | Inventory the legacy instance and generate the manifest, review documents, detach configuration, and adoption configuration. |
| `go run "$TFMIGRATE_PACKAGE" plan -state <state.json> -address <legacy address> -source-provider-version "$SOURCE_PROVIDER_VERSION" -target-provider-version "$TARGET_PROVIDER_VERSION" -out <dir> -remote` | Retry only when the tool reports IDs missing from state. Performs GET-only reads and requires `DESCOPE_MANAGEMENT_KEY`; `DESCOPE_BASE_URL` is optional. |
| `go run "$TFMIGRATE_PACKAGE" verify-plan -plan <plan.json> -manifest <dir>/manifest.json` | Verify that the saved adoption plan is complete and non-destructive. |

| Exit code | Meaning |
|-----------|---------|
| `0` | Success. |
| `1` | Migration blockers or plan verification failure. |
| `2` | Usage or I/O error. |

## Workflow Steps

### Step 1: Freeze Writers

Announce a change freeze for both the Terraform workspace and the Descope console for the affected project. Confirm that no CI pipeline, scheduled job, or other operator can apply Terraform during the migration. Do not continue until the freeze is active.

### Step 2: Capture Protected Backups

Capture the current state into a new owner-only file before generating anything:

```bash
(umask 077; set -o noclobber; terraform state pull > legacy.tfstate.json)
```

Raw state records the provider configuration address. `terraform show -json` omits aliases, so when it is used for inventory, explicitly confirm the module's default `descope` provider is the same configuration that owned the legacy resource.

`terraform show -json` output is also accepted as inventory input, but it is not a restorable state backup: keep the raw `terraform state pull` file as well. Both files are SENSITIVE and may contain connector secrets, access keys, and secret variable values. Store them only where the team keeps secrets, restrict access according to team policy, and never commit them.

### Step 3: Run Inventory and Generation

Run the `plan` subcommand with the exact legacy resource address and the exact source and target provider versions:

```bash
go run "$TFMIGRATE_PACKAGE" plan -state legacy.tfstate.json -address <legacy address> -source-provider-version "$SOURCE_PROVIDER_VERSION" -target-provider-version "$TARGET_PROVIDER_VERSION" -out <dir>
```

`<dir>` must not exist. The utility atomically creates it with owner-only permissions and refuses existing paths or symlinks so the protected state copy cannot be redirected or mixed with older output.

If and only if the utility reports IDs missing from state, provide `DESCOPE_MANAGEMENT_KEY`, optionally set `DESCOPE_BASE_URL`, and rerun with `-remote`. Remote backfill performs one GET-only project read. A custom base URL must use HTTPS, a DNS hostname, and the default TLS port; the utility rejects HTTP, IP literals, credentials in URLs, custom ports, and `TF_UNSAFE_LOGS`.

Review these generated artifacts before changing Terraform state:

- `manifest.json`
- `README.md`
- `SECRETS.md`
- `ROLLBACK.md`

Treat the manifest as the reviewable ownership and ID map. Confirm its legacy address, provider pins, entries, secrets, retained fields, unowned objects, and rollback metadata match the intended project.

### Step 4: Surface and Resolve Blockers

If the utility exits `1`, read `BLOCKERS.md`. Every blocker is fatal. Stop the migration until the underlying issue is resolved and a rerun exits `0`.

| Blocker kind | Required human action |
|--------------|-----------------------|
| `missing_id` | Locate the real backend ID from authoritative state or use the utility's GET-only `-remote` mode when offered. Never guess it. |
| `unsupported_legacy_field` | Review the populated legacy field with the owning team and update the migration logic or source configuration through an approved change before rerunning. |
| `duplicate_destination` | Determine why multiple legacy objects resolve to one destination address and correct the ownership or source configuration. |
| `ambiguous_ownership` | Establish whether Terraform, Descope, or another system owns the object and make that ownership unambiguous. |
| `unresolved_reference` | Identify the actual referenced remote object and repair the source reference or migration mapping. |
| `missing_payload` | Recover the required payload from an authoritative protected source or fix generation logic; do not synthesize content. |
| `unresolved_secret` | Obtain or rotate the real required secret and arrange to supply it from the operator's secret store. Never invent a value. |
| `unknown_legacy_attribute` | Review the populated attribute against the source provider schema and extend the utility's explicit handling before migration. |
| `invalid_legacy_value` | Correct the invalid source configuration through a supported workflow or update the utility to handle the valid legacy representation. |

Never invent an ID, payload, ownership decision, reference, or secret to clear a blocker.

### Step 5: HUMAN CHECKPOINT 1 - Approve Detach

Explicit operator approval is required before detach. Display the exact generated `removed` block from `detach/removed.tf` verbatim. Require the operator to confirm that the block targets the expected legacy local address and visibly contains `lifecycle { destroy = false }`. Do not reconstruct or alter the generated block.

### Step 6: Detach With the Old Provider Pinned

Keep the old source provider version pinned. Preserve the legacy HCL for review, remove the legacy `resource "descope_project"` block from the module, and put `detach/removed.tf` in its place. Then run:

```bash
terraform plan
```

Require the only planned change to be the single expected `forget` of the legacy address. Any create, update, delete, replacement, read, or additional forget stops the workflow. After the operator reviews the plan, the operator may run `terraform apply`.

The agent must never run `terraform apply`, `terraform state rm`, or `terraform state mv` automatically.

### Step 7: HUMAN CHECKPOINT 2 - Approve Adoption

Explicit operator approval is required before adoption. Show the reviewed manifest, generated adoption files, secret requirements, and detach result. Do not initialize or plan the adopt stage until the operator approves proceeding with the target provider version.

### Step 8: Adopt With the New Provider Pinned

Pin the new target provider version. Remove `detach/removed.tf`, then place these generated artifacts in the configuration:

- `adopt/main.tf`
- `adopt/imports.tf`
- `adopt/variables.tf`
- `adopt/versions.tf`
- `adopt/payloads/`

`adopt/main.tf`, `adopt/variables.tf`, and `adopt/payloads/` belong in the module that declared the legacy resource. `adopt/imports.tf` belongs in the root module because import blocks are root-only. Place `adopt/versions.tf` as directed by the generated `README.md` and preserve the relative payload paths.

Supply the sensitive variables listed in `SECRETS.md` from the operator's own secret store. For a root resource, `TF_VAR_<name>` is acceptable. For a child module, `TF_VAR_*` only populates root inputs: add matching sensitive variables at the root and every parent module, then pass each value through the existing module call chain exactly as `SECRETS.md` directs. Never place secret values in module blocks or a committed `.tfvars` file. Then run:

```bash
terraform init -upgrade
(umask 077; terraform plan -out=adopt.tfplan)
(umask 077; set -o noclobber; terraform show -json adopt.tfplan > adopt-plan.json)
go run "$TFMIGRATE_PACKAGE" verify-plan -plan adopt-plan.json -manifest <dir>/manifest.json
```

Both `adopt.tfplan` and `adopt-plan.json` contain sensitive values. Never commit or upload them; delete them according to the team's retention policy after verification and apply.

### Step 9: Verify the Final Plan Is Non-Destructive

`verify-plan` must exit `0`. An acceptable adoption plan is complete and contains only manifest-enumerated imports with no-op actions, unrelated already-converged no-op resources, and imported resources updated only for the exact secret attributes enumerated for their manifest entries. Missing or mismatched import IDs, unexpected imports, deferred changes, action invocations, incomplete or errored plans, creates, deletes, reads, forgets, replacements, updates without imports, and updates to any non-enumerated attribute are failures.

Only the operator may apply after `verify-plan` passes and after reviewing the final plan. The agent must not apply it automatically.

### Step 10: Retain Rollback Material

Keep the protected raw state backup and `ROLLBACK.md` for the team's agreed retention window. Rollback requires another explicit human checkpoint: verify workspace, lineage, and serial before a human restores state and re-pins the old provider. The agent and utility must never run `terraform state push`. Restoring state does not undo a write-only secret rehydration update; separately restore or rotate every such credential named by the adoption plan. No remote objects are created or deleted during detach or adoption.

## Provider-Side Guard

The post-split `descope_project` schema is version `1`. Its state upgrader fails closed when version `0` state still carries populated legacy nested ownership, and emits the diagnostic summary `Descope Project State Requires Migration` while listing the populated attributes. This is a guard, not a migration: the Terraform Plugin Framework cannot fan one resource into many resources, so the upgrader can only pass clean state through or refuse unsafe legacy state.

## Never

- Never run `terraform apply`, `terraform state rm`, or `terraform state mv` automatically.
- Never run `terraform state push` automatically; rollback state restoration is human-only after workspace, lineage, and serial checks.
- Never invent IDs, payloads, references, ownership decisions, or secrets.
- Never commit the state backup or generated secret values.
- Never skip `verify-plan`.
- Never use `-target` to work around a failing plan.
- Never admin-merge or bypass branch protection.
