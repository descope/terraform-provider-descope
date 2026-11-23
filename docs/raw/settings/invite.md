
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



invite_reminders
----------------

- Type: `bool`

Whether to automatically resend the invitation to users who have not accepted it yet. Requires a
custom email or SMS connector. Only invitations sent after this is turned on get reminders, and a
manual resend restarts them.



invite_reminder_interval
------------------------

- Type: `duration`

The time between reminders, measured from the last invitation sent to the user. Use values such
as "2 days", "1 week", etc. The minimum value is "1 day".



invite_reminder_max_count
-------------------------

- Type: `int`
- Default: `3`

The maximum number of reminders sent per invitation, from 1 to 5.



invite_reminder_email_template_id
---------------------------------

- Type: `string`

The ID of the email template to send reminders with, taken from a `descope_email_template` resource
with its `method` set to `invite`. An empty value (the default) sends reminders with the template the
invitation was sent with.
