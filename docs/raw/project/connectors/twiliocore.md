
TwilioCore
==========



account_sid
-----------

- Type: `string` (required)

// description for account_sid



senders
-------

- Type: `object` of `connectors.TwilioCoreSendersField` (required)

// description for senders



authentication
--------------

- Type: `object` of `connectors.TwilioAuthField` (required)

// description for authentication





TwilioCoreSendersField
======================



sms
----

- Type: `object` of `connectors.TwilioCoreSendersSMSField` (required)

// description for sms



voice
-----

- Type: `object` of `connectors.TwilioCoreSendersVoiceField` 

// description for voice





TwilioCoreSendersSMSField
=========================



phone_number
------------

- Type: `string` 

// description for phone_number



messaging_service_sid
---------------------

- Type: `string` 

// description for messaging_service_sid





TwilioCoreSendersVoiceField
===========================



phone_number
------------

- Type: `string` (required)

// description for phone_number





TwilioAuthField
===============



auth_token
----------

- Type: `secret` 

// description for auth_token



api_key
-------

- Type: `secret` 

// description for api_key



api_secret
----------

- Type: `secret` 

// description for api_secret
