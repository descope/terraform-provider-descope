
OIDC
====



id
----

- Type: `string`

An optional identifier for the OIDC application.



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

The Flow Hosting URL. Read more about using this parameter with custom
domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



claims
------

- Type: `list` of `string`

A list of supported claims. e.g. `sub`, `email`, `exp`.



force_authentication
--------------------

- Type: `bool`

This configuration overrides the default behavior of the SSO application and forces
the user to authenticate via the Descope flow, regardless of the SP's request.



client_id
---------

- Type: `string`

A dedicated OIDC `client_id` to import for this application. Optional and **immutable** — changing it
forces the application to be recreated. When omitted, the `client_id` is computed by the server; when
set, it must be unique within the project. Mirrors the inbound third-party application `client_id`.



client_secret
-------------

- Type: `string` (sensitive)

A dedicated OIDC `client_secret` to import for this application, applied on creation only. When omitted,
a secret is generated server-side. The value is sensitive and is not returned on subsequent reads.



client_type
-----------

- Type: `string`

OAuth client confidentiality. One of `""` (default — legacy access-key authentication),
`"confidential"` (a dedicated client secret is generated for the app), or `"public"`.



approved_redirect_urls
----------------------

- Type: `list` of `string`

A list of approved redirect URLs for this application (supports `*` wildcards). When set,
redirect URIs are validated against this per-app list; when empty, validation falls back to
the project's approved/trusted domains.



authorization_code_disabled
---------------------------

- Type: `bool`

Disables the `authorization_code` grant type for this application.



client_credentials_disabled
---------------------------

- Type: `bool`

Disables the `client_credentials` grant type for this application.



refresh_token_disabled
----------------------

- Type: `bool`

Disables the `refresh_token` grant type for this application.



jwt_bearer_disabled
-------------------

- Type: `bool`

Disables the `urn:ietf:params:oauth:grant-type:jwt-bearer` grant type for this application.



device_code_disabled
--------------------

- Type: `bool`

Disables the `urn:ietf:params:oauth:grant-type:device_code` grant type for this application.



force_pkce
----------

- Type: `bool`

When enabled, the authorization code flow requires PKCE in addition to the normal client authentication. A confidential client must then present both its client secret and a valid PKCE `code_verifier`. Public clients always use PKCE regardless of this setting.



default_audience
----------------

- Type: `string`

Controls the default `aud` claim of tokens issued for this application. One of `"projectId"` (the project ID only), `"clientId"` (the dedicated client ID only), or `""` (default — both). Only applies to modern apps that set a `client_type`; legacy apps always use the project ID, so the empty default leaves their behavior unchanged.
