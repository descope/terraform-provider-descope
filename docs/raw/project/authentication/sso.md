
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

sso_suite_settings: Configuration block for the SSO Suite. 




SSOSuite
========



style_id
--------

- Type: `string` 


Specifies the style ID to apply in the SSO Suite. Ensure a style with this ID exists in the console for it to be used.



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
Note: At least one of SAML or OIDC must remain enabled.



hide_oidc
---------

- Type: `bool` 


Setting this to `true` will hide the OIDC configuration option.
Note: At least one of SAML or OIDC must remain enabled.
