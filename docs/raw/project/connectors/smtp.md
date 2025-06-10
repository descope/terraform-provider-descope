
SMTP
====



sender
------

- Type: `object` of `connectors.SenderField` (required)

The sender details that should be displayed in the email message.



server
------

- Type: `object` of `connectors.ServerField` (required)

SMTP server connection details including hostname and port.



authentication
--------------

- Type: `object` of `connectors.SMTPAuthField` (required)

SMTP server authentication credentials and method.



use_static_ips
--------------

- Type: `bool` 

Whether the connector should send all requests from specific static IPs.





SMTPAuthField
=============



username
--------

- Type: `string` (required)

Username for SMTP server authentication.



password
--------

- Type: `secret` (required)

Password for SMTP server authentication.



method
------

- Type: `string` 
- Default: `"plain"`

SMTP authentication method (`plain` or `login`).
