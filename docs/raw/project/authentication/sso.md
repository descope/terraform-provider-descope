
SSO
====



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



merge_users
-----------

- Type: `bool` 

Whether to merge existing user accounts with new ones created through SSO authentication.



redirect_url
------------

- Type: `string` 

The URL the end user is redirected to after a successful authentication. If one is specified
in tenant level settings or SDK/API call, they will override this value.



sso_suite_settings
------------------

- Type: `object` of `authentication.SSOSuite` 

Configuration block for the SSO Suite.



allow_duplicate_domains
-----------------------

- Type: `bool` 

Whether to allow duplicate SSO domains across tenants.



allow_override_roles
--------------------

- Type: `bool` 

Whether to allow overriding user's roles with SSO related roles.



groups_priority
---------------

- Type: `bool` 

Whether to enable groups priority.



mandatory_user_attributes
-------------------------

- Type: `list` of `authentication.MandatoryUserAttribute` 

Define the required Descope attributes that must be populated when receiving SSO information.



limit_mapping_to_mandatory_attributes
-------------------------------------

- Type: `bool` 

Mapping to attributes not specified in `mandatory_user_attributes` is not allowed.



sso_domains_required
--------------------

- Type: `bool` 

When configuring SSO an SSO domain must be specified.



groups_attribute_name_required
------------------------------

- Type: `bool` 

When configuring SSO the groups attribute name must be specified.





MandatoryUserAttribute
======================



id
----

- Type: `string` (required)

The identifier for the attribute. This value is called `Machine Name` in the Descope console.



custom
------

- Type: `bool` 

Whether the attribute is a custom attribute defined in addition to the default Descope user attributes.





SSOSuite
========



style_id
--------

- Type: `string` 

Specifies the style ID to apply in the SSO Suite. Ensure a style with this ID exists in
the console for it to be used.



hide_scim
---------

- Type: `bool` 

Setting this to `true` will hide the SCIM configuration in the SSO Suite interface.



hide_groups_mapping
-------------------

- Type: `bool` 

Setting this to `true` will hide the groups mapping configuration section in the SSO Suite interface.



hide_domains
------------

- Type: `bool` 

Setting this to `true` will hide the domains configuration section in the SSO Suite interface.



hide_saml
---------

- Type: `bool` 

Setting this to `true` will hide the SAML configuration option.



hide_oidc
---------

- Type: `bool` 

Setting this to `true` will hide the OIDC configuration option.



force_domain_verification
-------------------------

- Type: `bool` 

Setting this to `true` will allow only verified domains to be used.
