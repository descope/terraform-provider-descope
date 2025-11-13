
InviteSettings
==============



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



send_email
----------

- Type: `bool` 
- Default: `true`

Whether to send invitation emails to users.



send_text
---------

- Type: `bool` 

Whether to send invitation SMS messages to users.



email_service
-------------

- Type: `object` of `templates.EmailService` 

Settings related to sending invitation emails.
