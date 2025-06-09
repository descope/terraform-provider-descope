package templates

import (
	"strings"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

func requireTemplateID(h *helpers.Handler, data map[string]any, typ string, name string) (string, bool) {
	list, ok := data[typ].([]any)
	if !ok {
		h.Error("Unexpected server response", "Expected to find list of templates in '%s' to match with '%s' template", typ, name)
		return "", false
	}

	for _, v := range list {
		if template, ok := v.(map[string]any); ok {
			if n, ok := template["name"].(string); ok && name == n {
				if id, ok := template["id"].(string); ok {
					return id, true
				}
			}
		}
	}

	h.Log("Expected to find template in '%s' to match with '%s' template, this could be due to a state file mismatch", typ, name)
	return "", false
}

func replaceConnectorIDWithReference(s *stringattr.Type, h *helpers.Handler) {
	if connector := strings.Split(s.ValueString(), ":"); len(connector) == 2 {
		ref := h.Refs.Name(connector[1])
		if ref != "" {
			*s = stringattr.Value(ref)
		}
	}
}
