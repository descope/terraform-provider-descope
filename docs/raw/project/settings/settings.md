
Settings
========



app_url
-------

- Type: `string` 

// description for app_url



custom_domain
-------------

- Type: `string` 

// description for custom_domain



approved_domains
----------------

- Type: `list` of `string` 

The list of approved domains that are allowed for redirect and verification URLs
for different authentication methods.



token_response_method
---------------------

- Type: `string` 
- Default: `"response_body"`

Configure how refresh tokens are managed by the Descope SDKs. Must be either `response_body`
or `cookies`. The default value is `response_body`.



cookie_policy
-------------

- Type: `string` 

Use "strict", "lax" or "none". To read more about custom domain and cookie policy
click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



cookie_domain
-------------

- Type: `string` 

The domain name for custom domain set up. To read more about custom domain and
cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



refresh_token_rotation
----------------------

- Type: `bool` 

Every time the user refreshes their session token via their refresh token, the
refresh token itself is also updated to a new one.



refresh_token_expiration
------------------------

- Type: `duration` 

The expiry time for the refresh token, after which the user must log in again. Use values
such as "4 weeks", "14 days", etc. The minimum value is "3 minutes".



session_token_expiration
------------------------

- Type: `duration` 

The expiry time of the session token, used for accessing the application's resources. The value
needs to be at least 3 minutes and can't be longer than the refresh token expiration.



step_up_token_expiration
------------------------

- Type: `duration` 

The expiry time for the step up token, after which it will not be valid and the user will
automatically go back to the session token.



trusted_device_token_expiration
-------------------------------

- Type: `duration` 

The expiry time for the trusted device token. The minimum value is "3 minutes".



access_key_session_token_expiration
-----------------------------------

- Type: `duration` 

The expiry time for access key session tokens. Use values such as "10 minutes", "4 hours", etc. The
value needs to be at least 3 minutes and can't be longer than 4 weeks.



enable_inactivity
-----------------

- Type: `bool` 

Use `True` to enable session inactivity. To read more about session inactivity
click [here](https://docs.descope.com/project-settings#session-inactivity).



inactivity_time
---------------

- Type: `duration` 

The session inactivity time. Use values such as "15 minutes", "1 hour", etc. The minimum
value is "10 minutes".



test_users_loginid_regexp
-------------------------

- Type: `string` 

Define a regular expression so that whenever a user is created with a matching login ID it will
automatically be marked as a test user.



test_users_verifier_regexp
--------------------------

- Type: `string` 

The pattern of the verifiers that will be used for testing.



test_users_static_otp
---------------------

- Type: `string` 

A 6 digit static OTP code for use with test users.



user_jwt_template
-----------------

- Type: `string` 

Name of the user JWT Template.



access_key_jwt_template
-----------------------

- Type: `string` 

Name of the access key JWT Template.
