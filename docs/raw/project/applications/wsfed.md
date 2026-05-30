
WSFed
=====



id
----

- Type: `string`

An optional identifier for the WS-Fed application.



name
----

- Type: `string` (required)

A name for the WS-Fed application.



description
-----------

- Type: `string`

A description for the WS-Fed application.



logo
----

- Type: `string`

A logo for the WS-Fed application. Should be a hosted image URL.



disabled
--------

- Type: `bool`

Whether the application should be enabled or disabled.



realm
-----

- Type: `string`

The WS-Fed realm identifier for the application.



reply_url
---------

- Type: `string`

The default reply URL where WS-Fed responses are sent. Used for IdP-initiated flows and when no `wreply` is supplied by the RP.



reply_allowed_callback_urls
---------------------------

- Type: `set` of `string`

Additional allowed `wreply` callback URLs beyond `reply_url`. Each entry may include the `*` wildcard. When the RP supplies a `wreply` parameter, it must match either the default `reply_url` or one of these patterns.



login_page_url
--------------

- Type: `string`

The Flow Hosting URL.



attribute_mapping
-----------------

- Type: `list` of `applications.AttributeMapping`

A list of attribute mappings from Descope user attributes to WS-Fed assertion attributes.



groups_mapping
--------------

- Type: `list` of `applications.GroupsMapping`

A list of group mappings from Descope roles to WS-Fed groups.



force_authentication
--------------------

- Type: `bool`

This configuration overrides the default behavior of the SSO application and forces the user to
authenticate via the Descope flow, regardless of the SP's request.



logout_redirect_url
-------------------

- Type: `string`

The URL to redirect to after logout.



error_redirect_url
------------------

- Type: `string`

The URL to redirect to when an error occurs.





GroupsMapping
=============



name
----

- Type: `string` (required)

The name of the groups mapping.



type
----

- Type: `string` (required)

The type of the groups mapping.



filter_type
-----------

- Type: `string` (required)

The filter type for the groups mapping.



value
-----

- Type: `string` (required)

The value of the groups mapping.



roles
-----

- Type: `list` of `applications.RoleGroupMapping`

The `RoleGroupMapping` object. A list of roles mapped to this group.





RoleGroupMapping
================



id
----

- Type: `string` (required)

The identifier of the role.



name
----

- Type: `string` (required)

The name of the role.
