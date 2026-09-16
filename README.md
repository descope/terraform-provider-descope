
<div align="center">
  <a href="https://github.com/descope/terraform-provider-descope">
    <img src=".github/images/descope-logo.png" alt="Descope Logo" width="160" height="160">
  </a>

  <h3 align="center">Descope Terraform Provider</h3>

  <p align="center">
    The official Terraform provider for managing Descope projects
  </p>
</div>

<br />

## About

Use the Descope Terraform Provider to manage your [Descope](https://www.descope.com) project
using Terraform configuration files.

* Modify project settings and authentication methods.
* Create connectors, roles, permissions, applications and other entities.
* Use custom themes and flows created in the Descope console.
* Reference entities from one another so they're created in the right order.

<br/>

## Getting Started

### Requirements

-   The [Terraform CLI](https://developer.hashicorp.com/terraform/install) installed.
-   A pro or enterprise tier license for your Descope company.
-   A valid management key for your Descope company. You can create one in the
    [Company section](https://app.descope.com/settings/company) of the Descope console.

### Usage

Declare the provider in your configuration and `terraform init` will automatically fetch and install the provider
for you from the [Terraform Registry](https://registry.terraform.io):

```hcl
terraform {
  required_providers {
    descope = {
      source = "descope/descope"
    }
  }
}
```

Configure the Descope provider with the management key as explained above and declare
a `descope_project` resource to create a new project for use with Terraform:

```hcl
provider "descope" {
  management_key = "K..."
}

resource "descope_project" "my_project" {
  name = "My Project"
}
```

Run `terraform plan` to ensure everything works, and then `terraform apply` if you want the project to actually
be created.

The `descope_project` resource manages the project itself and little else. Everything inside the project, from
authentication methods to roles, connectors and flows, is a separate resource that points back at it with a
`project_id` attribute. The examples below all assume the `my_project` resource declared above.

<br/>

## Examples

### Settings

Override the default values for specified project settings, in this case the session settings:

```hcl
resource "descope_session_settings" "my_settings" {
  project_id = descope_project.my_project.id

  refresh_token_expiration = "3 weeks"
  enable_inactivity = true
  inactivity_time = "1 hour"
}
```

The other settings resources work the same way, such as `descope_project_settings` for domains and security,
`descope_invite_settings` for user invitations, and `descope_otp_settings`, `descope_password_settings` and
the rest for the authentication methods.

### Authorization

Configure roles and permissions for users in the project. Roles refer to permissions by name, so use the `name`
attribute of the permission resources rather than hardcoding the names, and Terraform will know to create the
permissions first:

```hcl
resource "descope_permission" "build_apps" {
  project_id = descope_project.my_project.id
  name = "build-apps"
  description = "Allowed to build and sign applications"
}

resource "descope_permission" "upload_builds" {
  project_id = descope_project.my_project.id
  name = "upload-builds"
  description = "Allowed to upload new releases"
}

resource "descope_permission" "install_builds" {
  project_id = descope_project.my_project.id
  name = "install-builds"
  description = "Allowed to install beta releases"
}

resource "descope_role" "app_developer" {
  project_id = descope_project.my_project.id
  name = "App Developer"
  description = "Builds apps and uploads new beta builds"
  permissions = [
    descope_permission.build_apps.name,
    descope_permission.upload_builds.name,
    descope_permission.install_builds.name,
  ]
}

resource "descope_role" "app_tester" {
  project_id = descope_project.my_project.id
  name = "App Tester"
  description = "Installs and tests beta releases"
  permissions = [descope_permission.install_builds.name]
}
```

### Connectors and Flows

Setup a flow called `sign-up-or-in` by creating it in the Descope console in a development
project and exporting it as a `.json` file. Any entities the flow relies on need to be in the
plan as well, so in this example we also configure an HTTP connector with the expected name
`User Check` that the flow expects to be able to make use of. The names in the flow data are
matched against the entities in the project when the flow is imported, and nothing resolves them
again afterwards, so the connector has to exist by then and `depends_on` is what guarantees it.

```hcl
resource "descope_flow" "sign_up_or_in" {
  project_id = descope_project.my_project.id
  flow_id = "sign-up-or-in"
  data = file("flows/sign-up-or-in.json")

  depends_on = [descope_http_connector.user_check]
}

resource "descope_http_connector" "user_check" {
  project_id = descope_project.my_project.id
  name = "User Check"
  description = "A connector for checking if a new user is allowed to sign up"
  base_url = "https://example.com"

  authentication = {
    bearer_token = "<secret>"
  }
}
```

There's a resource for every connector type, named after the connector itself, so an SMTP connector is
a `descope_smtp_connector`, a Datadog connector is a `descope_datadog_connector`, and so on.

<br/>

## Development

See the [README](internal/README.md) file in the `internal` directory for more details about the development
process, architecture, and tools.

### Setup

Clone the repository and run `make dev` to prepare your local environment for development. This will ensure
you have the requisite `go` compiler, build and install the Descope Terraform Provider binary to `$GOPATH/bin`,
and create a `~/.terraformrc` override file to instruct `terraform` to use the local provider binary instead
of loading it from the Terraform registry.

```bash
git clone https://github.com/descope/terraform-provider-descope
cd terraform-provider-descope
make dev
```

### Build and Test

After making changes to source files, run `make install` to rebuild and install the provider. You can also run
the acceptance tests to ensure the provider works as expected.

```bash
# runs all unit and acceptance tests
make testacc

# or, to run all tests and compute code coverage
make testcoverage

# rebuild and install the provider
make install
```

<br/>

## Support

#### Contributing

If anything is missing or not working correctly please open an issue or pull request.

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Contact us

If you need help you can hop on our [Slack community](https://www.descope.com/community) or send an email to [Descope support](mailto:support@descope.com).
