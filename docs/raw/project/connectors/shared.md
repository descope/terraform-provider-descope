
SenderField
===========



email
-----

- Type: `string` (required)

The email address that will appear as the sender of the email.



name
----

- Type: `string` 

The display name that will appear as the sender of the email.





ServerField
===========



host
----

- Type: `string` (required)

The hostname or IP address of the SMTP server.



port
----

- Type: `int` 
- Default: `25`

The port number to connect to on the SMTP server.





AuditFilterField
================



key
----

- Type: `string` (required)

The field name to filter on (either 'actions' or 'tenants').



operator
--------

- Type: `string` (required)

The filter operation to apply ('includes' or 'excludes').



values
------

- Type: `list` of `string` (required)

The list of values to match against for the filter.





HTTPAuthField
=============



bearer_token
------------

- Type: `secret` 

Bearer token for HTTP authentication.



basic
-----

- Type: `object` of `connectors.HTTPAuthBasicField` 

Basic authentication credentials (username and password).



api_key
-------

- Type: `object` of `connectors.HTTPAuthAPIKeyField` 

API key authentication configuration.





HTTPAuthBasicField
==================



username
--------

- Type: `string` (required)

Username for basic HTTP authentication.



password
--------

- Type: `secret` (required)

Password for basic HTTP authentication.





HTTPAuthAPIKeyField
===================



key
----

- Type: `string` (required)

The API key.



token
-----

- Type: `secret` (required)

The API secret.
