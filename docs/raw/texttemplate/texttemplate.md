
TextTemplate
============



project_id
----------

- Type: `string` (required)

The ID of the project that the text template belongs to. Changing this value will require the
resource to be deleted and recreated.



method
------

- Type: `string` (required)

The authentication method the text template is used with, e.g. `magiclink` or `otp`. Changing this value will require the resource to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the text template that's unique among the templates of the same authentication method.



body
----

- Type: `string` (required)

The body of text messages sent with this template. Template macros such as `{{.code}}` can be used to insert dynamic values.
