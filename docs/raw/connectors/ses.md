
SES
====



access_key_id
-------------

- Type: `string` (required)

AWS Access key ID.



secret
------

- Type: `secret` (required)

AWS Secret Access Key.



region
------

- Type: `string` (required)

AWS region to send requests to (e.g. `us-west-2`).



endpoint
--------

- Type: `string` 

An optional endpoint URL (hostname only or fully qualified URI).



sender
------

- Type: `object` of `connectors.SenderField` (required)

The sender details that should be displayed in the email message.
