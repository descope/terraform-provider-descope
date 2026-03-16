
ApplicationScope
================



name
----

- Type: `string` (required)

A name for the scope.



description
-----------

- Type: `string`

A description for the scope.



optional
--------

- Type: `bool`

Whether this scope is optional. When `false`, the scope is mandatory and must be granted during
authorization. When `true`, the user may choose to withhold it.



values
------

- Type: `list` of `string`

The roles, user attributes, or connection scope identifiers associated with this scope, depending
on whether it is a permissions, attributes, or connections scope.
