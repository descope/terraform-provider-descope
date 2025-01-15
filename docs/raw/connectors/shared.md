
SenderField
===========



email
-----

- Type: `string` (required)

// description for email



name
----

- Type: `string` 

// description for name





ServerField
===========



host
----

- Type: `string` (required)

// description for host



port
----

- Type: `int` 
- Default: `25`

// description for port





AuditFilterField
================



key
----

- Type: `string` (required)

// description for key



operator
--------

- Type: `string` (required)

// description for operator



values
------

- Type: `list` of `string` (required)

// description for values





HTTPAuthField
=============



bearer_token
------------

- Type: `secret` 

// description for bearer_token



basic
-----

- Type: `object` of `connectors.HTTPAuthBasicField` 

// description for basic



api_key
-------

- Type: `object` of `connectors.HTTPAuthAPIKeyField` 

// description for api_key





HTTPAuthBasicField
==================



username
--------

- Type: `string` (required)

// description for username



password
--------

- Type: `secret` (required)

// description for password





HTTPAuthAPIKeyField
===================



key
----

- Type: `string` (required)

// description for key



token
-----

- Type: `secret` (required)

// description for token
