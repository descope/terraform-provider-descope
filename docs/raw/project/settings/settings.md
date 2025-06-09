
Settings
========



app_url
-------

- Type: `string` 

The URL which your application resides on.



custom_domain
-------------

- Type: `string` 

A custom CNAME that's configured to point to `cname.descope.com`. Read more about custom
domains and cookie policy [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



approved_domains
----------------

- Type: `set` of `string` 

The list of approved domains that are allowed for redirect and verification URLs
for different authentication methods.



default_no_sso_apps
-------------------

- Type: `bool` 

Define whether a user created with no federated apps, will have access to all apps,
or will not have access to any app.



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



refresh_token_response_method
-----------------------------

- Type: `string` 
- Default: `"response_body"`

Configure how refresh tokens are managed by the Descope SDKs. Must be either `response_body`
or `cookies`. The default value is `response_body`.



refresh_token_cookie_policy
---------------------------

- Type: `string` 
- Default: `"none"`

Use `strict`, `lax` or `none`. Read more about custom domains and cookie policy
[here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



refresh_token_cookie_domain
---------------------------

- Type: `string` 

The domain name for refresh token cookies. To read more about custom domain and
cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



session_token_expiration
------------------------

- Type: `duration` 

The expiry time of the session token, used for accessing the application's resources. The value
needs to be at least 3 minutes and can't be longer than the refresh token expiration.



session_token_response_method
-----------------------------

- Type: `string` 
- Default: `"response_body"`

Configure how sessions tokens are managed by the Descope SDKs. Must be either `response_body`
or `cookies`. The default value is `response_body`.



session_token_cookie_policy
---------------------------

- Type: `string` 
- Default: `"none"`

Use `strict`, `lax` or `none`. Read more about custom domains and cookie policy
[here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



session_token_cookie_domain
---------------------------

- Type: `string` 

The domain name for session token cookies. To read more about custom domain and
cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



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



session_migration
-----------------

- Type: `object` of `settings.SessionMigration` 

Configure seamless migration of existing user sessions from another vendor to Descope.
