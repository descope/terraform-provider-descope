
FlowTranslations
================



project_id
----------

- Type: `string` (required)

The ID of the project that the flow belongs to. Changing this value will require the resource to be
deleted and recreated.



flow_id
-------

- Type: `string` (required)

The ID of the flow whose texts are translated (e.g., `sign-up-or-in`). Changing this value will
require the resource to be deleted and recreated.



languages
---------

- Type: `string` (required)

A JSON document with the translated texts, the same document as the `languages` object in an
export of the flow's localization: an object with a key for each language, each holding the translated
texts as an object of strings keyed by the text's translation key. The document replaces all of the
flow's translations, and destroying the resource removes them. Only list the languages that the flow's
translation settings name as its target languages, since Descope keeps no others.
