
Permission
==========



project_id
----------

- Type: `string` (required)

The ID of the project that this permission belongs to. Changing this value will require the
resource to be deleted and recreated.



name
----

- Type: `string` (required)

The name of the permission. Permission names must be unique in the project, and roles reference
permissions by name. The names of the system permissions that are created in every project are
reserved and cannot be used: `Impersonate`, `SSO Admin`, `Super User`, and `User Admin`.



description
-----------

- Type: `string`

An optional description for the permission.
