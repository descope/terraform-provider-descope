package outboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Extra query parameters appended to the authorization or token request, matching the backend's repeated URLParam.

var URLParamAttributes = map[string]schema.Attribute{
	"key":   stringattr.Required(stringattr.StandardLenValidator),
	"value": stringattr.Default("", stringvalidator.LengthAtMost(4096)),
}

type URLParamModel struct {
	Key   stringattr.Type `tfsdk:"key"`
	Value stringattr.Type `tfsdk:"value"`
}

func (m *URLParamModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Key, data, "key")
	stringattr.Get(m.Value, data, "value")
	return data
}

func (m *URLParamModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Key, data, "key")
	stringattr.Set(&m.Value, data, "value")
}
