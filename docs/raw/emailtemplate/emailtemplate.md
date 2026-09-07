
EmailTemplate
=============



project_id
----------

- Type: `string` (required)

The ID of the project that the email template belongs to. Changing this value will require the
resource to be deleted and recreated.



method
------

- Type: `string` (required)

The authentication method the email template is used with, e.g. `magiclink` or `otp`. Changing this value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the email template that's unique among the templates of the same authentication method.



subject
-------

- Type: `string` (required)

The subject line of emails sent with this template.



html_body
---------

- Type: `string`

The HTML body of emails sent with this template. Required unless `use_plain_text_body` is enabled.
Template macros such as `{{.code}}` can be used to insert dynamic values.



plain_text_body
---------------

- Type: `string`

The plain text body of emails sent with this template. Required when `use_plain_text_body` is
enabled.



use_plain_text_body
-------------------

- Type: `bool`

Whether emails are sent with the plain text body instead of the HTML body.
