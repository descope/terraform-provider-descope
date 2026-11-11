
EmailVerificationSettings
=========================



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



enabled
-------

- Type: `bool`

Whether flows can sign users in with an email verification token, a proof from the user's email provider
that the browser hands the flow instead of a code. Email verification is in beta and can only be
enabled for projects it is available to.
