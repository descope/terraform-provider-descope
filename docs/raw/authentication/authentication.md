
Authentication
==============



otp
----

- Type: `object` of `authentication.OTP` 

A dynamically generated set of numbers, granting the user one-time access.



magic_link
----------

- Type: `object` of `authentication.MagicLink` 

An authentication method where a user receives a unique link via email to log in.



enchanted_link
--------------

- Type: `object` of `authentication.EnchantedLink` 

An enhanced and more secure version of Magic Link, enabling users to start the authentication
process on one device and execute the verification on another.



embedded_link
-------------

- Type: `object` of `authentication.EmbeddedLink` 

Make the authentication experience smoother for the user by generating their initial token in a
way that does not require the end user to initiate the process, requiring only verification.



password
--------

- Type: `object` of `authentication.Password` 

The classic username and password combination used for authentication.



oauth
-----

- Type: `object` of `authentication.OAuth` 

Authentication using Open Authorization, which allows users to authenticate with various external
services.



sso
----

- Type: `object` of `authentication.SSO` 

Single Sign-On (SSO) authentication method that enables users to access multiple applications with
a single set of credentials.



totp
----

- Type: `object` of `authentication.TOTP` 

A one-time code generated for the user using a shared secret and time.



passkeys
--------

- Type: `object` of `authentication.Passkeys` 

Device-based passwordless authentication, using fingerprint, face scan, and more.
