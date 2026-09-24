# tfexport

`tfexport` generates a Terraform configuration for an existing Descope project, so you can start managing it with
the [Descope Terraform Provider](https://registry.terraform.io/providers/descope/descope) without writing the
configuration by hand.

It reads your project and writes a directory of `.tf` files that describe it, along with `import` blocks that let
Terraform adopt the existing entities instead of creating new ones. It only reads from your project and never
changes anything in it.

## What it does

- Generates a resource block for every supported entity in the project: settings, authentication methods, roles
  and permissions, connectors, applications, flows, styles, templates, and more.
- Writes the values a person would have written, leaving out anything that matches the provider's defaults.
- Links resources to each other with references instead of hardcoded IDs, so the configuration reads naturally.
- Puts large payloads in their own files, such as flows in `flows/`, widgets in `widgets/`, the theme in
  `styles.json`, and long message template bodies in `templates/`.
- Generates an `import` block for each resource, so the first `terraform apply` adopts your project as it is.
- Never writes secret values. Connector passwords, API keys and similar values become sensitive Terraform variables
  that you fill in yourself, and the tool warns about any it can't represent that way.

## What it doesn't do

- It doesn't export access keys, management keys, console users (`descope_descoper`), or engines. Access keys and
  management keys carry key material, console users are company-level rather than project-level, and engine export
  isn't supported yet.
- It doesn't export users or tenants. The Terraform provider manages project configuration, not the data in it.
- It can't recover secret values, since the Descope API never returns them.

## Requirements

- [Go](https://go.dev/doc/install) 1.26 or later, to install the tool.
- [Terraform](https://developer.hashicorp.com/terraform/install) 1.5 or later, for `import` block support.
- A Descope management key with access to the project. You can create one in the
  [Company section](https://app.descope.com/settings/company) of the Descope console.

## Setup

Install the tool with `go install`, which fetches and builds it without cloning the repository:

```bash
go install github.com/descope/terraform-provider-descope/tools/tfexport@main
```

This puts a `tfexport` binary in `$(go env GOPATH)/bin`, usually `~/go/bin`, so make sure that directory is in your
`PATH`. Run the same command again to update to the latest version.

## Usage

Set your management key and run the tool with the ID of the project to export, which you can find in the
[Project section](https://app.descope.com/settings/project) of the Descope console, and an output directory. The
directory is created if it doesn't exist.

```bash
export DESCOPE_MANAGEMENT_KEY="K..."

tfexport -project P... -out ./my-project
```

When it finishes, it prints how many resources it generated, along with any warnings.

### Flags

| Flag | Description |
|---|---|
| `-project` | The ID of the project to export. Required. |
| `-out` | The directory to write the generated files into. Required, and must be empty or not exist yet. |
| `-force` | Write into the output directory even if it already contains files. Files the export doesn't overwrite are left in place, so prefer a fresh directory. |
| `-only` | Limit the export to resource types whose name contains this text, e.g. `-only connector`. |
| `-project-address` | Reference an existing `descope_project` resource instead of exporting the project, e.g. `-project-address descope_project.main`. The export then has no project block, project import or `provider.tf`. See [Migrating from v0.3.x](#migrating-from-v03x). |
| `-import-prefix` | The module path the generated resources will live in, e.g. `-import-prefix module.auth`. It's prepended to the `to` address of every import block. |
| `-name-prefix` | Prefix every generated resource name, variable name and file name except the shared `provider.tf`, e.g. `-name-prefix prod`, so that exports of several projects can share one directory. |

### Warnings and exit codes

The tool exits with `0` on success and `1` if the export failed. It exits with `2` when the export completed but
some entities couldn't be read and are missing from it, each printed as an `incomplete:` line. The generated
configuration doesn't cover those entities, so check them before relying on it.

Lines starting with `warning:` don't affect the exit code, but read them before applying. Most tell you which
secrets to supply through their generated variables. A warning that a secret has a stored value the configuration
cannot carry means you must set that value manually before applying, otherwise the apply clears it.

## Applying the configuration

1. Provide values for the variables declared in `variables.tf`. These are the secrets the export couldn't write,
   plus `project_id` when the export didn't include the project itself, such as with `-only`. Set them in a
   `terraform.tfvars` file or as `TF_VAR_<name>` environment variables, and keep secrets out of version control.

2. Initialize and review the plan:

    ```bash
    cd my-project
    terraform init
    terraform plan
    ```

    The plan should only import resources and update secrets, and end with a summary like this one:

    ```
    Plan: 42 to import, 0 to add, 3 to change, 0 to destroy.
    ```

    The API reports that a secret is stored but not its value, so every resource with a stored secret plans an
    in-place update that sets it to the value of its variable, shown as `(sensitive value)`. If the plan proposes any
    other change, the configuration doesn't yet match your project, most often because a variable is missing.

3. Apply to adopt the project into your Terraform state:

    ```bash
    terraform apply
    ```

    Once applied, the `import` blocks have done their job. You can delete `import.tf`, or keep it, since Terraform
    ignores import blocks for resources that are already in state.

From here on, make changes by editing the `.tf` files and running `terraform apply`.

## Migrating from v0.3.x

In v0.3.x of the provider, the `descope_project` resource managed the whole project configuration. The new
provider only manages the project's name, environment and tags with it, and everything else with standalone
resources. After upgrading, use `tfexport` to generate those resources and adopt them into your existing state,
next to the `descope_project` resource you already have. The full procedure is in the provider's
[Upgrading from v0.3.x](https://registry.terraform.io/providers/descope/descope/latest/docs/guides/upgrading-from-v0.3)
guide.

```bash
tfexport -project P... -out ./generated -project-address descope_project.main
```

Copy the generated `.tf` files, and the `flows`, `widgets` and `templates` directories and `styles.json` when present,
into the directory that has your `descope_project` resource. Then follow [Applying the configuration](#applying-the-configuration).

If the project resource is in a module, pass its address as seen from inside that module, and the module's path with
`-import-prefix`. Terraform only allows import blocks in the root module, so move `import.tf` there:

```bash
tfexport -project P... -out ./generated -project-address descope_project.main -import-prefix module.auth
```

The variables in the generated `variables.tf` become inputs of that module, so pass their values in the `module` block,
and rename any that clash with inputs the module already has.

If your configuration manages several projects, for example with `for_each` or one `descope_project` resource per
environment, run the tool once per project with a different `-name-prefix`. The prefix is added to every generated
resource name, variable name and file name, so the exports can be copied into the same directory:

```bash
tfexport -project P... -out ./generated-dev -project-address 'descope_project.main["dev"]' -name-prefix dev
tfexport -project P... -out ./generated-prod -project-address 'descope_project.main["prod"]' -name-prefix prod
```
