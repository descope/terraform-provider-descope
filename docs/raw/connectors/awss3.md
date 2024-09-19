
AWSS3
=====



access_key_id
-------------

- Type: `secret` (required)

The unique AWS access key ID.



secret_access_key
-----------------

- Type: `secret` (required)

The secret AWS access key.



region
------

- Type: `string` (required)

The AWS S3 region, e.g. `us-east-1`.



bucket
------

- Type: `string` (required)

The AWS S3 bucket. This bucket should already exist for the connector to work.



audit_enabled
-------------

- Type: `bool` 
- Default: `true`

// description for audit_enabled



audit_filters
-------------

- Type: `string` 

// description for audit_filters



troubleshoot_log_enabled
------------------------

- Type: `bool` 

// description for troubleshoot_log_enabled
