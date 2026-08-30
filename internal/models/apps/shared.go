package apps

import (
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var appIDValidator = stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9\-_]{1,30}$`), "must be 1 to 30 alphanumeric, hyphen or underscore characters")

func sharedAppValues(_ *helpers.Handler, id, name, description, logo stringattr.Type, disabled boolattr.Type) map[string]any {
	data := map[string]any{}
	stringattr.Get(id, data, "id")
	stringattr.Get(name, data, "name")
	stringattr.Get(description, data, "description")
	stringattr.Get(logo, data, "logo")
	boolattr.GetNot(disabled, data, "enabled")
	return data
}

func setSharedAppValues(_ *helpers.Handler, data map[string]any, id, name, description, logo *stringattr.Type, disabled *boolattr.Type) {
	stringattr.Set(id, data, "id")
	stringattr.Set(name, data, "name")
	stringattr.Set(description, data, "description")
	stringattr.Set(logo, data, "logo")
	boolattr.SetNot(disabled, data, "enabled")
}

// Attribute Mapping

var AttributeMappingAttributes = map[string]schema.Attribute{
	"name":  stringattr.Required(),
	"value": stringattr.Required(),
}

type AttributeMappingModel struct {
	Name  stringattr.Type `tfsdk:"name"`
	Value stringattr.Type `tfsdk:"value"`
}

func (m *AttributeMappingModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Value, data, "value")
	return data
}

func (m *AttributeMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Value, data, "value")
}
