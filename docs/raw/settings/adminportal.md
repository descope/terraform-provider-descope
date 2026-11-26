
AdminPortal
===========



project_id
----------

- Type: `string` (required)

The ID of the project that this configuration belongs to. Changing this value will require the
resource to be deleted and recreated.



enabled
-------

- Type: `bool`

Whether the Admin Portal is enabled. At least one widget must be listed in the `widgets`
attribute when this is set to `true`.



style_id
--------

- Type: `string`

The name of the style the Admin Portal uses. Styles are defined in the project theme, which
can be managed with the `descope_styles` resource. When left empty, the default style is used.



widgets
-------

- Type: `list` of `settings.AdminPortalWidget`

The widgets to show in the Admin Portal, in display order. The `widget_id` values usually
reference widgets managed with the `descope_widget` resource, e.g.,
`widget_id = descope_widget.mywidget.widget_id`.





AdminPortalWidget
=================



widget_id
---------

- Type: `string` (required)

The unique identifier of the Widget



type
----

- Type: `string` (required)

The type of the Widget
