
RBac
====



is_company_admin
----------------

- Type: `bool` 

Whether this descoper has company-wide admin access. When set to `true`,
the descoper cannot have `tag_roles` or `project_roles`.



project_roles
-------------

- Type: `list` of `descoper.DescoperProjectRole` 

A list of roles that are granted to the descoper for specific
projects by their project ID.



tag_roles
---------

- Type: `list` of `descoper.DescoperTagRole` 

A list of roles that are granted to the descoper for all projects
that have a specific tag.
