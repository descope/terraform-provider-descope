
IDJAGTarget
===========



authorization_server
--------------------

- Type: `string` (required)

The issuer of the external authorization server, used as the assertion audience.



target_tenant
-------------

- Type: `string`

The tenant ID as the target authorization server knows it. Empty reuses the user's current tenant ID.



continuation_trust
------------------

- Type: `list` of `apiresource.ContinuationTrust`

Issuers outside this project whose identity-continuation assertions are accepted for this
authorization server. Applies to the whole project. At most 20 entries.





ContinuationTrust
=================



issuer
------

- Type: `string` (required)

The https issuer URL of the trusted authorization server, matched exactly against the `iss` claim
of the assertion. It cannot be this project's own issuer.



jwks_uri
--------

- Type: `string` (required)

The https URL of the issuer's JWKS. Only keys published there can sign a trusted assertion.



sign_algorithm
--------------

- Type: `string` (required)

The signing algorithm of the assertion: `RS256`, `ES256`, `ES384` or `ES512`.



actor_client_id
---------------

- Type: `string` (required)

The client ID of the only actor allowed to present assertions from this issuer. Each issuer and
actor client pair may appear once.



enabled
-------

- Type: `bool`
- Default: `true`

Whether the entry is trusted. A disabled entry is kept but never trusted.
