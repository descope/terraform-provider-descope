package attribute

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var InboundAppAttributeSchema = schema.Schema{
	MarkdownDescription: "Manages a single inbound app custom attribute definition in a Descope project.",
	Attributes:          InboundAppAttributeAttributes,
}

var InboundAppAttributeAttributes = map[string]schema.Attribute{
	"id":             stringattr.Required(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20), stringplanmodifier.RequiresReplace()),
	"project_id":     stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":           stringattr.Required(stringattr.StandardLenValidator),
	"type":           stringattr.Required(typeValidator, stringplanmodifier.RequiresReplace()),
	"select_options": strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1))),
}

type InboundAppAttributeModel struct {
	AttributeModel
}
