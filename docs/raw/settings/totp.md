
TOTPSettings
============



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using TOTP authentication directly via API and SDK calls. Note
that this does not affect authentication flows that are configured to use TOTP.



service_label
-------------

- Type: `string`

The template for the service issuer label (issuer) shown in the authenticator app.



account_label
-------------

- Type: `string`

The template for the account name shown in the authenticator app, for example `{{user.email}}`
or `{{user.customAttributes.employeeId}}`. Supports `{{user.email}}`, `{{user.phone}}`,
`{{user.name}}`, `{{user.loginId}}` and `{{user.customAttributes.<name>}}`. When empty, or when
it renders empty for a user, the user's login ID is used.
