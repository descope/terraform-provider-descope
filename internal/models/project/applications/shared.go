package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
)

func sharedApplicationData(_ *helpers.Handler, id, name, description, logo stringattr.Type, disabled boolattr.Type) map[string]any {
	data := map[string]any{}
	stringattr.Get(id, data, "id")
	stringattr.Get(name, data, "name")
	stringattr.Get(description, data, "description")
	stringattr.Get(logo, data, "logo")
	boolattr.GetNot(disabled, data, "enabled")
	return data
}

func setSharedApplicationData(_ *helpers.Handler, data map[string]any, id, name, description, logo *stringattr.Type, disabled *boolattr.Type) {
	stringattr.Set(id, data, "id")
	stringattr.Set(name, data, "name")
	stringattr.Set(description, data, "description")
	stringattr.Set(logo, data, "logo")
	boolattr.SetNot(disabled, data, "enabled")
}
