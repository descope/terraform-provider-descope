
FGASchema
=========



project_id
----------

- Type: `string` (required)

The ID of the project that this FGA schema belongs to. Changing this value will require the
resource to be deleted and recreated.



schema
------

- Type: `string`

The project's FGA (fine-grained authorization) schema, configured in the
[Descope console](https://app.descope.com/authorization/fga) under the FGA tab. Use the code view
to get the schema text and paste it as the value for this attribute. When set, the schema must
start with `model AuthZ`. Setting an empty schema, or destroying the resource, clears the
project's FGA schema, which also removes any FGA relations that depend on it.
