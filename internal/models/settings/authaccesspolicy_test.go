package settings_test

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthAccessPolicy(t *testing.T) {
	proj := testacc.Project(t)
	project := `project_id = ` + proj.Path() + `.id`
	p := testacc.AuthAccessPolicy(t)
	rules := map[string]*testacc.Resource{}
	for _, id := range []string{"a", "b", "c", "d", "e"} {
		r := testacc.AuthAccessRule(t)
		r.ID = id
		r.Name += "-" + id
		rules[id] = r
	}
	config := func(defaultAction string, order []string, extra ...string) string {
		s := proj.Config()
		ids := []string{}
		for _, id := range order {
			s += rules[id].Config(project, `effect = "allow"`)
			ids = append(ids, rules[id].Path()+".id")
		}
		for _, id := range extra {
			s += rules[id].Config(project, `effect = "block"`)
		}
		return s + p.Block(project, `default_action = "`+defaultAction+`"`, `rule_ids = [`+strings.Join(ids, ", ")+`]`)
	}
	testacc.Run(t,
		// create rules and their order
		resource.TestStep{
			Config: config("block", []string{"a", "b", "c"}),
			Check:  resource.ComposeTestCheckFunc(p.Check(map[string]any{"default_action": "block"}), resource.TestCheckResourceAttrPair(p.Path(), "id", proj.Path(), "id"), orderChecks(p, rules, "a", "b", "c")),
		},
		// reorder without touching the rules
		resource.TestStep{
			Config: config("block", []string{"c", "a", "b"}),
			Check:  orderChecks(p, rules, "c", "a", "b"),
		},
		// delete a rule along with its reference
		resource.TestStep{
			Config: config("block", []string{"c", "a"}),
			Check:  orderChecks(p, rules, "c", "a"),
		},
		// insert a new rule in the middle
		resource.TestStep{
			Config: config("allow", []string{"c", "d", "a"}),
			Check:  resource.ComposeTestCheckFunc(p.Check(map[string]any{"default_action": "allow"}), orderChecks(p, rules, "c", "d", "a")),
		},
		// a rule left out of the order is created but leaves the policy out of sync
		resource.TestStep{
			Config:             config("allow", []string{"c", "d", "a"}, "e"),
			ExpectNonEmptyPlan: true,
		},
		// applying the policy without the rule fails
		resource.TestStep{
			Config:      config("allow", []string{"c", "d", "a"}, "e"),
			ExpectError: regexp.MustCompile(`Missing access rules`),
		},
		// listing the rule fixes it
		resource.TestStep{
			Config: config("allow", []string{"c", "d", "a", "e"}),
			Check:  orderChecks(p, rules, "c", "d", "a", "e"),
		},
		// import
		resource.TestStep{
			ResourceName:      p.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(p.Path(), "project_id"),
		},
		// duplicates are rejected at plan time
		resource.TestStep{
			Config:      config("allow", []string{"c", "d", "a", "e", "a"}),
			ExpectError: regexp.MustCompile(`unique`),
		},
		// remove every rule and restore the default action
		resource.TestStep{
			Config: proj.Config() + p.Block(project),
			Check:  p.Check(map[string]any{"default_action": "allow", "rule_ids.#": 0}),
		},
	)
}

func orderChecks(p *testacc.Resource, rules map[string]*testacc.Resource, order ...string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(p.Path(), "rule_ids.#", strconv.Itoa(len(order)))}
	for i, id := range order {
		checks = append(checks, resource.TestCheckResourceAttrPair(p.Path(), "rule_ids."+strconv.Itoa(i), rules[id].Path(), "id"))
	}
	return resource.ComposeTestCheckFunc(checks...)
}
