
ScopeClaimMapping
=================



scope
-----

- Type: `string` (required)

The scope this mapping applies to.



claims
------

- Type: `map` of `string`

The claims the scope contributes, mapping each claim name to the attribute it is taken from.



description
-----------

- Type: `string`

A description of what the scope grants access to.



use_project_mapping
-------------------

- Type: `bool`

Whether the entry inherits the project-wide claim mapping instead of the claims listed here.



mandatory
---------

- Type: `bool`

Whether the scope is non-optional. A token request that is denied a mandatory scope fails outright,
rather than succeeding with the subset of scopes that were granted.
