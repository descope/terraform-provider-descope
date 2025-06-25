
JWTTemplate
===========



name
----

- Type: `string` (required)

Name of the JWT Template.



description
-----------

- Type: `string` 

Description of the JWT Template.



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



template
--------

- Type: `string` (required)

The JSON template defining the structure and claims of the JWT token. This is expected
to be a valid JSON object given as a `string` value.
