
SMTP
====



sender
------

- Type: `object` of `connectors.SenderField` (required)

The sender details that should be displayed in the email message.



server
------

- Type: `object` of `connectors.ServerField` (required)

// description for server



authentication
--------------

- Type: `object` of `connectors.SMTPAuthField` (required)

// description for authentication



use_static_ips
--------------

- Type: `bool` 

Whether the connector should send all requests from specific static IPs.





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
