Custom Language
===============



project_id
----------

- Type: `string` (required)

The ID of the Descope project this custom language belongs to. Changing this value will require the
resource to be deleted and recreated.



language
--------

- Type: `string` (required)

The language code sent by clients at runtime to request this language (e.g. `phl`). Changing this
value will require the resource to be deleted and recreated.



region
------

- Type: `string`

An optional region subtag appended to the language code to distinguish variants, such as `en-UK`.
Changing this value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

The human-readable display name shown for this language in the Descope console.
