
> **Community Fork** — This is an independently maintained fork of [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope).
> It is **not** an official Descope product. For the official provider, see the [upstream repository](https://github.com/descope/terraform-provider-descope).

<div align="center">
  <h3 align="center">Descope Terraform Provider (Community Fork)</h3>

  <p align="center">
    A community-maintained Terraform provider for managing Descope projects with extended resource coverage
  </p>
</div>

<br />

## Why This Fork?

The official Descope Terraform provider supports project-level configuration through the `descope_project` resource. However, the Descope Management API exposes many additional entities that are not yet available as Terraform resources.

This fork aims to close that gap by adding standalone resources for the full Management API surface:

| Resource | Upstream | This Fork |
|----------|:--------:|:---------:|
| Project settings, auth methods, connectors, flows | Yes | Yes |
| Roles & permissions (nested in project) | Yes | Yes |
| Applications (nested in project) | Yes | Yes |
| **Tenants** (standalone CRUD) | No | Planned |
| **Users** (standalone CRUD) | No | Planned |
| **Access Keys** (M2M) | No | Planned |
| **SSO Configuration** (per-tenant) | No | Planned |
| **SSO Applications** (standalone) | No | Planned |
| **Groups** (SCIM) | No | Planned |
| **Fine-Grained Authorization** (FGA) | No | Planned |
| **Third-Party Applications** | No | Planned |
| **Outbound Applications** | No | Planned |
| **Audit** (data source) | No | Planned |
| **Analytics** (data source) | No | Planned |
| **Password Settings** (standalone) | No | Planned |
| **Standalone Roles & Permissions** | No | Planned |
| **Standalone Flows** | No | Planned |
| **Project Export/Import** | No | Planned |

See the [open issues](https://github.com/jamescrowley321/terraform-provider-descope/issues) for the full roadmap.

<br/>

## About

Use this Terraform provider to manage your [Descope](https://www.descope.com) project using Terraform configuration files.

* Modify project settings and authentication methods.
* Create connectors, roles, permissions, applications and other entities.
* Use custom themes and flows created in the Descope console.
* Ensure dependencies between entities are satisfied.
* **Manage tenants, users, access keys, SSO, and more as standalone resources** (coming soon).

<br/>

## Getting Started

### Requirements

-   The [Terraform CLI](https://developer.hashicorp.com/terraform/install) installed.
-   A pro or enterprise tier license for your Descope company.
-   A valid management key for your Descope company. You can create one in the
    [Company section](https://app.descope.com/settings/company) of the Descope console.

### Usage

> **Note:** Until this fork is published to a Terraform registry, you must build and install the provider locally.
> See [Development](#development) below.

Configure the Descope provider with the management key and declare
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
    "http" = [
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

See the [README](internal/README.md) file in the `internal` directory for more details about the development
process, architecture, and tools.

### Setup

Clone this repository and run `make dev` to prepare your local environment for development. This will ensure
you have the requisite `go` compiler, build and install the Descope Terraform Provider binary to `$GOPATH/bin`,
and create a `~/.terraformrc` override file to instruct `terraform` to use the local provider binary instead
of loading it from the Terraform registry.

```bash
git clone https://github.com/jamescrowley321/terraform-provider-descope
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

## Contributing

Contributions are welcome, including AI-assisted submissions. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

<br/>

## Relationship to Upstream

This fork tracks the upstream [descope/terraform-provider-descope](https://github.com/descope/terraform-provider-descope) repository. Upstream changes will be merged periodically to stay current. New resources developed here may be proposed back to the upstream project via pull requests.

### Attribution

This project is based on the work of [Descope](https://www.descope.com) and the original provider authors. The original source code is licensed under the [MIT License](LICENSE).

<br/>

## Support

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Issues

If anything is missing or not working correctly please [open an issue](https://github.com/jamescrowley321/terraform-provider-descope/issues).

#### Descope Community

For general Descope questions (not specific to this fork) you can use the [Slack community](https://www.descope.com/community) or contact [Descope support](mailto:support@descope.com).
