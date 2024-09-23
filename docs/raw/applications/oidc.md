
OIDC
====



id
----

- Type: `string` 

An identifier for the OIDC application.



name
----

- Type: `string` (required)

A name for the OIDC application.



description
-----------

- Type: `string` 

A description for the OIDC application.



logo
----

- Type: `string` 

A logo for the OIDC application. Should be a hosted image URL.



disabled
--------

- Type: `bool` 

Whether the application should be enabled or disabled.



login_page_url
--------------

- Type: `string` 

The Flow Hosting URL. Read more about using this parameter with custom domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



claims
------

- Type: `list` of `string` 

A list of supported claims. e.g. `sub`, `email`, `exp`.
