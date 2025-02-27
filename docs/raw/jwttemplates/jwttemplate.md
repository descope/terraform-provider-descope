
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

The authorization claims format - `default`, `tenantOnly` or `none`. Read more about schema types [here](https://docs.descope.com/project-settings/jwt-templates).



empty_claim_policy
------------------

- Type: `string` 
- Default: `"none"`

Policy for empty claims - `none`, `nil` or `delete`.



empty_claim_policy
------------------

- Type: `string` 
- Default: `"none"`

Policy for empty claims - "none", "nil" or "delete".



conformance_issuer
------------------

- Type: `bool` 

// description for conformance_issuer



enforce_issuer
--------------

- Type: `bool` 

// description for enforce_issuer



template
--------

- Type: `string` (required)

// description for template
