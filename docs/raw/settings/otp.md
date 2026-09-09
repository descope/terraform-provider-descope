
OTPSettings
===========



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using OTP authentication directly via API and SDK calls. Note
that this does not affect authentication flows that are configured to use OTP.



domain
------

- Type: `string`

The domain to embed in OTP messages.



expiration_time
---------------

- Type: `duration`

The amount of time that an OTP code will be valid for.



email_connector_id
------------------

- Type: `string`

The ID of an email connector to use for sending emails, or `Descope` (the default) for the
built-in Descope delivery service.



text_connector_id
-----------------

- Type: `string`

The ID of an SMS connector to use for sending text messages, or `Descope` (the default) for the
built-in Descope delivery service.



voice_connector_id
------------------

- Type: `string`

The ID of a voice call connector to use for making voice calls, or `Descope` (the default) for
the built-in Descope delivery service.



email_template_id
-----------------

- Type: `string`

The ID of the email template to send to users, taken from a `descope_email_template` resource with
its `method` set to `otp`. An empty value (the default) selects the built-in System template.



text_template_id
----------------

- Type: `string`

The ID of the text template to send to users, taken from a `descope_text_template` resource with
its `method` set to `otp`. An empty value (the default) selects the built-in System template.



voice_template_id
-----------------

- Type: `string`

The ID of the voice template to use in calls to users, taken from a `descope_voice_template` resource
with its `method` set to `otp`. An empty value (the default) selects the built-in System template.
