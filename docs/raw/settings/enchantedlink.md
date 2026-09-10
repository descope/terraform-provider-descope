
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



email_service
-------------

- Type: `object` of `settings.EmailServiceRef`

Settings related to sending emails as part of the enchanted link authentication.



text_service
------------

- Type: `object` of `settings.TextServiceRef`

Settings related to sending SMS messages as part of the enchanted link authentication. Omitting this
block selects the built-in Descope delivery service, which also clears any custom text template that
had been selected, so a project configured through the console will be reset by an apply that does
not declare it.



email_template_id
-----------------

- Type: `string`

The ID of the email template to send to users, taken from a `descope_email_template` resource with
its `method` set to `enchantedlink`. An empty value (the default) selects the built-in System
template.



text_template_id
----------------

- Type: `string`

The ID of the text template to send to users, taken from a `descope_text_template` resource with its
`method` set to `enchantedlink`. An empty value (the default) selects the built-in System template.
