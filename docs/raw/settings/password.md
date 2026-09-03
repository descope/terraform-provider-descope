
PasswordSettings
================



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using password authentication directly via API and SDK calls.
Note that this does not affect authentication flows that are configured to use passwords.



min_length
----------

- Type: `int`
- Default: `8`

The minimum length of the password that users are required to use. The maximum length is always
`64`.



lowercase
---------

- Type: `bool`
- Default: `true`

Whether passwords must contain at least one lowercase letter.



uppercase
---------

- Type: `bool`
- Default: `true`

Whether passwords must contain at least one uppercase letter.



number
------

- Type: `bool`
- Default: `true`

Whether passwords must contain at least one number.



non_alphanumeric
----------------

- Type: `bool`
- Default: `true`

Whether passwords must contain at least one non-alphanumeric character (e.g. `!`, `@`, `#`).



any_letter
----------

- Type: `bool`

Whether passwords must contain at least one letter, either uppercase or lowercase.



disallowed_characters
---------------------

- Type: `string`

Reject passwords containing any of these characters. Each character in the string is treated as a
forbidden literal (e.g., `"'"` to reject single and double quotes).



disallow_email_match
--------------------

- Type: `bool`

Whether to reject passwords that match the user's email address or its local-part (the segment
before `@`), case-insensitively. The check is skipped if the user's email is not known at validation
time.



expiration
----------

- Type: `bool`

Whether users are required to change their password periodically.



expiration_weeks
----------------

- Type: `int`
- Default: `20`

The number of weeks after which a user's password expires and they need to replace it.



reuse
-----

- Type: `bool`

Whether to forbid password reuse when users change their password.



reuse_amount
------------

- Type: `int`
- Default: `10`

The number of previous passwords whose hashes are kept to prevent users from reusing old passwords.



lock
----

- Type: `bool`

Whether the user account should be locked after a specified number of failed login attempts.



lock_attempts
-------------

- Type: `int`
- Default: `5`

The number of failed login attempts allowed before an account is locked.



temporary_lock
--------------

- Type: `bool`

Whether the user account should be temporarily locked after a specified number of failed login
attempts.



temporary_lock_attempts
-----------------------

- Type: `int`
- Default: `3`

The number of failed login attempts allowed before an account is temporarily locked.



temporary_lock_duration
-----------------------

- Type: `duration`

The amount of time before the user can sign in again after the account is temporarily locked.



enforce_strength
----------------

- Type: `string`
- Default: `"none"`

Use zxcvbn to calculate the strength of a given password and enforce a minimum level of strength.



mask_errors
-----------

- Type: `bool`
- Default: `true`

Prevents information about user accounts from being revealed in error messages, e.g., whether a user
already exists.



email_service
-------------

- Type: `object` of `settings.EmailServiceRef`

Settings related to sending password reset emails as part of the password feature.



email_template_id
-----------------

- Type: `string`

The ID of the email template for password reset emails, taken from a `descope_email_template`
resource with its `method` set to `password`. The same template serves both the magic link and
enchanted link reset flows. An empty value (the default) selects the built-in System template.
