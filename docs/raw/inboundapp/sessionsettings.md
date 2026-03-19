
SessionSettings
===============



enabled
-------

- Type: `bool`

Whether to override the project's session settings.



refresh_token_expiration
------------------------

- Type: `duration`

The expiration duration for refresh tokens issued to this inbound app.



session_token_expiration
------------------------

- Type: `duration`

The expiration duration for session tokens issued to this inbound app.



key_session_token_expiration
----------------------------

- Type: `duration`

The expiration duration for access key session tokens. Must be between 3 minutes and one month.



user_template_id
----------------

- Type: `string`

The ID of the JWT template to use for user JWTs issued to this inbound app.



key_template_id
---------------

- Type: `string`

The ID of the JWT template to use for access key JWTs issued to this inbound app.
