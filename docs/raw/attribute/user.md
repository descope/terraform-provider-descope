
UserAttribute
=============



id
----

- Type: `string` (required)

The immutable machine name that identifies the attribute. This value is called `Machine Name` in the
Descope console. Changing this value will require the resource to be deleted and recreated.



project_id
----------

- Type: `string` (required)

The ID of the project that this user attribute belongs to. Changing this value will require the
resource to be deleted and recreated.



name
----

- Type: `string` (required)

The display name of the attribute. This value is called `Display Name` in the Descope console.



type
----

- Type: `string` (required)

The type of the attribute. Choose one of "string", "number", "boolean", "singleselect",
"multiselect", "date", "monthday". Changing this value will require the resource to be deleted and
recreated.



select_options
--------------

- Type: `set` of `string`

The list of options to choose from when the attribute type is "singleselect" or "multiselect".



widget_authorization
--------------------

- Type: `object` of `attribute.WidgetAuthorization`

Determines the permissions users are required to have to access this attribute in the user management
widget.
