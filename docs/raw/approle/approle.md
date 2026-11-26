
AppRole
=======



project_id
----------

- Type: `string` (required)

The ID of the project that this role belongs to. Changing this value will require the
resource to be deleted and recreated.



app_id
------

- Type: `string` (required)

The ID of the federated application that this role is scoped to. Reference the `id`
attribute of the application resource (e.g., `descope_oidc_app.example.id`). Changing this
value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

The name of the role. Role names must be unique per application, and a role is renamed in
place, so user assignments and service account references to it are preserved.



description
-----------

- Type: `string`

An optional description for the role.



permission_ids
--------------

- Type: `set` of `string`

The IDs of the application permissions this role grants. Reference the `id` attribute of
`descope_app_permission` resources on the same application (e.g.,
`descope_app_permission.example.id`) so that Terraform orders the operations correctly.



role_mappings
-------------

- Type: `set` of `string`

The IDs of project-level roles that map users to this application role at login. Reference
the `id` attribute of `descope_role` resources (e.g., `descope_role.example.id`).
