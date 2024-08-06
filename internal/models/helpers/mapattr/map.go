package mapattr

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Optional(attributes map[string]schema.Attribute, validators ...validator.Object) schema.MapNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.MapNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		Default:      mapdefault.StaticValue(types.MapNull(nested.Type())),
	}
}

func StringOptional(validators ...validator.Map) schema.MapAttribute {
	return optionalTypeMap(types.StringType, validators...)
}

func optionalTypeMap(elementType attr.Type, validators ...validator.Map) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: elementType,
		Default:     mapdefault.StaticValue(types.MapNull(elementType)),
		Validators:  validators,
	}
}
