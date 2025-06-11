# Development

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
DESCOPE_PROJECT_ID=P...                   # required for testacc
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
- `MatchableModel[T]`: Models with name/ID matching for friendly diffs
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
