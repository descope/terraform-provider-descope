
Widget
======



project_id
----------

- Type: `string` (required)

The ID of the project that the widget belongs to. Changing this value will require the resource to
be deleted and recreated.



widget_id
---------

- Type: `string` (required)

The machine identifier of the widget (e.g., `user-management`), used to embed the widget and to
reference it in other resources. Changing this value will require the resource to be deleted and
recreated.



data
----

- Type: `string` (required)

The JSON data of the widget in its exported representation, including its metadata and screens. This
will usually be exported as a `.json` file from the Descope console, and set in the `.tf` file using
the `data = file("...")` syntax.
