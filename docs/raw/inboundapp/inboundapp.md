
InboundApp
==========



project_id
----------

- Type: `string` (required)

The ID of the Descope project this inbound app belongs to. Changing this value will require the
resource to be deleted and recreated.



deletion_protection
-------------------

- Type: `bool`

Protects the inbound app from being accidentally destroyed or replaced. Destroying or
replacing an inbound app issues a new client ID and client secret, breaking any
integrations that rely on the existing ones, so inbound apps are protected by default.
To allow the resource to be destroyed or replaced, for example in ephemeral test
environments, set this attribute to `false` and apply the change first.



name
----

- Type: `string` (required)

A name for the inbound app.



description
-----------

- Type: `string`

A description for the inbound app.



logo_url
--------

- Type: `string`

A URL to the inbound app's logo image.



login_page_url
--------------

- Type: `string`

The Flow Hosting URL.



approved_callback_urls
----------------------

- Type: `set` of `string`

A set of approved redirect URIs that the inbound app is allowed to redirect to after authorization.



permissions_scopes
------------------

- Type: `list` of `inboundapp.ApplicationScope`

A list of permission scopes that the inbound app can request. Permission scopes provide the app with
the ability to act on behalf of a user based on their roles and permissions.



attributes_scopes
-----------------

- Type: `list` of `inboundapp.ApplicationScope`

A list of user information scopes that the inbound app can request. Attribute scopes provide the app
with access to user profile data such as email, phone, or custom attributes.



connections_scopes
------------------

- Type: `list` of `inboundapp.ApplicationScope`

A list of connection scopes that the inbound app can request. Connection scopes provide the app with
the ability to access external tokens based on the mapped scopes.



session_settings
----------------

- Type: `object` of `inboundapp.SessionSettings`

Custom session management settings for this inbound app, overriding the project defaults.



audience_whitelist
------------------

- Type: `set` of `string`

A set of allowed custom `aud` claim values that the inbound app can request via the `resource`
parameter, per RFC 8707.



force_add_all_authorization_info
--------------------------------

- Type: `bool`

When enabled, all of the user's tenants, roles, and permissions will always be included in issued tokens.



force_dpop
----------

- Type: `bool`

Require clients to use DPoP (Demonstrating Proof of Possession), binding access tokens to a key held by the client so a stolen token cannot be used by anyone else.



default_audience
----------------

- Type: `string`

The default `aud` claim to include in tokens issued for this app. Use `projectId` to set the project ID
as the audience, `clientId` to set the app's client ID, or leave empty to include both.



client_type
-----------

- Type: `string`

The OAuth client type: `confidential` for a client that authenticates with a secret, or `public` for
one that authenticates with PKCE and has no secret. Setting it explicitly is recommended. The empty
value `""` is the legacy mode. If this attribute is omitted when the app is created, the account's
default client type is used and the stored value shows which one it is, and removing the attribute
from the configuration later keeps the current type. Moving from the legacy mode to `confidential`
updates the app in place and keeps its client ID, and any other change, including a move to the legacy
mode, replaces the resource, which issues a new client ID and secret. An app that was never issued a
client secret, such as one created with `non_confidential_client = true` in v0.3, needs one after
moving to `confidential`: rotate its client secret, for example in the Descope console, or force a
replacement with `terraform apply -replace` to get a new client ID and secret.



client_id
---------

- Type: `string`

A custom client ID for the inbound app. If not set, an ID will be generated automatically. Changing
this value after creation will require the resource to be replaced.



client_secret
-------------

- Type: `secret`

The client secret for authenticating this inbound app. This value is generated automatically and
cannot be retrieved after the resource is created. Store this value securely.



force_pkce
----------

- Type: `bool`

When enabled, the authorization code flow requires PKCE in addition to the normal client authentication. A confidential client must then present both its client secret and a valid PKCE `code_verifier`. Public clients always use PKCE regardless of this setting.



allowed_tenants
---------------

- Type: `set` of `string`

Restricts the app to these tenant IDs. Leave empty to allow all tenants.



scope_claim_mapping
-------------------

- Type: `list` of `inboundapp.ScopeClaimMapping`

Maps the claims that each scope contributes to the tokens issued for this app. This is the structured
form of `attributes_scopes` and takes precedence over it: when a requested scope appears in both, the
mapping here wins.
