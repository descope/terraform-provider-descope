
MagicLink
=========



enabled
-------

- Type: `bool` 

// description for enabled



expiration_time
---------------

- Type: `duration` 

// description for expiration_time



redirect_url
------------

- Type: `string` 

The URL to redirect users to after they log in using the magic link.



email_service
-------------

- Type: `object` of `templates.EmailService` 

Settings related to sending emails as part of the magic link authentication.



text_service
------------

- Type: `object` of `templates.TextService` 

Settings related to sending SMS messages as part of the magic link authentication.
