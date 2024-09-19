
Attributes
==========



tenant
------

- Type: `list` of `attributes.TenantAttribute` 

// description for tenant



user
----

- Type: `list` of `attributes.UserAttribute` 

// description for user





TenantAttribute
===============



name
----

- Type: `string` (required)

// description for name



type
----

- Type: `string` (required)

// description for type



select_options
--------------

- Type: `list` of `string` 

// description for select_options





UserAttribute
=============



name
----

- Type: `string` (required)

// description for name



type
----

- Type: `string` (required)

// description for type



select_options
--------------

- Type: `list` of `string` 

// description for select_options



widget_authorization
--------------------

- Type: `object` of `attributes.UserAttributeWidgetAuthorization` 

// description for widget_authorization





UserAttributeWidgetAuthorization
================================



view_permissions
----------------

- Type: `list` of `string` 

// description for view_permissions



edit_permissions
----------------

- Type: `list` of `string` 

// description for edit_permissions
