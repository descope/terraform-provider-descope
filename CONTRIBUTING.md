# Contributing

Thank you for your interest in contributing to this community fork of the Descope Terraform Provider.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (see `go.mod` for the required version)
- [Terraform CLI](https://developer.hashicorp.com/terraform/install)
- A Descope management key (required for acceptance tests)

### Setup

```bash
git clone https://github.com/jamescrowley321/terraform-provider-descope
cd terraform-provider-descope
make dev
```

This builds the provider binary and configures a local `~/.terraformrc` override so Terraform uses your local build.

### Build and Test

```bash
make install       # Rebuild and install the provider
make testacc       # Run acceptance tests
make testcoverage  # Run tests with coverage report
make lint          # Run linting and security checks
```

## Workflow

1. **Open an issue first** — Describe the feature, bug, or change you want to make. This avoids duplicate work and lets us align on the approach.
2. **Fork and branch** — Create a feature branch from `main`:
   ```bash
   git checkout -b feat/my-feature main
   ```
3. **Make your changes** — Follow the architecture patterns described in `internal/README.md`.
4. **Test** — Add or update unit tests and acceptance tests. All tests must pass.
5. **Lint** — Run `make lint` before committing.
6. **Commit** — Use [Conventional Commits](https://www.conventionalcommits.org/) (see below).
7. **Push and open a PR** — Push your branch and open a pull request against `main`.

### Conventional Commits

All commit messages must follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<optional scope>): <description>

<optional body>

<optional footer(s)>
```

**Types:** `feat`, `fix`, `docs`, `test`, `chore`, `refactor`, `ci`, `build`, `perf`, `style`

**Examples:**

```
feat(tenant): add descope_tenant resource with full CRUD
fix(access-key): handle missing expiration in API response
docs: update README with new resource table
test(user): add acceptance tests for user creation
```

**Breaking changes:** Add `!` after the type/scope and include a `BREAKING CHANGE:` footer:

```
feat(api)!: rename descope_project settings block

BREAKING CHANGE: The `project_settings` block has been renamed to `settings`.
```

### Branch Naming

Use the same type prefixes for branch names:

```
feat/tenant-resource
fix/access-key-expiration
docs/update-readme
```

## AI-Assisted Contributions

Contributions that use AI tools (GitHub Copilot, Claude, ChatGPT, Cursor, etc.) are welcome. We apply the same quality standards to all contributions regardless of how they were authored.

### Requirements for AI-Assisted PRs

- **All CI checks must pass** — lint, tests, security scans. No exceptions.
- **Tests are required** — AI-generated code must include unit tests and acceptance tests with adequate coverage.
- **Disclose tooling** — Note which AI tools were used in the PR description. A simple line like "Co-authored with Claude Code" is sufficient.
- **Human review is required** — All PRs require at least one human reviewer approval before merge.
- **You are responsible** — The submitter is accountable for the correctness, security, and quality of the code, regardless of whether it was AI-generated.

### What We Look For

- Code follows the existing architecture patterns (see `internal/README.md`)
- No hallucinated APIs or invented SDK methods — verify against the [Descope Go SDK](https://github.com/descope/go-sdk)
- Tests actually run and cover the new functionality
- Documentation is accurate and complete

## Adding New Resources

When adding a new Terraform resource (e.g., `descope_tenant`), follow this checklist:

1. [ ] Define the model in `internal/models/` following the existing patterns
2. [ ] Implement the entity in `internal/entities/`
3. [ ] Add the resource to the provider in `internal/resources/`
4. [ ] Add unit tests
5. [ ] Add acceptance tests
6. [ ] Generate documentation with `make docs`
7. [ ] Update the resource table in `README.md`

See `internal/README.md` for detailed architecture documentation.

## Code of Conduct

Be respectful, constructive, and collaborative. We're building something useful together.

## Questions?

Open a [discussion or issue](https://github.com/jamescrowley321/terraform-provider-descope/issues) and we'll help you get started.
