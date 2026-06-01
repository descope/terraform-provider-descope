
SessionMigration
================



vendor
------

- Type: `string`

The name of the vendor the sessions are migrated from, in all lowercase.



client_id
---------

- Type: `string`

The unique client ID for the vendor.



domain
------

- Type: `string`

The domain value if needed by the vendor.



audience
--------

- Type: `string`

The audience value if needed by the vendor.



issuer
------

- Type: `string`

An issuer URL if needed by the vendor.



api_token
---------

- Type: `secret`

An API token for the vendor, required when `vendor` is set to `okta`.



loginid_matched_attributes
--------------------------

- Type: `set` of `string`

A set of attributes from the vendor's user that should be used to match with
the Descope user's login ID.



user_sync_type
--------------

- Type: `string`

The type of user synchronization to perform. Valid values are `matchOnly` (match existing users only) and `jit` (just-in-time provisioning).



user_mapping
------------

- Type: `list` of `settings.ExternalAuthUserMappingItem`

A list of attribute mappings from the external vendor's user to Descope user attributes.





ExternalAuthUserMappingItem
===========================



external_key
------------

- Type: `string` (required)

The attribute key in the external vendor's user object.



descope_key
-----------

- Type: `string` (required)

The Descope user attribute to map the external key to.
