
Engine
======



project_id
----------

- Type: `string` (required)

The ID of the Descope project this engine belongs to. Changing this value will require the
resource to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the engine.



created_time
------------

- Type: `int`

The creation time of the engine as a Unix timestamp.



secret
------

- Type: `secret`

The plaintext secret for the engine. This is only available after the engine is created and
cannot be retrieved later. Store this value securely as it is used to authenticate the engine.
