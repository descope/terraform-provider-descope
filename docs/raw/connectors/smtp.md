
SMTP
====



sender
------

- Type: `object` of `connectors.SenderField` (required)

// description for sender



server
------

- Type: `object` of `connectors.ServerField` (required)

// description for server



authentication
--------------

- Type: `object` of `connectors.SMTPAuthField` (required)

// description for authentication





SMTPAuthField
=============



username
--------

- Type: `string` (required)

// description for username



password
--------

- Type: `secret` (required)

// description for password



method
------

- Type: `string` 
- Default: `"plain"`

// description for method