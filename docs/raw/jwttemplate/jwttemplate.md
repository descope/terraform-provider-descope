
JWTTemplate
===========



project_id
----------

- Type: `string` (required)

The ID of the project that this JWT template belongs to. Changing this value will require
the resource to be deleted and recreated.



name
----

- Type: `string` (required)

Name of the JWT template. Must be unique among the project's JWT templates of both types.



description
-----------

- Type: `string`

Description of the JWT template.



type
----

- Type: `string` (required)

The kind of tokens this template applies to - `user` for user session JWTs or `key` for
access key JWTs. This value can be changed in place without recreating the template.



issuer_type
-----------

- Type: `string`
- Default: `"legacy"`

The issuer format for the JWT - `legacy`, `inbound` or `federated`. The `federated` issuer
type cannot be combined with the `conformance_issuer` attribute.



auth_schema
-----------

- Type: `string`
- Default: `"default"`

The authorization claims format - `default`, `tenantOnly` or `none`.
Read more about schema types [here](https://docs.descope.com/project-settings/jwt-templates).



empty_claim_policy
------------------

- Type: `string`
- Default: `"none"`

Policy for empty claims - `none`, `nil` or `delete`.



auto_tenant_claim
-----------------

- Type: `bool`

When a user is associated with a single tenant, the tenant will be set as the user's
active tenant, using the `dct` (Descope Current Tenant) claim in their JWT.



conformance_issuer
------------------

- Type: `bool`

Whether to use OIDC conformance for the JWT issuer field.



enforce_issuer
--------------

- Type: `bool`

Whether to enforce that the JWT issuer matches the project configuration.



exclude_permission_claim
------------------------

- Type: `bool`

When enabled, permissions will not be included in the JWT token.



override_subject_claim
----------------------

- Type: `bool`

Switching on will allow you to add a custom subject claim to the JWT. A default new `dsub` claim
will be added with the user ID.



add_jti_claim
-------------

- Type: `bool`

When enabled, a unique JWT ID (jti) claim will be added to the token for tracking and preventing replay attacks.



template
--------

- Type: `string` (required)

The JSON template defining the structure and claims of the JWT token. This is expected
to be a valid JSON object given as a `string` value, preferably with the `jsonencode`
function so the value converges with what the Descope service returns.
