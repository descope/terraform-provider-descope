
Settings
========



cookie_policy
-------------

- Type: `string` 

Use "strict", "lax" or "none". To read more about custom domain and cookie policy
click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



domain
------

- Type: `string` 

The Domain name for custom domain set up. To read more about custom domain and
cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



approved_domains
----------------

- Type: `list` of `string` 

The list of approved domains that are allowed for redirect and verification URLs
for different authentication methods.



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



refresh_token_expiration
------------------------

- Type: `duration` 

The expiry time for the refresh token, after which the user must log in again. Use values
such as "4 weeks", "14 days", etc. The minimum value is "3 minutes".



access_key_session_token_expiration
-----------------------------------

- Type: `duration` 

The expiry time for access key session tokens. Use values such as "10 minutes", "4 hours", etc. The
value needs to be at least 3 minutes and can't be longer than 4 weeks.



user_jwt_template
-----------------

- Type: `string` 

Name of the user JWT Template.



access_key_jwt_template
-----------------------

- Type: `string` 

Name of the access key JWT Template.
