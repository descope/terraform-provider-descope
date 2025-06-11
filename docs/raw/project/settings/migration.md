
SessionMigration
================



vendor
------

- Type: `string` (required)

The name of the vendor the sessions are migrated from, in all lowercase.



client_id
---------

- Type: `string` (required)

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



loginid_matched_attributes
--------------------------

- Type: `set` of `string` (required)

A set of attributes from the vendor's user that should be used to match with
the Descope user's login ID.
