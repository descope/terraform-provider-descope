
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
* Create connectors, role, permissions, applications and other entities.
* Use custom themes and flows created in the Descope console.
* Ensure dependencies between entities are satisfied.

<br/>

## Getting Started

### Requirements

-   The [Terraform CLI](https://developer.hashicorp.com/terraform/install) installed.
-   A pro or enterprise tier license for your Descope company.
-   A valid management key for your Descope company. You can create one in the
    [Company section](https://app.descope.com/settings/company) of the Descope console.
-   The project ID of one of the existing projects in your Descope company. You can
    find it in the [Project section](https://app.descope.com/settings/project) of the
    Descope console. This project will not be modified in any way by the provider.

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

Configure the Descope provider with the project ID and management key above and declare a `descope_project` resource
to create a new project for use with Terraform:

```hcl
provider "descope" {
  project_id = "P..."
  management_key = "K..."
}

resource "descope_project" "my_project" {
  name = "My Project"
}
```

Run `terraform plan` to ensure everything works, and then `terraform apply` if you want the project to actually
be created.

<br/>

## Examples

### Project Settings

Override the default values for specified project settings:

```hcl
resource "descope_project" "my_project" {
  name = "My Project"

  project_settings = {
    refresh_token_expiration = "3 weeks"
    enable_inactivity = true
    inactivity_time = "1 hour"
  }
}
```

### Authorization

Configure roles and permissions for users in the project:

```hcl
resource "descope_project" "my_project" {
  name = "My Project"

  authorization = {
    roles = [
      {
        name = "App Developer"
        description = "Builds apps and uploads new beta builds"
        permissions = ["build-apps", "upload-builds", "install-builds"]
      },
      {
        name = "App Tester"
        description = "Installs and tests beta releases"
        permissions = ["install-builds"]
      },
    ]
    permissions = [
      {
        name = "build-apps"
        description = "Allowed to build and sign applications"
      },
      {
        name = "upload-builds"
        description = "Allowed to upload new releases"
      },
      {
        name = "install-builds"
        description = "Allowed to install beta releases"
      },
    ]
  }
}
```

### Connectors and Flows

Setup a flow called `sign-up-or-in` by creating it in the Descope console in a development
project and exporting it as a `.json` file. The provider will ensure that any entities used
by the flow such as connectors will be provided by the plan. In this example, we also configure
an HTTP connector with the expected name `User Check` that the flow expects to be able to
make use of.

```hcl
resource "descope_project" "my_project" {
  name = "My Project"

  flows = {
    "sign-up-or-in" = {
      data = file("flows/sign-up-or-in.json")
    }
  }

  connectors = {
    "http": [
      {
        name = "User Check"
        description = "A connector for checking if a new user is allowed to sign up"
        base_url = "https://example.com"
        bearer_token = "<secret>"
      }
    ]
  }
}
```

<br/>

## Development

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

After making changes to source files, run `make install` to rebuild the and install the provider. You can also run
the acceptance tests to ensure the provider works as expected.

```bash
# runs all unit and acceptance tests
make testacc

# or, to run all tests and compute code coverage
make testcoverage

# rebuild and install the provider
make install
```

### Terragen

The provider uses code generation where possible to ensure models, schemas, tests and docs are up to date and
synchronized. After making changes to the provider source code be sure to run the `terragen` tool to ensure all
schemas and docs are updated, and to be alerted if anything is missing.

```bash
make terragen
```

<br/>

## Support

#### Contributing

If anything is missing or not working correctly please open an issue or pull request.

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Contact us

If you need help you can hop on our [Slack community](https://www.descope.com/community) or send an email to [Descope support](mailto:support@descope.com).
