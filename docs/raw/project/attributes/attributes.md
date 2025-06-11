
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

When the attribute type is "multiselect". A list of options to choose from.



authorization
-------------

- Type: `object` of `attributes.TenantAttributeAuthorization` 

Determines the required permissions for this tenant.





TenantAttributeAuthorization
============================



view_permissions
----------------

- Type: `set` of `string` 

Determines the required permissions for this tenant.





UserAttribute
=============



name
----

- Type: `string` (required)

The name of the attribute.



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

- Type: `object` of `attributes.UserAttributeWidgetAuthorization` 

Determines the permissions users are required to have to access this attribute
in the user management widget.





UserAttributeWidgetAuthorization
================================



view_permissions
----------------

- Type: `set` of `string` 

The permissions users are required to have to view this attribute in the user management widget.



edit_permissions
----------------

- Type: `set` of `string` 

The permissions users are required to have to edit this attribute in the user management widget.
