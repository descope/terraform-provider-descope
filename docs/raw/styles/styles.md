
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
which exports one style per file, and set in the `.tf` file using the `data = file("...")` syntax.
Applying this resource creates and overwrites the styles the data carries, and leaves every other
style in the project untouched, so styles created in the console are safe from a Terraform apply.
For the same reason removing a style from the data does not remove it from the project - delete it
in the console or through the management API instead. An apply is rejected if the data's
`componentsVersion` differs from the project's, in which case export the styles again or leave
`componentsVersion` out, and if it changes the type of a style that already exists in the project.
