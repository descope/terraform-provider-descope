
Role
====



project_id
----------

- Type: `string` (required)

The ID of the project that this role belongs to. Changing this value will require the resource
to be deleted and recreated.



name
----

- Type: `string` (required)

The name of the role. Role names must be unique in the project, and a role can be renamed
freely without affecting resources that reference it.



description
-----------

- Type: `string`

An optional description for the role.



permissions
-----------

- Type: `set` of `string`

The names of the permissions this role grants. Permissions are referenced by name, and their
existence is validated when the role is created or updated. Reference the `name` attribute of
`descope_permission` resources (e.g., `descope_permission.example.name`) so that Terraform
orders the operations correctly. Note that if a permission is deleted, it is silently removed
from any roles that still reference it.



default
-------

- Type: `bool`

Whether this role is a default role in the project.



private
-------

- Type: `bool`

Whether this role is private.
