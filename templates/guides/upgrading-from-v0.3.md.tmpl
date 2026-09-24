---
page_title: "Upgrading from v0.3.x - descope Provider"
description: |-
  Upgrade a configuration that manages a project with a v0.3.x version of the provider to the standalone resources introduced in v1.0.
---

# Upgrading from v0.3.x

In v0.3.x, a single `descope_project` resource managed your whole project through nested attributes like
`authentication`, `connectors` and `flows`. In v1.0, each of those parts has its own standalone resource that points
at the project with `project_id`, and `descope_project` only manages the project's `name`, `environment` and `tags`.

The upgrade happens in two parts:

1. **Upgrade the provider.** Bump the version and trim `descope_project` down to the attributes it still supports.
2. **Adopt the rest of the project.** Generate the standalone resources with the `tfexport` tool and import them.

The provider upgrade doesn't change anything in Descope, and the adoption plan shows exactly what it will change. You
also keep your existing `descope_project` resource and its state, so this works fine in a state that manages other
infrastructure too.

## Before you start

- You need Terraform 1.5 or later, for `import` blocks.
- If more than one workspace shares the same state, upgrade them all together. Once you apply with v1.0, previous
  versions of the provider can't read the state anymore.

## Part 1: Upgrading the provider

1. **Back up the state.** With your current provider version, run:

   ```shell
   terraform state pull > backup.tfstate
   ```

   You'll only need this if you want to [roll back](#rolling-back).

2. **Switch to v1.0 and trim the project resource.** Change the provider's version constraint to `~> 1.0` and run
   `terraform init -upgrade`. Then remove every nested attribute from `descope_project`, so that only `name`,
   `environment` and `tags` are left. If you use `descope_inbound_app`, also remove its `non_confidential_client`
   attribute, which v1.0 no longer has, and see [Known issues](#known-issues).

   Removing the nested attributes doesn't delete anything in Descope. Terraform just stops tracking those settings
   until you adopt them again in Part 2.

   If you run `terraform plan` before removing them, you'll get an `Unsupported argument` error for each one. That's
   expected, and nothing is changed in Descope or in your state.

3. **Plan and apply.** `terraform plan` should show no changes to the project, along with a warning about the
   configuration Terraform no longer tracks. Run `terraform apply` to save the upgraded state.

## Part 2: Adopting the rest of the project

1. **Generate the standalone resources.** The [tfexport](https://github.com/descope/terraform-provider-descope/tree/main/tools/tfexport)
   tool reads your project and writes the matching resources, plus `import` blocks that adopt what's already there.
   Its README explains how to install it and which management key it needs. Pass it the address of your existing
   project resource so the generated resources reference it:

   ```shell
   tfexport -project P... -out ./generated -project-address descope_project.main
   ```

   Copy the generated files into the directory that has your `descope_project` resource. If your configuration
   already has `descope_inbound_app` resources, see [Known issues](#known-issues) first.

   If your `descope_project` resource is inside a module, or you manage several projects in one configuration, see
   [Special cases](#special-cases).

2. **Fill in secrets.** The Descope API never returns secrets like connector credentials, so the generated
   configuration declares them as sensitive variables. Set their values, for example in a `terraform.tfvars` file that
   you keep out of version control.

3. **Plan and review.** `terraform plan` should import everything and change nothing but secrets:

   ```
   Plan: 42 to import, 0 to add, 3 to change, 0 to destroy.
   ```

   The in-place changes are normal. The API tells the provider that a secret is stored, but not its value, so every
   imported resource with a stored secret plans an update that sets it from its variable. It shows up as
   `(sensitive value)`, and applying it just resends the secret.

   Don't apply the plan if it:

   - creates a resource that already exists in the project
   - changes a value you've configured
   - clears a secret, which shows as a sensitive value changing to `null` or an empty string

   See [Known issues](#known-issues) for the other in-place changes you might see.

4. **Apply and clean up.** Run `terraform apply` to adopt the resources, then delete the `import` blocks.

That's it. Your project is managed by Terraform again.

**Important:** In v1.0 the project is protected from deletion. Destroying a `descope_project` fails unless you set
`deletion_protection = false`.

## Known issues

### Inbound apps can end up managed twice

`descope_inbound_app` was already a standalone resource in v0.3.x, and tfexport exports every inbound app in the
project. If your configuration already manages an inbound app, delete the generated resource for it and its `import`
block, or two resources will manage the same app. tfexport prints a warning for each inbound app it exports as a
reminder.

### Adding `client_type` to an existing inbound app replaces it

v1.0 replaces the `non_confidential_client` attribute of `descope_inbound_app` with `client_type`. Don't add
`client_type` to an app that already exists: changing it replaces the app, which gives it a new client ID and secret.
Removing `non_confidential_client` on its own plans no changes.

### HTTP connectors created with v0.3.x plan an update

The adoption plan shows an in-place update that sets `use_mtls` and `use_static_ips` to `false` on HTTP connectors that
were created with v0.3.x, because those connectors have no stored value for them. The update is safe: `false` is how
the connector already behaves.

## Special cases

### The project resource is inside a module

Terraform only allows `import` blocks in the root module, so run tfexport with `-import-prefix module.<name>`, copy the
generated resources into the module, and move `import.tf` to the root module.

The variables in the generated `variables.tf` become inputs of the module, so pass their values in the `module` block.
Rename any that clash with inputs the module already has.

### You manage several projects

If your configuration manages several projects, for example with `for_each` or one `descope_project` resource per
environment, run tfexport once per project with a different `-name-prefix`. The prefix is added to every generated
resource, variable and file name, so the exports don't clash and can go in the same directory:

```shell
tfexport -project P... -out ./generated-dev -project-address 'descope_project.main["dev"]' -name-prefix dev
tfexport -project P... -out ./generated-prod -project-address 'descope_project.main["prod"]' -name-prefix prod
```

### Upgrading with a single apply

You can skip the apply in Part 1 by making all the configuration changes first (trimming `descope_project`, generating
the resources and filling in secrets) and then running one plan and apply.

We recommend doing it in two steps anyway. The first plan confirms that the provider upgrade by itself leaves the
project unchanged, which makes the second, bigger plan easier to review.

## Writing the resources by hand

If you'd rather write the standalone resources yourself, use the [attribute mapping](#attribute-mapping) below to find
the resources and import IDs that replace each part of your v0.3.x configuration.

**Note:** Always import the settings resources, `descope_styles`, `descope_fga_schema`, flows and widgets. Creating one
of them without an import overwrites the existing configuration in the project with the configuration in the resource.
The provider warns during planning when one of these is about to be created in an existing project.

A few more things to watch for:

- References that used names now use IDs, e.g. `user_jwt_template` in `descope_session_settings` and the connector used
  by a messaging method's settings resource.
- Every project has built-in entities that Descope creates automatically, such as the default OIDC application, JWT
  templates and flows. You only need them in your configuration if you actually use them. If you do, import them with
  their existing IDs (see the Import ID column below) rather than creating them.
- The role and permission IDs in the v0.3.x state aren't the IDs that `descope_role` and `descope_permission` import
  with, so look them up by name with the management API.
- `descope_outbound_app` is new in v1.0 and has no v0.3.x counterpart.

## Rolling back

Until your first `terraform apply` with v1.0, you can roll back by reverting the version constraint and your
configuration changes, then running `terraform init -upgrade`. After that apply, v0.3.x refuses the upgraded state, so
you'll also need to restore the backup from Part 1:

```shell
terraform state push -force backup.tfstate
```

Pushing the backup replaces the whole state, and `-force` skips the check that would normally stop an older state from
overwriting a newer one.

**Important:** Only do this if nothing else was applied to the state since you took the backup, including
changes to resources that aren't managed by Descope. Otherwise, contact Descope support.

## Attribute mapping

| v0.3.x attribute | Standalone resources | Import ID |
|---|---|---|
| `project_settings` | `descope_project_settings`, `descope_session_settings`, `descope_session_migration` | `<project_id>` |
| `invite_settings` | `descope_invite_settings`, `descope_email_template` | `<project_id>`, `<project_id>/invite/<template_id>` |
| `authentication.otp` | `descope_otp_settings`, `descope_email_template`, `descope_text_template`, `descope_voice_template` | `<project_id>`, `<project_id>/otp/<template_id>` |
| `authentication.magic_link` | `descope_magiclink_settings`, `descope_email_template`, `descope_text_template` | `<project_id>`, `<project_id>/magiclink/<template_id>` |
| `authentication.enchanted_link` | `descope_enchantedlink_settings`, `descope_email_template` | `<project_id>`, `<project_id>/enchantedlink/<template_id>` |
| `authentication.embedded_link` | `descope_embeddedlink_settings` | `<project_id>` |
| `authentication.password` | `descope_password_settings`, `descope_email_template` | `<project_id>`, `<project_id>/password/<template_id>` |
| `authentication.sso` | `descope_sso_settings`, `descope_email_template` | `<project_id>`, `<project_id>/sso/<template_id>` |
| `authentication.totp` | `descope_totp_settings` | `<project_id>` |
| `authentication.passkeys` | `descope_passkey_settings` | `<project_id>` |
| `authentication.oauth` | `descope_oauth_settings`, `descope_oauth_provider` | `<project_id>`, `<project_id>/<provider_name>` |
| `authorization.roles`, `authorization.permissions` | `descope_role`, `descope_permission` | `<project_id>/<id>` |
| `authorization.fga` | `descope_fga_schema` | `<project_id>` |
| `attributes` | `descope_user_attribute`, `descope_tenant_attribute`, `descope_access_key_attribute` | `<project_id>/<id>` |
| `connectors.<type>` | `descope_<type>_connector` | `<project_id>/<id>` |
| `applications` | `descope_oidc_app`, `descope_saml_app`, `descope_wsfed_app` | `<project_id>/<id>` |
| `applications.*.roles`, `applications.*.permissions` | `descope_app_role`, `descope_app_permission` | `<project_id>/<app_id>/<id>` |
| `jwt_templates` | `descope_jwt_template` | `<project_id>/<id>` |
| `styles` | `descope_styles` | `<project_id>` |
| `flows` | `descope_flow` | `<project_id>/<flow_id>` |
| `widgets` | `descope_widget` | `<project_id>/<widget_id>` |
| `lists` | `descope_list` | `<project_id>/<id>` |
| `admin_portal` | `descope_admin_portal` | `<project_id>` |
