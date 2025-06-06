
Attributes
==========



tenant
------

- Type: `list` of `attributes.TenantAttribute` 

A list of `TenantAttribute`. Read the description below.



user
----

- Type: `list` of `attributes.UserAttribute` 

A list of `UserAttribute`. Read the description below.





TenantAttribute
===============



name
----

- Type: `string` (required)

The name of the attribute.



type
----

- Type: `string` (required)

The type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".



select_options
--------------

- Type: `set` of `string` 

When the attribute type is "multiselect". A list of options to chose from.



authorization
-------------

- Type: `object` of `attributes.TenantAttributeAuthorization` 

// description for authorization





TenantAttributeAuthorization
============================



view_permissions
----------------

- Type: `set` of `string` 

// description for view_permissions





UserAttribute
=============



name
----

- Type: `string` (required)

The name of the attribute.



type
----

- Type: `string` (required)

The type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".



select_options
--------------

- Type: `set` of `string` 

When the attribute type is "multiselect". A list of options to chose from.



widget_authorization
--------------------

- Type: `object` of `attributes.UserAttributeWidgetAuthorization` 

The `UserAttributeWidgetAuthorization` object. Read the description below.





UserAttributeWidgetAuthorization
================================



view_permissions
----------------

- Type: `set` of `string` 

A list of permissions by name to set viewing permissions to the attribute in widgets. e.g "SSO Admin".



edit_permissions
----------------

- Type: `set` of `string` 

A list of permissions by name to set editing permissions to the attribute in widgets. e.g "SSO Admin".
