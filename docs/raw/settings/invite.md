
InviteSettings
==============



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



require_invitation
------------------

- Type: `bool`

Whether users must be invited before they can sign up to the project.



invite_url
----------

- Type: `string`

Custom URL to include in the message sent to invited users.



add_magiclink_token
-------------------

- Type: `bool`

Whether to include a magic link token in invitation messages.



expire_invited_users
--------------------

- Type: `bool`

Expire the user account if the invitation is not accepted within the expiration time.



invite_expiration
-----------------

- Type: `duration`

The expiry time for the invitation, meant to be used together with
`expire_invited_users` and/or `add_magiclink_token`. Use values such
as "2 weeks", "4 days", etc. The minimum value is "1 hour".



send_email
----------

- Type: `bool`
- Default: `true`

Whether to send invitation emails to users.



send_text
---------

- Type: `bool`

Whether to send invitation SMS messages to users.



email_connector_id
------------------

- Type: `string`

The ID of an email connector to use for sending invitation emails. An empty value (the default) selects the built-in Descope delivery service.



email_template_id
-----------------

- Type: `string`

The ID of the email template to send to invited users, taken from a `descope_email_template` resource
with its `method` set to `invite`. An empty value (the default) selects the built-in System template.



send_reminders
--------------

- Type: `bool`

Whether to automatically resend the invitation to users who have not accepted it yet, every
`reminder_interval` and up to `max_reminders` times. Reminders are only sent through a custom
connector, so turning this on requires `email_connector_id` to be set when invitations are sent
by email. Turning reminders on only affects invitations sent from then on.



reminder_interval
-----------------

- Type: `duration`

How long to wait before sending a reminder, measured from the last invitation or reminder sent
to the user. Use values such as "2 days", "1 week", etc. The minimum value is "1 day".



max_reminders
-------------

- Type: `int`
- Default: `3`

The maximum number of reminders to send for each invitation, between 1 and 5.



reminder_email_template_id
--------------------------

- Type: `string`

The ID of the email template to send reminders with, taken from a `descope_email_template` resource
with its `method` set to `invite`. An empty value (the default) sends reminders with the same template
as the invitation, and the value `System` selects the built-in System template.



reminder_text_template_id
-------------------------

- Type: `string`

The ID of the invitation SMS template to send reminders with. An empty value (the default) sends
reminders with the same template as the invitation, and the value `System` selects the built-in
System template.
