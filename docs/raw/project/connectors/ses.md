
SES
====



auth_type
---------

- Type: `string` 
- Default: `"credentials"`

The authentication type to use.



access_key_id
-------------

- Type: `secret` 

AWS Access key ID.



secret
------

- Type: `secret` 

AWS Secret Access Key.



role_arn
--------

- Type: `string` 

The Amazon Resource Name (ARN) of the role to assume.



external_id
-----------

- Type: `string` 

The external ID to use when assuming the role.



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
