
OIDC
====



id
----

- Type: `string` 

// Specify an identifier for the OIDC application.



name
----

- Type: `string` (required)

// Specify a name for the OIDC application.



description
-----------

- Type: `string` 

// Specify a description for the OIDC application.



logo
----

- Type: `string` 

// Specify a logo for the OIDC application. Should be a hosted image URL.



disabled
--------

- Type: `bool` 

// Specify whether the application should be enabled or disabled.



login_page_url
--------------

- Type: `string` 

// Specify the Flow Hosting URL. Read more about using this parameter with custom domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



claims
------

- Type: `list` of `string` 

// Specify a list of supported claims. e.g. `sub`, `email`, `exp`.
