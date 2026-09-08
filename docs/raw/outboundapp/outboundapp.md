
OutboundApp
===========



project_id
----------

- Type: `string` (required)

The ID of the Descope project this outbound application belongs to. Changing this value will
require the resource to be deleted and recreated.



tenant_id
---------

- Type: `string`

The ID of the tenant that owns this outbound application, for applications scoped to a single
tenant rather than the whole project. An empty value (the default) creates a project level
application. This value cannot be changed after creation, and changing it is rejected when the plan
is generated. Moving an application to a different tenant means creating a new one, because deleting
an outbound application also deletes every access token stored against it.



name
----

- Type: `string` (required)

A name for the outbound application.



description
-----------

- Type: `string`

A description for the outbound application.



logo
----

- Type: `string`

A URL for an image to display as the outbound application's logo.



app_type
--------

- Type: `string`
- Default: `"oauth"`

The kind of credential this application holds. Use `oauth` for an application that runs the OAuth
flow, so that Descope obtains and refreshes tokens on the user's behalf, or `apikey` for an
application whose tokens are uploaded directly and never refreshed.



client_id
---------

- Type: `string`

The OAuth client ID issued by the third party provider.



client_secret
-------------

- Type: `secret`

The OAuth client secret issued by the third party provider. The Descope API never returns this
value, so it cannot be read back and is not verified against the server. Omitting the field leaves
any previously stored secret in place, while setting it to a new value replaces it. An empty
string is not accepted.



discovery_url
-------------

- Type: `string`

The provider's OpenID Connect discovery endpoint, used to look up its other endpoints.



authorization_url
-----------------

- Type: `string`

The provider's authorization endpoint, which users are redirected to in order to grant consent.



authorization_url_params
------------------------

- Type: `list` of `outboundapp.URLParam`

Additional query parameters to append to the authorization request, for providers that require
them.



token_url
---------

- Type: `string`

The provider's token endpoint, used to exchange an authorization code and to refresh tokens.



token_url_params
----------------

- Type: `list` of `outboundapp.URLParam`

Additional query parameters to append to the token request, for providers that require them.



revocation_url
--------------

- Type: `string`

The provider's token revocation endpoint.



default_scopes
--------------

- Type: `set` of `string`

The scopes requested from the provider when no scopes are given for a specific connection.



default_redirect_url
--------------------

- Type: `string`

The URL that users are returned to once they have finished granting or denying consent at the
provider. This can be overridden per connection.



callback_domain
---------------

- Type: `string`

The domain that the provider redirects back to after consent. Leave this empty to use the project's
configured domain.



pkce
----

- Type: `bool`

Whether to use PKCE when exchanging the authorization code, which is required by some providers and
recommended by all of them.



access_type
-----------

- Type: `string`

Whether the provider should issue a refresh token alongside the access token. Use `offline` to
request one, so that Descope can keep the connection alive without the user present, or `online`
for an access token only.



prompt
------

- Type: `set` of `string`

How the provider should behave when the user is already signed in. Common values are `none`,
`login`, `consent` and `select_account`.
