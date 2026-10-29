
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



email_connector_id
------------------

- Type: `string`

The ID of an email connector to use for sending emails. An empty value (the default) selects the built-in Descope delivery service.



email_template_id
-----------------

- Type: `string`

The ID of the email template for password reset emails, taken from a `descope_email_template`
resource with its `method` set to `password`. The same template serves both the magic link and
enchanted link reset flows. An empty value (the default) selects the built-in System template.



role_policies
-------------

- Type: `list` of `settings.PasswordRolePolicy`

Stricter password requirements for users who hold specific roles. When a user holds any of the
roles of a policy, at the project level or in any of their tenants, each rule of the policy is
combined with the matching rule of the password settings that apply to the user, and the strictest
value wins. A policy can therefore only make the requirements stricter, and when a user matches
more than one policy all of them apply. The policies are enforced when a signed-in user sets or
changes their own password, and not at sign-up or to passwords set for a user by an administrator.
Existing passwords are not checked again when roles or policies change.





PasswordRolePolicy
==================



roles
-----

- Type: `set` of `string` (required)

The names of the roles the policy applies to. Reference the `name` attribute of `descope_role`
resources (e.g., `descope_role.admin.name`) so that Terraform orders the operations correctly.
Each role must exist when the policy is applied.



min_length
----------

- Type: `int`

The minimum length of the password, up to `64`. The default value of `0` adds no requirement.



lowercase
---------

- Type: `bool`

Whether passwords must contain at least one lowercase letter.



uppercase
---------

- Type: `bool`

Whether passwords must contain at least one uppercase letter.



number
------

- Type: `bool`

Whether passwords must contain at least one number.



non_alphanumeric
----------------

- Type: `bool`

Whether passwords must contain at least one non-alphanumeric character (e.g. `!`, `@`, `#`).



any_letter
----------

- Type: `bool`

Whether passwords must contain at least one letter, either uppercase or lowercase.



disallowed_characters
---------------------

- Type: `string`

Reject passwords containing any of these characters, in addition to the characters that the
password settings disallow. Each character in the string is treated as a forbidden literal.



disallow_email_match
--------------------

- Type: `bool`

Whether to reject passwords that match the user's email address or its local-part (the segment
before `@`), case-insensitively.



enforce_strength
----------------

- Type: `string`
- Default: `"none"`

Use zxcvbn to calculate the strength of a given password and enforce a minimum level of strength.
The default value of `none` adds no requirement.
