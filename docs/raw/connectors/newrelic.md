
NewRelic
========



api_key
-------

- Type: `secret` (required)

Ingest License Key of the account you want to report data to.



data_center
-----------

- Type: `string` 

The New Relic data center the account belongs to. Possible values are: `US`,
`EU`, `FedRAMP`. Default is `US`.



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



override_logs_prefix
--------------------

- Type: `bool` 

Enable this option to use a custom prefix for log fields.



logs_prefix
-----------

- Type: `string` 
- Default: `"descope."`

A custom prefix for all log fields. The default prefix is `descope.`.
