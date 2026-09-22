
OAuthSettings
=============



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using OAuth authentication directly via API and SDK calls. Note
that this does not affect authentication flows that are configured to use OAuth.
