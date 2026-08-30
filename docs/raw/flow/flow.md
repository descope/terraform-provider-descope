
Flow
====



project_id
----------

- Type: `string` (required)

The ID of the project that the flow belongs to. Changing this value will require the resource to be
deleted and recreated.



flow_id
-------

- Type: `string` (required)

The machine identifier of the flow (e.g., `sign-up-or-in`), used to run the flow and to reference it
in other resources. Changing this value will require the resource to be deleted and recreated.



data
----

- Type: `string` (required)

The JSON data of the flow in its exported representation, including its metadata, contents, and
references. This will usually be exported as a `.json` file from the Descope console, and set in the
`.tf` file using the `data = file("...")` syntax.
