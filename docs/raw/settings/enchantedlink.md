
EnchantedLinkSettings
=====================



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using enchanted link authentication directly via API and SDK
calls. Note that this does not affect authentication flows that are configured to use enchanted
links.



expiration_time
---------------

- Type: `duration`

How long the enchanted link remains valid before it expires.



redirect_url
------------

- Type: `string`

The URL to redirect users to after they log in using the enchanted link.



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



email_template_id
-----------------

- Type: `string`

The ID of the email template to send to users, taken from a `descope_email_template` resource with
its `method` set to `enchantedlink`. An empty value (the default) selects the built-in System
template.



text_template_id
----------------

- Type: `string`

The ID of the text template to send to users, taken from a `descope_text_template` resource with
its `method` set to `enchantedlink`. An empty value (the default) selects the built-in System template.
