
SAML
====



id
----

- Type: `string` 

// description for id



name
----

- Type: `string` (required)

// description for name



description
-----------

- Type: `string` 

// description for description



logo
----

- Type: `string` 

// description for logo



disabled
--------

- Type: `bool` 

// description for disabled



login_page_url
--------------

- Type: `string` 

// description for login_page_url



dynamic_configuration
---------------------

- Type: `object` of `applications.DynamicConfiguration` 

// description for dynamic_configuration



manual_configuration
--------------------

- Type: `object` of `applications.ManualConfiguration` 

// description for manual_configuration



acs_allowed_callback_urls
-------------------------

- Type: `list` of `string` 

// description for acs_allowed_callback_urls



subject_name_id_type
--------------------

- Type: `string` 

// description for subject_name_id_type



subject_name_id_format
----------------------

- Type: `string` 

// description for subject_name_id_format



default_relay_state
-------------------

- Type: `string` 

// description for default_relay_state



attribute_mapping
-----------------

- Type: `list` of `applications.AttributeMapping` 

// description for attribute_mapping





AttributeMapping
================



name
----

- Type: `string` (required)

// description for name



value
-----

- Type: `string` (required)

// description for value





DynamicConfiguration
====================



metadata_url
------------

- Type: `string` (required)

// description for metadata_url





ManualConfiguration
===================



acs_url
-------

- Type: `string` (required)

// description for acs_url



entity_id
---------

- Type: `string` (required)

// description for entity_id



certificate
-----------

- Type: `string` (required)

// description for certificate
