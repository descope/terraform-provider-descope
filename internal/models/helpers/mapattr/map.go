package mapattr

import (
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Optional(attributes map[string]schema.Attribute, extras ...any) schema.MapNestedAttribute {
	mapValidators, objectValidators := parseExtras(extras)
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: objectValidators,
	}
	return schema.MapNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		Default:      mapdefault.StaticValue(types.MapNull(nested.Type())),
		Validators:   mapValidators,
	}
}

func StringOptional(validators ...validator.Map) schema.MapAttribute {
	return optionalTypeMap(types.StringType, validators)
}

func optionalTypeMap(elementType attr.Type, validators []validator.Map) schema.MapAttribute {
	return schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: elementType,
		Default:     mapdefault.StaticValue(types.MapNull(elementType)),
		Validators:  validators,
	}
}

func Get[T any, M helpers.Model[T]](m M, data map[string]any, key string, h *helpers.Handler) {
	if m != nil {
		data[key] = m.Values(h)
	}
}

func Set[T any, M helpers.Model[T]](m *M, data map[string]any, key string, h *helpers.Handler) { // using *M to stay consistent with other Set functions
	if v, ok := data[key].(map[string]any); ok {
		if *m != nil {
			(*m).SetValues(h, v)
		}
	}
}

func parseExtras(extras []any) (mapValidators []validator.Map, objectValidators []validator.Object) {
	for _, e := range extras {
		matched := false
		if v, ok := e.(validator.Map); ok {
			matched = true
			mapValidators = append(mapValidators, v)
		}
		if v, ok := e.(validator.Object); ok {
			matched = true
			objectValidators = append(objectValidators, v)
		}
		if !matched {
			panic(fmt.Sprintf("Unexpected extra value of type %T in map attribute", e))
		}
	}
	return
}
