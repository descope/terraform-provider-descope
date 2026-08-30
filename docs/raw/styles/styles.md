
Styles
======



project_id
----------

- Type: `string` (required)

The ID of the project that the styles belong to. Changing this value will require the resource to be
deleted and recreated.



data
----

- Type: `string` (required)

The JSON data of the styles in their exported theme representation, defining the visual styling of
the project's flow pages. This will usually be exported as a `.json` file from the Descope console,
and set in the `.tf` file using the `data = file("...")` syntax.
