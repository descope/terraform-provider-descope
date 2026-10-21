
Governance
==========



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



configured
----------

- Type: `bool`

Whether the Agentic Governance Suite has been set up for the project.



auto_approval
-------------

- Type: `bool`

Setting this to `true` approves access requests automatically instead of waiting for an admin to
review them.



suite_disabled
--------------

- Type: `bool`

Setting this to `true` turns off the Agentic Governance Suite for the project.
