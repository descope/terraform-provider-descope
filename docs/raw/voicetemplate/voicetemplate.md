
VoiceTemplate
=============



project_id
----------

- Type: `string` (required)

The ID of the project that the voice template belongs to. Changing this value will require the
resource to be deleted and recreated.



method
------

- Type: `string` (required)

The authentication method the voice template is used with. Currently only `otp` is supported.
Changing this value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the voice template that's unique among the templates of the same authentication method.



body
----

- Type: `string` (required)

The text that's read aloud in voice calls made with this template. Template macros such as `{{.code}}` can be used to insert dynamic values.
