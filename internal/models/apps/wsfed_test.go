package apps_test

import (
	"fmt"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWSFedApp(t *testing.T) {
	projectID := testacc.ProjectID(t)
	a := testacc.WSFedApp(t)
	r := testacc.Role(t)
	project := `project_id = "` + projectID + `"`

	var appID string
	captureID := func(s string) error {
		appID = s
		return nil
	}
	sameID := func(s string) error {
		if s != appID {
			return fmt.Errorf("expected id to be preserved as %s, got %s", appID, s)
		}
		return nil
	}

	createStep := resource.TestStep{
		Config: a.Config(project,
			`deletion_protection = false`,
			`realm = "urn:example"`,
			`reply_url = "https://sp.example.com/wsfed"`,
		),
		Check: a.Check(map[string]any{
			"id":             captureID,
			"project_id":     projectID,
			"name":           a.Name,
			"realm":          "urn:example",
			"reply_url":      "https://sp.example.com/wsfed",
			"login_page_url": testacc.AttributeIsSet,
		}),
	}

	a.Name += "-renamed"
	updateStep := resource.TestStep{
		Config: r.Config(project) + a.Config(project,
			`deletion_protection = false`,
			`realm = "urn:example2"`,
			`reply_url = "https://sp.example.com/wsfed2"`,
			`reply_allowed_callback_urls = ["https://sp.example.com/*"]`,
			`attribute_mapping = [{ name = "email", value = "user.email" }]`,
			`groups_mapping = [{
				name = "grp"
				type = "roles"
				filter_type = "roles"
				value = "grp-value"
				roles = [{ id = `+r.Path()+`.id, name = `+r.Path()+`.name }]
			}]`,
			`logout_redirect_url = "https://sp.example.com/logout"`,
			`error_redirect_url = "https://sp.example.com/error"`,
		),
		Check: a.Check(map[string]any{
			"id":                            sameID,
			"name":                          a.Name,
			"realm":                         "urn:example2",
			"reply_allowed_callback_urls":   []string{"https://sp.example.com/*"},
			"attribute_mapping.0.name":      "email",
			"groups_mapping.0.name":         "grp",
			"groups_mapping.0.roles.0.name": r.Name,
			"logout_redirect_url":           "https://sp.example.com/logout",
			"error_redirect_url":            "https://sp.example.com/error",
		}),
	}

	importStep := resource.TestStep{
		ResourceName:            a.Path(),
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateIdFunc:       testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		ImportStateVerifyIgnore: []string{"deletion_protection"},
	}

	testacc.RunWithDestroyCheck(t, "descope_wsfed_app", createStep, updateStep, importStep)
}
