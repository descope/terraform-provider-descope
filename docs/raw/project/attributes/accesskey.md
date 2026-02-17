
AccessKeyAttribute
==================



id
----

- Type: `string` 

An optional identifier for the attribute. This value is called `Machine Name` in the Descope console.
If a value is not provided then an appropriate one will be created from the value of `name`.



name
----

- Type: `string` (required)

The name of the attribute. This value is called `Display Name` in the Descope console.



type
----

- Type: `string` (required)

The type of the attribute. Choose one of "string", "number", "boolean",
"singleselect", "multiselect", "date".



select_options
--------------

- Type: `set` of `string` 

When the attribute type is "multiselect". A list of options to choose from.



widget_authorization
--------------------

- Type: `object` of `attributes.AccessKeyAttributeWidgetAuthorization` 

Determines the permissions access key are required to have to access this attribute
in the access key management widget.





AccessKeyAttributeWidgetAuthorization
=====================================



view_permissions
----------------

- Type: `set` of `string` 

The permissions users are required to have to view this attribute in the access key management widget.



edit_permissions
----------------

- Type: `set` of `string` 

The permissions users are required to have to edit this attribute in the access key management widget.
