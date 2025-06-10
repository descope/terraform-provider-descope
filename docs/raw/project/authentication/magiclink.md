
MagicLink
=========



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



expiration_time
---------------

- Type: `duration` 

How long the magic link remains valid before it expires.



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
