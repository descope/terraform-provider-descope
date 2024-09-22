
Attributes
==========



tenant
------

- Type: `list` of `attributes.TenantAttribute` 

// Specify a list of `TenantAttribute`. Read the description below.



user
----

- Type: `list` of `attributes.UserAttribute` 

// Specify a list of `UserAttribute`. Read the description below.





TenantAttribute
===============



name
----

- Type: `string` (required)

// Specify the name of the attribute.



type
----

- Type: `string` (required)

// Specify the type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".



select_options
--------------

- Type: `list` of `string` 

// Specify when the attribute type is "multiselect". A list of options to chose from.





UserAttribute
=============



name
----

- Type: `string` (required)

// Specify the name of the attribute.



type
----

- Type: `string` (required)

// Specify the type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".



select_options
--------------

- Type: `list` of `string` 

// Specify when the attribute type is "multiselect". A list of options to chose from.



widget_authorization
--------------------

- Type: `object` of `attributes.UserAttributeWidgetAuthorization` 

// Specify the `UserAttributeWidgetAuthorization` object. Read the description below.





UserAttributeWidgetAuthorization
================================



view_permissions
----------------

- Type: `list` of `string` 

// Specify a list of permissions by name to set viewing permissions to the attribute in widgets. e.g "SSO Admin".



edit_permissions
----------------

- Type: `list` of `string` 

// Specify a list of permissions by name to set editing permissions to the attribute in widgets. e.g "SSO Admin".