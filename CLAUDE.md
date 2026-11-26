# CLAUDE.md

Guidance for Claude Code (claude.ai/code) when working in this repository.

## Code style

### Comments

Comments are held to an extremely high bar. Never add any comments to the codebase unless you can justify their
addition against to all the guidelines below:

- Remove any comment whose information the code already conveys. When in doubt remove the comment.
- Comments must never repeat what the code does.
- Never use comments to justify a change to a reviewer, and never refer to past code that has been changed.
- ASCII only. Soft cap of 150 characters per line for comments. The cap is deliberately soft, and slightly
  exceeding it to keep a comment on one line is better than wrapping. Wrap only when the comment genuinely
  needs more than one line.
- Only use multi-line comments next to code for a critical concept, background knowledge about backend behavior,
  or to explain a backwards compatibility reason.
- A function whose name says what it does gets no doc comment, or one line stating its single non-obvious behaviour.
- Test code should never need comments, other than a single line above each test case in an acceptance tests file,
  to assist readers to easily grasp what the test case checks.
- Structural markers such as `// Model` and `// Utils` organise code rather than document it. Keep them where a
  file already uses them, and note that terragen emits them into every generated connector.
- Comments in this repo should never divulge implementation details about private repositories and code, private
  issues, private conversations, etc.

To ensure the guidelines above are adhered to, don't add any comments to the source files directly when writing
code or while committing. Instead, any comments that you decide are actually needed should be written to a temporary
file (for example, `~/Desktop/comments-SLUG.md` with a slug of the branch or task name) with the filename and lines
and relevant code (like a patch file format). When work on a feature or branch is complete, we'll review which
comments merit putting into the code.

### Conventions

- Prefer verbosity over inconsistent abstractions.
- Models stay transport-free. Nothing under `internal/models/` imports `internal/infra`.
- Minimal diffs. No refactors that were not asked for.
- Use a soft cap of around 120 characters per line, but prefer dense code over pointless line wrapping.
- Keep function signatures and function calls on a single line unless wrapping is absolutely necessary.

### Testing

Write the tests before the code they test, so they prove the change works rather than just agreeing with it:

1. Write the tests that assert the desired behavior: acceptance tests, positive and negative unit tests, or both,
   whichever fits the change.
2. Run them and confirm they fail, and fail for the expected reason.
3. Make the change.
4. Run them again, along with the rest of the affected suite, and confirm everything passes.

When asserting errors with `testify`, check which error it is: use `require.ErrorIs` for a specific error, and
`require.ErrorContains` (alone or alongside it) for a specific message. `require.Error` and `assert.Error` only prove
that something failed, so there's almost never a reason to use them.

## Commands

- `make dev` - prepares the development environment (runs `make install` and `make terraformrc`)
- `make install` - builds and installs the provider to `$GOPATH/bin`
- `make terragen` - code generation for connectors, models and documentation
- `make docs` - registry documentation via tfplugindocs
- `make testacc` - acceptance and unit tests, against a real backend
- `make testcleanup` - removes leftover test projects
- `make sweep` - removes leftover `testacc-` entities from the shared test project
- `make lint` - golangci-lint and gitleaks

Run a subset with `tests=pattern`, e.g. `make testacc tests=TestPermission`. Do not use `$` anchors in that
pattern; make eats them.

## Environment

Acceptance tests and terragen read `tools/config.env`:

```bash
DESCOPE_MANAGEMENT_KEY=K...               # required for testacc
DESCOPE_BASE_URL=https://api.descope.com  # optional, point at a local stack to develop against it
DESCOPE_TESTACC_PROJECT_ID=P...           # existing project that project-scoped acceptance tests run in
DESCOPE_TEMPLATES_PATH=...                # required for terragen
```

Set `TF_UNSAFE_LOGS` to any non-empty value to log request and response payloads. They may contain secrets, which
is why it is deliberately separate from the `TF_LOG` variables.

## Generated code

Run `make terragen` after changing any model or connector template, then `make docs`. Never hand-edit generated
output: `docs/resources/*.md`, `internal/docs/docs.go`, the generated connector models, or
`internal/resources/connectors.go`. Attribute descriptions are authored in `docs/raw/` and compiled from there.

terragen has two requirements that produce confusing failures when unmet: `DESCOPE_TEMPLATES_PATH` must be set, and
the working directory's basename must be `terraform-provider-descope`, which bites inside a worktree. Use
`make terragen flags='--skip-validate'` to get past missing or invalid template data. Its phases are connector
generation (`conngen/`), schema parsing (`schema/`), documentation (`docgen/`) and source generation (`srcgen/`).

## Traps

**Local backend.** When `DESCOPE_BASE_URL` points at a local stack, the services run in docker and are rebuilt by
hand, so the running binary is often older than the backend source. Before concluding anything from probing a local
endpoint, compare the container's creation time (`docker ps --format '{{.Names}}\t{{.CreatedAt}}'`) against
`git log` for the file involved: a stale service looks exactly like a provider bug, and has cost a wrong diagnosis
and several wasted test runs.

**Removing an attribute from config does not always reset it.** `stringattr.Default("")` plans a change back to the
default when the attribute is deleted from a configuration. `stringattr.Optional()` is Optional+Computed with no
default, so it silently keeps the stored value instead, and the attribute cannot be cleared by deletion at all.
Choose deliberately; most attributes want `Default`.

**Acceptance tests hit a real backend** and create real projects. Clean up with `make testcleanup` after failures.

## Orientation

One Terraform resource per entity. `descope_project` carries only the project itself; every other resource
references it by `project_id`. Models under `internal/models/` define the schema and convert to and from the API
payload, `internal/resources/` wires them to CRUD operations, and `internal/infra/` is the HTTP client. Connector
models are generated from templates rather than written. The generic `/v1/mgmt/infra` entity endpoint is frozen:
new resources use dedicated endpoints.

Adding or changing a model attribute means reading `internal/models/CLAUDE.md` first, and adding a resource means
reading `internal/resources/CLAUDE.md`. The `tools/tfexport` tool has its own `CLAUDE.md`.
