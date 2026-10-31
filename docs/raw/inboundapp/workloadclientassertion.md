
WorkloadClientAssertion
=======================



federation_type
---------------

- Type: `string` (required)

The kind of issuer this entry trusts: `spiffe` for a SPIFFE trust domain, or `oidc` for an OIDC issuer.



trust_domain
------------

- Type: `string`

The SPIFFE trust domain, such as `prod.example.com`. Required when `federation_type` is `spiffe`.



issuer
------

- Type: `string`

The HTTPS issuer URL. Required when `federation_type` is `oidc`. For a `spiffe` entry it is the trust
domain's OIDC Discovery Provider and replaces the bundle endpoint in `jwks_uri`.



jwks_uri
--------

- Type: `string`

For a `spiffe` entry, the HTTPS URL of the trust domain's bundle endpoint, required unless `issuer` is set.
For an `oidc` entry, an optional JWKS URL that overrides the one found by discovery.



allowed_subjects
----------------

- Type: `list` of `string` (required)

The workloads this app accepts, matched against the assertion's subject. For a `spiffe` entry each value
is an exact SPIFFE ID or one ending with `/*` to accept the workloads under that path. For an `oidc`
entry each value is an anchored regular expression.



sign_algorithm
--------------

- Type: `string`

Pins the JWS algorithm of the assertion, such as `RS256` or `ES256`. Leave empty to accept every
supported asymmetric algorithm.
