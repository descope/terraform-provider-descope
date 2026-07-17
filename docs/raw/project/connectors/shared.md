
SenderField
===========



email
-----

- Type: `string` (required)

The email address that will appear as the sender of the email.



name
----

- Type: `string`

The display name that will appear as the sender of the email.





ServerField
===========



host
----

- Type: `string` (required)

The hostname or IP address of the SMTP server.



port
----

- Type: `int`
- Default: `25`

The port number to connect to on the SMTP server.





AuditFilterField
================



key
----

- Type: `string` (required)

The field name to filter on (either 'actions' or 'tenants').



operator
--------

- Type: `string` (required)

The filter operation to apply ('includes' or 'excludes').



values
------

- Type: `list` of `string` (required)

The list of values to match against for the filter.





HTTPAuthField
=============



bearer_token
------------

- Type: `secret`

Bearer token for HTTP authentication.



basic
-----

- Type: `object` of `connectors.HTTPAuthBasicField`

Basic authentication credentials (username and password).



api_key
-------

- Type: `object` of `connectors.HTTPAuthAPIKeyField`

API key authentication configuration.



oauth2_client_credentials
-------------------------

- Type: `object` of `connectors.HTTPAuthOAuth2ClientCredentialsField`

OAuth 2.0 client credentials configuration used to fetch an access token before making requests.





HTTPAuthBasicField
==================



username
--------

- Type: `string` (required)

Username for basic HTTP authentication.



password
--------

- Type: `secret` (required)

Password for basic HTTP authentication.





HTTPAuthAPIKeyField
===================



key
----

- Type: `string` (required)

The API key.



token
-----

- Type: `secret` (required)

The API secret.





HTTPAuthOAuth2ClientCredentialsField
====================================



client_id
---------

- Type: `string` (required)

The OAuth 2.0 client ID used to authenticate against the token endpoint.



client_secret
-------------

- Type: `secret` (required)

The OAuth 2.0 client secret used to authenticate against the token endpoint.



auth_url
--------

- Type: `string` (required)

The token endpoint URL used to request an access token.



auth_style
----------

- Type: `string`
- Default: `"header"`

How the client credentials are sent to the token endpoint. Either `header` to send them in the
`Authorization` header, or `params` to send them in the request body.



scopes
------

- Type: `string`

A space-separated list of OAuth scopes to request when fetching the access token.



token_request_headers
---------------------

- Type: `map` of `string`

Additional headers to include in the token request sent to the token endpoint.
