
ClientAssertion
===============



key_source
----------

- Type: `string` (required)

Where the key that signed the assertion comes from: `jwks` for a JWKS URL, or `publicKey` for a PEM
public key.



jwks_uri
--------

- Type: `string`

The HTTPS URL of the JWKS document. Required when `key_source` is `jwks`.



public_key
----------

- Type: `string`

A PEM encoded public key or certificate. Required when `key_source` is `publicKey`.



sign_algorithm
--------------

- Type: `string`

Pins the JWS algorithm of the assertion, such as `RS256` or `ES256`. Leave empty to accept every
supported asymmetric algorithm.
