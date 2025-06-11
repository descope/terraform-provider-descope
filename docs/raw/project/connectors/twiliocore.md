
TwilioCore
==========



account_sid
-----------

- Type: `string` (required)

Twilio Account SID from your Twilio Console.



senders
-------

- Type: `object` of `connectors.TwilioCoreSendersField` (required)

Configuration for SMS and voice message senders.



authentication
--------------

- Type: `object` of `connectors.TwilioAuthField` (required)

Twilio authentication credentials (either auth token or API key/secret).





TwilioCoreSendersField
======================



sms
----

- Type: `object` of `connectors.TwilioCoreSendersSMSField` (required)

SMS sender configuration using either a phone number or messaging service.



voice
-----

- Type: `object` of `connectors.TwilioCoreSendersVoiceField` 

Voice call sender configuration.





TwilioCoreSendersSMSField
=========================



phone_number
------------

- Type: `string` 

Twilio phone number for sending SMS messages.



messaging_service_sid
---------------------

- Type: `string` 

Twilio Messaging Service SID for sending SMS messages.





TwilioCoreSendersVoiceField
===========================



phone_number
------------

- Type: `string` (required)

Twilio phone number for making voice calls.





TwilioAuthField
===============



auth_token
----------

- Type: `secret` 

Twilio Auth Token for authentication.



api_key
-------

- Type: `secret` 

Twilio API Key for authentication (used with API Secret).



api_secret
----------

- Type: `secret` 

Twilio API Secret for authentication (used with API Key).
