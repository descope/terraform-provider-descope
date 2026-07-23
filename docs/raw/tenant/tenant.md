
Tenant
======



project_id
----------

- Type: `string` (required)

The Descope project that owns the tenant. Changing this value requires replacing the tenant resource.



id
----

- Type: `string` 

An optional tenant identifier. If omitted, Descope generates an identifier. Changing this value requires replacing the tenant resource.



name
----

- Type: `string` (required)

The tenant name.



self_provisioning_domains
-------------------------

- Type: `list` of `string` 

Domains whose users may self-provision into the tenant.



disabled
--------

- Type: `bool` 

Whether the tenant is disabled.



enforce_sso
-----------

- Type: `bool` 

Whether SSO is enforced for the tenant.



enforce_sso_exclusions
----------------------

- Type: `list` of `string` 

Login IDs excluded from tenant SSO enforcement.



federated_application_ids
-------------------------

- Type: `list` of `string` 

Federated application identifiers associated with the tenant.



parent
------

- Type: `string` 

The parent tenant identifier. Changing this value requires replacing the tenant resource.



role_inheritance
----------------

- Type: `string` 

Role inheritance from the parent tenant. Valid values are `none` and `userOnly`.
