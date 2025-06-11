
## Terragen

The provider uses code generation to maintain consistency and make it easier to update models
and documentation. The `terragen` tool has multiple phases:

1. **Connector Generation** (`conngen/`): Parses connector templates and generates Go models/tests
2. **Schema Parsing** (`schema/`): Extracts schema from existing Go models  
3. **Documentation Generation** (`docgen/`): Generates markdown docs from schemas
4. **Source Generation** (`srcgen/`): Creates compiled documentation Go files

Run the tool with `make terragen`. You can also use `make terragen flags='--skip-validate'` to prevent
it from aborting when missing or invalid data is encountered.
