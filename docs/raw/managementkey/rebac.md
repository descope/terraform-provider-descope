
ReBac
=====



company_roles
-------------

- Type: `set` of `string` 

A list of company-level role names that are granted to the management key. This
attribute is mutually exclusive with `tag_roles` and `project_roles`.



project_roles
-------------

- Type: `list` of `managementkey.ProjectRole` 

A list of project-level role names that are granted to the management key for
specific projects by their project ID.



tag_roles
---------

- Type: `list` of `managementkey.TagRole` 

A list of project-level role names that are granted to the management key for
all projects that have a specific tag.
