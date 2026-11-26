
SessionMigrationSettings
========================



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



vendor
------

- Type: `string`

The vendor to migrate user sessions from, either `auth0` or `okta`. Session migration is
disabled when this isn't set.



client_id
---------

- Type: `string`

The client ID from the vendor's authentication configuration.



domain
------

- Type: `string`

The domain from the vendor's authentication configuration. Only used with the `auth0` vendor.



audience
--------

- Type: `string`

An optional audience value from the vendor's authentication configuration. Only used with the
`auth0` vendor.



issuer
------

- Type: `string`

The issuer from the vendor's authentication configuration. Only used with the `okta` vendor.



api_token
---------

- Type: `secret`

The API token for the vendor's management API. Only used with the `okta` vendor. This value is
write-only and never read back from the server.



loginid_matched_attributes
--------------------------

- Type: `set` of `string`

Which attributes from the vendor's user should be used to match the Descope user's login ID
(e.g. `email`).



user_sync_type
--------------

- Type: `string`

The type of user synchronization to perform, either `matchOnly` to only match existing users or
`jit` for just-in-time provisioning.



user_mapping
------------

- Type: `list` of `settings.UserMappingItem`

A list of attribute mappings from the external vendor's user to Descope user attributes.





UserMappingItem
===============



external_key
------------

- Type: `string` (required)

The attribute key in the external vendor's user object.



descope_key
-----------

- Type: `string` (required)

The Descope user attribute to map the external key to.
