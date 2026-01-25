
RBac
====



is_company_admin
----------------

- Type: `bool` 

Whether this user has company-wide admin access. When set to `true`, the user
cannot have `tag_roles` or `project_roles`.



tag_roles
---------

- Type: `list` of `descoper.DescoperTagRole` 

A list of role assignments for projects matching specific tags.



project_roles
-------------

- Type: `list` of `descoper.DescoperProjectRole` 

A list of role assignments for specific projects by their ID.
