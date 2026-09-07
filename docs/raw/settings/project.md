
ProjectSettings
===============



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



app_url
-------

- Type: `string`

The URL where your application resides.



custom_domain
-------------

- Type: `string`

A custom CNAME domain for the project, used instead of the default Descope base URL for
authentication traffic. Must be a subdomain of the `app_url` domain.



approved_domains
----------------

- Type: `set` of `string`

List of approved domains that are allowed for redirect and verification URLs for different
authentication methods.



default_no_sso_apps
-------------------

- Type: `bool`

Whether new users should not be given access to any federated applications by default.



tenant_user_isolation
---------------------

- Type: `bool`

Isolate users per tenant so the same login ID is treated as a separate identity in each tenant,
with independent credentials and MFA state.



allow_auth_hosting_iframe_embedding
-----------------------------------

- Type: `bool`

When enabled, Descope-hosted flows can be displayed within an iframe on
your website. This modifies the security headers that typically prevent the page from
being embedded.



test_users_loginid_regexp
-------------------------

- Type: `string`

A regular expression pattern that, when a user is created with a matching login ID, will
automatically mark the user as a test user.



test_users_static_otp
---------------------

- Type: `string`

A static OTP code that can be used by test users instead of a real verification code. Must
be set together with `test_users_verifier_regexp`.



test_users_verifier_regexp
--------------------------

- Type: `string`

A regular expression pattern that determines which test user verifiers (email addresses or
phone numbers) the static OTP code applies to.
