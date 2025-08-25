
Password
========



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



expiration
----------

- Type: `bool` 

Whether users are required to change their password periodically.



expiration_weeks
----------------

- Type: `int` 

The number of weeks after which a user's password expires and they need to replace it.



lock
----

- Type: `bool` 

Whether the user account should be locked after a specified number of failed login attempts.



lock_attempts
-------------

- Type: `int` 

The number of failed login attempts allowed before an account is locked.



lowercase
---------

- Type: `bool` 

Whether passwords must contain at least one lowercase letter.



min_length
----------

- Type: `int` 

The minimum length of the password that users are required to use. The maximum length is always `64`.



non_alphanumeric
----------------

- Type: `bool` 

Whether passwords must contain at least one non-alphanumeric character (e.g. `!`, `@`, `#`).



number
------

- Type: `bool` 

Whether passwords must contain at least one number.



reuse
-----

- Type: `bool` 

Whether to forbid password reuse when users change their password.



reuse_amount
------------

- Type: `int` 

The number of previous passwords whose hashes are kept to prevent users from
reusing old passwords.



uppercase
---------

- Type: `bool` 

Whether passwords must contain at least one uppercase letter.



mask_errors
-----------

- Type: `bool` 

Prevents information about user accounts from being revealed in error messages, e.g.,
whether a user already exists.



email_service
-------------

- Type: `object` of `templates.EmailService` 

Settings related to sending password reset emails as part of the password feature.
