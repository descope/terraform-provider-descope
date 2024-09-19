
Project
=======



name
----

- Type: `string` (required)

The name of the Descope project.



environment
-----------

- Type: `string` 

This can be set to `production` to mark production projects, otherwise this should be left unset for development or staging projects.



project_settings
----------------

- Type: `object` of `settings.Settings` 

General settings for the Descope project.



authentication
--------------

- Type: `object` of `authentication.Authentication` 

Settings for each authentication method.



authorization
-------------

- Type: `object` of `authorization.Authorization` 

// description for authorization



attributes
----------

- Type: `object` of `attributes.Attributes` 

// description for attributes



connectors
----------

- Type: `object` of `connectors.Connectors` 

// description for connectors



applications
------------

- Type: `object` of `applications.Application` 

// description for applications



jwt_templates
-------------

- Type: `object` of `jwttemplates.JWTTemplates` 

// description for jwt_templates



styles
------

- Type: `object` of `flows.Styles` 

// description for styles



flows
-----

- Type: `map` of `flows.Flow` 

// description for flows
