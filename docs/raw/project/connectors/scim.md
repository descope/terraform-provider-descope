
SCIM
====



disabled
--------

- Type: `bool`

Whether to disable this SCIM connector. When disabled, provisioning events will not be
sent to the configured endpoint.



federated_app_id
----------------

- Type: `string` (required)

The ID of the federated SSO application this SCIM connector is associated with.



base_url
--------

- Type: `string` (required)

The base URL of the SCIM v2 endpoint that user provisioning events will be sent to.



authentication
--------------

- Type: `object` of `connectors.HTTPAuthField`

Authentication credentials used when sending requests to the SCIM endpoint.



headers
-------

- Type: `map` of `string`

Custom HTTP headers to send with each provisioning request.



hmac_secret
-----------

- Type: `secret`

HMAC is a method for message signing with a symmetrical key. This secret will be
used to sign the base64 encoded payload, and the resulting signature will be sent
in the `x-descope-webhook-s256` header. The receiving service should use this
secret to verify the integrity and authenticity of the payload by checking the
provided signature.



insecure
--------

- Type: `bool`

Will ignore certificate errors raised by the client.
