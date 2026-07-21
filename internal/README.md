# Development

## Getting Started

This guide helps you set up a local development environment and work with existing Descope projects.

### Prerequisites

1. Run `make dev` to build and install the local provider (see [Setup and Installation](#setup-and-installation))
2. Obtain a management key from the Descope console (see [Configuration](#configuration) for environment setup)
3. Create a working directory with a `main.tf` file containing the below provider configuration.
4. Run `terraform init` to initialize the Terraform provider in the working directory.
5. Continue to the next steps depending on whether you want to create a new project or import an existing one.

```hcl
terraform {
  required_providers {
    descope = {
      source = "descope/descope"
    }
  }
}

provider "descope" {
  management_key = "K..."
  base_url       = "https://api.descope.com"  # or your local instance
}
```

### Importing an Existing Project

#### Add an import block:

In your `main.tf` file, append an `import` block for the plan generation command to be able to know which
project's configuration to read:

```hcl
import {
  to = descope_project.my_project
  id = "P..."  # your project ID
}
```

#### Generate configuration from your existing project:

Ignore any errors, as long as the `export.tf` file is created this is working as intended.

```bash
terraform plan -generate-config-out="export.tf"
```

#### Clean up the generated file:

- Remove any `oauth` system configurations that shouldn't be managed.
- Remove any inline `flows`, `styles`, email templates, etc.
- Replace any secret placeholders with the actual secret values.

#### Import state for the existing project:

Move the contents of the `export.tf` into the `main.tf` file, replacing the `import` block completely, then
run this command with the actual projectId to initialize the Terraform state for the existing project. This
is required as otherwise Terraform will try to create the project.

```bash
terraform import descope_project.my_project P...
```

#### Make changes and apply:

Edit your `.tf` file (e.g., change the project name) and run:

```bash
terraform plan   # preview changes
terraform apply  # apply changes
```

### Creating a New Project

#### Add a project resource to your `main.tf`:

```hcl
resource "descope_project" "my_project" {
  name = "My New Project"
}
```

#### Initialize and apply:

```bash
terraform plan   # preview changes
terraform apply  # create the project
```

## Commands

### Setup and Installation

- `make dev` - Prepares development environment (runs `make install` and `make terraformrc`)
- `make install` - Builds and installs terraform-provider-descope to $GOPATH/bin
- `make terraformrc` - Creates ~/.terraformrc to use local provider binary instead of registry

### Testing

- `make testacc` - Runs acceptance tests (requires environment configuration, see below)
- `make testcoverage` - Runs all tests with coverage analysis and generates coverage.html
- `make testcleanup` - Cleans up test projects using descope CLI (sometimes needed after failures)
- Use `tests=pattern` flag to run specific tests: `make testacc tests=TestProjectResource`

### Code Generation and Validation

- `make terragen` - Runs code generation for connectors, models, and documentation
- `make docs` - Generates Terraform registry documentation using tfplugindocs
- `make lint` - Runs golangci-lint and gitleaks security checks

## Configuration

Some makefile commands require these environment variables or a config file at `tools/config.env` with:

```bash
DESCOPE_MANAGEMENT_KEY=K...               # required for testacc
DESCOPE_BASE_URL=https://api.descope.com  # optional for testacc
DESCOPE_TEMPLATES_PATH=...                # required for terragen
```

## Sources

### Project Structure

The project source files are organized in this manner, though usually changes are only needed in the `models` layer:

- **Resources Layer** (`internal/resources/`): Terraform resource implementations (CRUD operations)
- **Entities Layer** (`internal/entities/`): Business logic layer that handles schema, validation, and API conversion
- **Models Layer** (`internal/models/`): Core data structures with Terraform Framework schema definitions
- **Infrastructure Layer** (`internal/infra/`): HTTP client and API communication

### Model Interfaces

Key interfaces reside in `internal/models/helpers/model.go` and are used by the model structs:

- `Model[T]`: Basic model with Values/SetValues methods for API serialization
- `NamedModel[T]`: Models with name/ID matching for friendly diffs
- `KeyedModel[T]`: Models with key matching for preserving IDs
- `CollectReferencesModel[T]`: Models that reference other models
- `UpdateReferencesModel[T]`: Models needing post-creation reference updates

### Model Implementation

Model implementations follow a consistent pattern. For example, for a model `Foo` we'll find:

- `FooAttributes`: Map of attributes that define the Terraform schema
- `FooModel`: A Go struct that's instantiated by Terraform according to the schema
- `Values`: A function on `FooModel` that returns a `map[string]any` representation of the model
- `SetValues`: A function on `FooModel` that updates its attributes with the server response

### Reference Resolution System

The provider tracks references between models:

1. `CollectReferences()` - Gathers existing model references
2. `Values()` - Converts model to API format using collected references  
3. `SetValues()` - Updates model from API response
4. `UpdateReferences()` - Resolves server IDs back to local references

### Connector System

Connectors are dynamically generated from templates in the Descope API schema. The generation process:

- Parses connector metadata from API templates
- Generates Go models with proper Terraform schema
- Creates test files and documentation
- Updates naming mappings in `naming.json`
- Custom naming by editing `naming.json` and rerunning `make terragen`

Note that some providers (like `smtp`, `sendgrid`, etc) have custom implementations.
