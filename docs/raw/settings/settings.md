
Settings
========



cookie_policy
-------------

- Type: `string` 

// Optional. Use "strict", "lax" or "none". To read more about custom domain and cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



domain
------

- Type: `string` 

// Optional. The Domain name for custom domain set up. To read more about custom domain and cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).



enable_inactivity
-----------------

- Type: `bool`
- Default: `False`

// Optional. Specify True to enable session inactivity. To read more about session inactivity click [here](https://docs.descope.com/project-settings#session-inactivity).



inactivity_time
---------------

- Type: `duration` 

// Optional. Specify the inactivity timer, e.g "15 minutes", "1 hour". Minimum is "10 minutes".



refresh_token_expiration
------------------------

- Type: `duration` 

// Optional. Specify the refresh token expiration timer.  e.g "15 minutes", "1 hour". Minimum is "2 minutes".



user_jwt_template
-----------------

- Type: `string` 

// description for user_jwt_template



access_key_jwt_template
-----------------------

- Type: `string` 

// description for access_key_jwt_template
