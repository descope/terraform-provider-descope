package listattr

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Required(attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Required:     true,
		NestedObject: nested,
	}
}

func Optional(attributes map[string]schema.Attribute, validators ...validator.Object) schema.ListNestedAttribute {
	nested := schema.NestedAttributeObject{
		Attributes: attributes,
		Validators: validators,
	}
	return schema.ListNestedAttribute{
		Optional:     true,
		Computed:     true,
		NestedObject: nested,
		Default:      listdefault.StaticValue(types.ListNull(nested.Type())),
	}
}

func Get[T any, M helpers.Model[T]](list []M, data map[string]any, key string, h *helpers.Handler) {
	data[key] = valuesFromModels(h, list)
}

func valuesFromModels[T any, M helpers.Model[T]](h *helpers.Handler, list []M) []any {
	var values []any
	for _, v := range list {
		values = append(values, v.Values(h))
	}
	return values
}

func Set[T any, M helpers.Model[T]](list *[]M, data map[string]any, key string, h *helpers.Handler) bool {
	changed := false
	if vs, ok := data[key].([]any); ok {
		*list = []M{} // start from scratch to avoid unwanted duplications / collisions
		for _, v := range vs {
			if m, ok := v.(map[string]any); ok {
				var model M
				model = utils.ZVL(model)
				model.SetValues(h, m)
				*list = append(*list, model)
				changed = true
			}
		}
	}
	return changed
}
