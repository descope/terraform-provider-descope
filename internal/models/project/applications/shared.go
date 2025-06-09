package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

type Handler = helpers.Handler

func sharedApplicationData(_ *Handler, id, name, description, logo stringattr.Type, disabled boolattr.Type) map[string]any {
	data := map[string]any{}
	stringattr.Get(id, data, "id")
	stringattr.Get(name, data, "name")
	stringattr.Get(description, data, "description")
	stringattr.Get(logo, data, "logo")
	boolattr.GetNot(disabled, data, "enabled")
	return data
}

func setSharedApplicationData(_ *Handler, data map[string]any, id, name, description, logo *stringattr.Type, disabled *boolattr.Type) {
	stringattr.Set(id, data, "id")
	stringattr.Set(name, data, "name")
	stringattr.Set(description, data, "description")
	stringattr.Set(logo, data, "logo")
	boolattr.SetNot(disabled, data, "enabled")
}

// Match identifiers

func RequireID(h *Handler, data map[string]any, key string, name stringattr.Type, id *stringattr.Type) {
	n := name.ValueString()
	h.Log("Looking for %s application named '%s'", key, n)
	if appID, ok := requireID(h, data, key, n); ok {
		value := stringattr.Value(appID)
		if !id.Equal(value) {
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
