
AppPermission
=============



project_id
----------

- Type: `string` (required)

The ID of the project that this permission belongs to. Changing this value will require the
resource to be deleted and recreated.



app_id
------

- Type: `string` (required)

The ID of the federated application that this permission is scoped to. Reference the `id`
attribute of the application resource (e.g., `descope_oidc_app.example.id`). Changing this
value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

The name of the permission. Permission names must be unique per application, and a permission
can be renamed freely without affecting roles that reference it. Note that if a permission is
deleted, it is silently removed from any application roles that still reference it.



description
-----------

- Type: `string`

An optional description for the permission.
