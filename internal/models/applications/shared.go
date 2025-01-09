package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Handler = helpers.Handler

func sharedApplicationData(_ *Handler, id, name, description, logo types.String, disabled types.Bool) map[string]any {
	data := map[string]any{}
	stringattr.Get(id, data, "id")
	stringattr.Get(name, data, "name")
	stringattr.Get(description, data, "description")
	stringattr.Get(logo, data, "logo")
	data["enabled"] = !disabled.ValueBool()
	return data
}

// Match identifiers

func RequireID(h *Handler, data map[string]any, key string, name types.String, id *types.String) {
	n := name.ValueString()
	h.Log("Looking for %s application named '%s'", key, n)
	if appID, ok := requireID(h, data, key, n); ok {
		value := types.StringValue(appID)
		if !(*id).Equal(value) {
			h.Log("Setting new ID '%s' for %s application named '%s'", appID, key, n)
			*id = value
		} else {
			h.Log("Keeping existing ID '%s' for %s application named '%s'", appID, key, n)
		}
	}
}

func requireID(h *Handler, data map[string]any, key string, name string) (string, bool) {
	list, ok := data[key].([]any)
	if !ok {
		h.Error("Unexpected server response", "Expected to find list of '%s' applications to match with '%s' application", key, name)
		return "", false
	}

	for _, v := range list {
		if app, ok := v.(map[string]any); ok {
			if n, ok := app["name"].(string); ok && name == n {
				if id, ok := app["id"].(string); ok {
					return id, true
				}
			}
		}
	}

	h.Error("Application not found", "Expected to find application of type '%s' to match with '%s' application", key, name)
	return "", false
}
