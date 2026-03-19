
InboundApp
==========



project_id
----------

- Type: `string` (required)

The ID of the Descope project this inbound app belongs to. Changing this value will require the
resource to be deleted and recreated.



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



default_audience
----------------

- Type: `string`

The default `aud` claim to include in tokens issued for this app. Use `projectId` to set the project ID
as the audience, `clientId` to set the app's client ID, or leave empty to include both.



non_confidential_client
-----------------------

- Type: `bool`

Whether this is a public (non-confidential) client that does not use a client secret. Changing this
value after creation will require the resource to be replaced.



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
