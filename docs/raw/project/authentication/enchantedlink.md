
EnchantedLink
=============



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



expiration_time
---------------

- Type: `duration`

How long the enchanted link remains valid before it expires.



redirect_url
------------

- Type: `string`

The URL to redirect users to after they log in using the enchanted link.



allow_unverified_recipients
---------------------------

- Type: `bool`

By default, enchanted links are only sent to verified email addresses. Enabling this allows
sending them to unverified email addresses as well, which may increase the risk of spam and fraud.



email_service
-------------

- Type: `object` of `templates.EmailService`

Settings related to sending emails as part of the enchanted link authentication.
