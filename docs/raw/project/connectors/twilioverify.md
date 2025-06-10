
TwilioVerify
============



account_sid
-----------

- Type: `string` (required)

Twilio Account SID from your Twilio Console.



service_sid
-----------

- Type: `string` (required)

Twilio Verify Service SID for verification services.



sender
------

- Type: `string` 

Optional sender identifier for verification messages.



authentication
--------------

- Type: `object` of `connectors.TwilioAuthField` (required)

Twilio authentication credentials (either auth token or API key/secret).
