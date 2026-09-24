---
page_title: "Upgrading from v0.3.x - descope Provider"
description: |-
  Upgrade a configuration that manages a project with a v0.3.x version of the provider to the standalone resources introduced in v1.0.
---

# Upgrading from v0.3.x

In v0.3.x, the `descope_project` resource managed the entire project configuration with nested attributes such as
`authentication`, `connectors` and `flows`. Starting with v1.0, `descope_project` only manages the project's `name`,
`environment`, `tags` and `deletion_protection`, and every other part of the project is managed by a standalone
resource that references the project with `project_id`.

The upgrade keeps your existing `descope_project` resource and its state, so it works in a state that manages other
infrastructure as well. It never changes the project's configuration in Descope. What it does require is adding the
standalone resources to your configuration and adopting the existing entities into them with `import` blocks.

-> Terraform 1.5 or later is required for `import` blocks.

## What to expect after upgrading to v1.0

- **Your v0.3.x configuration stops working.** With v1.0, `terraform plan` fails with an `Unsupported argument` error
  for each nested attribute of `descope_project` that no longer exists, such as `authentication` or `connectors`.
  The plan fails before anything is changed in Descope or in your state.
- **Removing the nested attributes is safe.** It doesn't delete or change anything in Descope, since v1.0's
  `descope_project` only ever sends the project's `name`, `environment` and `tags`. Terraform just stops managing the
  rest of the project's configuration until you adopt it with the standalone resources, as described in the steps
  below. Until the first `terraform apply` with v1.0, each plan warns about this.
- **The state can't be used with v0.3.x anymore.** Once you apply with v1.0, a v0.3.x version of the provider refuses
  the state, so upgrade every workspace that shares it at the same time. To roll back, restore the state backup from
  the first step.
- **The project is protected from deletion.** In v1.0, destroying a `descope_project` fails unless you set
  `deletion_protection = false`.

## Upgrade steps

1. **Back up the state.** Run `terraform state pull > backup.tfstate` with your current provider version.

2. **Upgrade the provider and remove the nested attributes.** Change the version constraint to `~> 1.0` and run
   `terraform init -upgrade`. Remove every nested attribute from the `descope_project` resource, keeping only `name`,
   `environment` and `tags`.

3. **Plan and apply.** `terraform plan` should report no changes for the project, along with the warning about the
   configuration that's no longer tracked. Run `terraform apply` to save the upgraded state.

4. **Generate the standalone resources.** Use the [tfexport](https://github.com/descope/terraform-provider-descope/tree/main/tools/tfexport)
   tool to read the project and write the matching resources and `import` blocks, referencing your existing project
   resource:

    ```bash
    tfexport -project P... -out ./generated -project-address descope_project.main
    ```

    Copy the generated files into the directory that has your project resource. If that's a module, add
    `-import-prefix module.<name>` and move `import.tf` to the root module, since Terraform only allows `import`
    blocks there. The variables in the generated `variables.tf` then become inputs of that module, so pass their
    values in the `module` block, and rename any that clash with inputs the module already has. Alternatively, write
    the resources yourself using the table below.

5. **Supply secrets.** The Descope API never returns secrets such as connector credentials, so the generated
   configuration declares them as sensitive variables. Provide their values, e.g. in a `terraform.tfvars` file that's
   kept out of version control.

6. **Plan and review.** `terraform plan` should import every resource and change nothing but secrets:

    ```
    Plan: 42 to import, 0 to add, 3 to change, 0 to destroy.
    ```

    The API reports that a secret is stored but not its value, so a sensitive update per imported secret is expected:
    every resource with a stored secret plans an in-place update that sets it to the value of its variable, shown as
    `(sensitive value)`, and applying it resends the secret.

    Don't apply a plan that creates a resource that already exists in the project, or that changes a value you've
    configured or clears a secret, which shows as a sensitive value that changes to `null` or an empty string. An
    update that only sets an attribute that has no stored value to its default, such as `use_mtls` and
    `use_static_ips` on HTTP connectors created with v0.3.x, is safe.

7. **Apply** to adopt the resources, then delete the `import` blocks.

Steps 3 and 6 can be merged into a single apply, by doing step 2 and steps 4 and 5 together. Doing them separately is
recommended, so that the first plan confirms the upgrade itself leaves the project unchanged.

~> Always import the settings resources, `descope_styles`, `descope_fga_schema`, flows and widgets. Creating them
without an import overwrites the existing configuration in the project with the configuration in the resource. The
provider warns during planning when one of these resources is about to be created in an existing project.

## Writing the resources by hand

If you write the standalone resources yourself instead of generating them with tfexport, keep in mind:

- References that used names now use IDs, e.g. `user_jwt_template` in `descope_session_settings` and the connector
  used by a messaging method's settings resource.
- The role and permission IDs in the v0.3.x state aren't the IDs that `descope_role` and `descope_permission` import
  with, so look them up by name with the management API.
- Every project has built-in entities, such as the default OIDC application, JWT templates and flows. Import them
  rather than creating them.
- `descope_inbound_app`, which was already a standalone resource in v0.3.x, replaces `non_confidential_client` with
  `client_type`.
- `descope_outbound_app` is new in v1.0 and has no v0.3.x counterpart.

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
