
EmailTemplate
=============



active
------

- Type: `bool` 

Whether this email template is currently active and in use.



name
----

- Type: `string` (required)

Unique name for this email template.



subject
-------

- Type: `string` (required)

Subject line of the email message.



html_body
---------

- Type: `string` 

HTML content of the email message body, required if `use_plain_text_body` isn't set.



plain_text_body
---------------

- Type: `string` 

Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.



use_plain_text_body
-------------------

- Type: `bool` 

Whether to use the plain text body instead of HTML for the email.
