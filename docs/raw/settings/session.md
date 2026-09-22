
SessionSettings
===============



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



user_jwt_template
-----------------

- Type: `string`

The ID of the JWT template used for user session tokens (e.g. `descope_jwt_template.example.id`).



access_key_jwt_template
-----------------------

- Type: `string`

The ID of the JWT template used for access key session tokens.



refresh_token_expiration
------------------------

- Type: `duration`

How long refresh tokens remain valid before they expire.



refresh_token_rotation
----------------------

- Type: `bool`

Whether to rotate the refresh token every time it's used to refresh the session.



refresh_token_response_method
-----------------------------

- Type: `string`
- Default: `"response_body"`

How refresh tokens are managed by the Descope SDKs, either `response_body` or `cookies`.



refresh_token_cookie_policy
---------------------------

- Type: `string`
- Default: `"none"`

The SameSite cookie policy for the refresh token cookie, either `strict`, `lax` or `none`.



refresh_token_cookie_domain
---------------------------

- Type: `string`

The domain the refresh token cookie is restricted to.



session_token_expiration
------------------------

- Type: `duration`

How long session tokens remain valid before they expire.



session_token_response_method
-----------------------------

- Type: `string`
- Default: `"response_body"`

How session tokens are managed by the Descope SDKs, either `response_body` or `cookies`.



session_token_cookie_policy
---------------------------

- Type: `string`
- Default: `"none"`

The SameSite cookie policy for the session token cookie, either `strict`, `lax` or `none`.



session_token_cookie_domain
---------------------------

- Type: `string`

The domain the session token cookie is restricted to.



step_up_token_expiration
------------------------

- Type: `duration`

How long step up tokens remain valid before they expire.



trusted_device_token_expiration
-------------------------------

- Type: `duration`

How long trusted device tokens remain valid before they expire.



access_key_session_token_expiration
-----------------------------------

- Type: `duration`

How long access key session tokens remain valid before they expire.



enable_inactivity
-----------------

- Type: `bool`

Whether to detect idle sessions and close them on behalf of the user.



inactivity_time
---------------

- Type: `duration`

How long a session can be inactive before it's closed, when `enable_inactivity` is set.
