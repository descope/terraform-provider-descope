package list

import (
	"context"
	"fmt"
	"net"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single list of allowed or denied values in a Descope project. A list holds either plain text values, IP addresses and CIDR ranges, or an arbitrary JSON object, and can be referenced from flows.",
	Attributes:          ListAttributes,
}

var ListAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"project_id":  stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":        stringattr.Required(stringvalidator.LengthAtMost(100)),
	"description": stringattr.Default("", stringattr.StandardLenValidator),
	"texts":       strsetattr.Default(stringattr.NonEmptyValidator),
	"ips":         strsetattr.Default(ipValidator),
	"json":        jsonattr.Default(""),
}

type ListModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	ProjectID   stringattr.Type `tfsdk:"project_id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Texts       strsetattr.Type `tfsdk:"texts"`
	IPs         strsetattr.Type `tfsdk:"ips"`
	JSON        jsonattr.Type   `tfsdk:"json"`
}

func (m *ListModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")

	switch {
	case !m.Texts.IsEmpty():
		data["type"] = "texts"
		strsetattr.Get(m.Texts, data, "data", h)
	case !m.IPs.IsEmpty():
		data["type"] = "ips"
		strsetattr.Get(m.IPs, data, "data", h)
	case m.JSON.ValueString() != "":
		data["type"] = "json"
		jsonattr.Get(m.JSON, data, "data")
	}

	return data
}

func (m *ListModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")

	switch listType, _ := data["type"].(string); listType {
	case "texts":
		strsetattr.Set(&m.Texts, data, "data", h) // the server dedups these, which a set absorbs
		m.IPs, m.JSON = strsetattr.Empty(), jsonattr.Value("")
	case "ips":
		strsetattr.Set(&m.IPs, data, "data", h)
		m.Texts, m.JSON = strsetattr.Empty(), jsonattr.Value("")
	case "json":
		jsonattr.Set(&m.JSON, data, "data")
		m.Texts, m.IPs = strsetattr.Empty(), strsetattr.Empty()
	}
}

func (m *ListModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Texts, m.IPs, m.JSON) {
		return
	}
	chosen := 0
	for _, chose := range []bool{!m.Texts.IsEmpty(), !m.IPs.IsEmpty(), m.JSON.ValueString() != ""} {
		if chose {
			chosen++
		}
	}
	if chosen != 1 {
		h.Invalid("Exactly one of the 'texts', 'ips' or 'json' attributes must be set")
	}
}

// IP

var ipValidator = ipAddressValidator{}

type ipAddressValidator struct{}

func (v ipAddressValidator) Description(_ context.Context) string {
	return "must be an IP address or CIDR range"
}

func (v ipAddressValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ipAddressValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if _, _, err := net.ParseCIDR(value); err == nil {
		return
	}
	if net.ParseIP(value) == nil {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid Attribute Value", fmt.Sprintf("Attribute %s must be an IP address or CIDR range", req.Path)))
	}
}

func (m *ListModel) GetID() stringattr.Type {
	return m.ID
}

func (m *ListModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *ListModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}
