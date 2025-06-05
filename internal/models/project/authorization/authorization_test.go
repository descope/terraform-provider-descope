package authorization_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAuthorization(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				authorization = {
					roles = [
						{
							name = "App Developer"
							description = "Builds apps and uploads new beta builds"
							permissions = ["build-apps", "upload-builds", "install-builds"]
						},
						{
							name = "App Tester"
							description = "Installs and tests beta releases"
							permissions = ["install-builds"]
						},
					]
					permissions = [
						{
							name = "build-apps"
							description = "Allowed to build and sign applications"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Missing Permission`),
		},
		resource.TestStep{
			Config: p.Config(`
				authorization = {
					roles = [
						{
							name = "App Developer"
							description = "Builds apps and uploads new beta builds"
							permissions = ["build-apps", "upload-builds", "install-builds"]
						},
						{
							name = "App Tester"
							description = "Installs and tests beta releases"
							permissions = ["install-builds"]
						},
					]
					permissions = [
						{
							name = "build-apps"
							description = "Allowed to build and sign applications"
						},
						{
							name = "upload-builds"
							description = "Allowed to upload new releases"
						},
						{
							name = "install-builds"
							description = "Allowed to install beta releases"
						},
					]
				}
			`),
			Check: p.Check(map[string]any{
				"authorization.roles.#":                   2,
				"authorization.roles.0.name":              "App Developer",
				"authorization.roles.0.description":       "Builds apps and uploads new beta builds",
				"authorization.roles.0.permissions":       []string{"build-apps", "install-builds", "upload-builds"},
				"authorization.permissions.#":             3,
				"authorization.permissions.0.name":        "build-apps",
				"authorization.permissions.0.description": "Allowed to build and sign applications",
			}),
		},
	)
}
