
OTP
====



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



domain
------

- Type: `string` 

The domain to embed in OTP messages.



expiration_time
---------------

- Type: `duration` 

The amount of time that an OTP code will be valid for.



email_service
-------------

- Type: `object` of `templates.EmailService` 

Settings related to sending emails with OTP codes.



text_service
------------

- Type: `object` of `templates.TextService` 

Settings related to sending SMS messages with OTP codes.



voice_service
-------------

- Type: `object` of `templates.VoiceService` 

Settings related to voice calls with OTP codes.
