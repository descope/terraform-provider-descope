package stringattr_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	cases := []struct {
		name     string
		current  types.String
		data     map[string]any
		expected types.String
	}{
		{"value", types.StringNull(), map[string]any{"key": "a"}, types.StringValue("a")},
		{"empty value", types.StringNull(), map[string]any{"key": ""}, types.StringValue("")},
		{"missing key", types.StringNull(), map[string]any{}, types.StringNull()},
		{"missing key over value", types.StringValue("a"), map[string]any{}, types.StringValue("a")},
		{"missing key over unknown", types.StringUnknown(), map[string]any{}, types.StringValue("")},
		{"wrong type", types.StringNull(), map[string]any{"key": 42}, types.StringNull()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.current
			stringattr.Set(&s, c.data, "key")
			assert.Equal(t, c.expected, s)
		})
	}
}

func TestSetDefault(t *testing.T) {
	cases := []struct {
		name     string
		current  types.String
		data     map[string]any
		expected types.String
	}{
		{"value", types.StringNull(), map[string]any{"key": "a"}, types.StringValue("a")},
		{"empty value", types.StringNull(), map[string]any{"key": ""}, types.StringValue("d")},
		{"missing key", types.StringNull(), map[string]any{}, types.StringValue("d")},
		{"missing key over value", types.StringValue("a"), map[string]any{}, types.StringValue("a")},
		{"missing key over unknown", types.StringUnknown(), map[string]any{}, types.StringValue("d")},
		{"wrong type", types.StringNull(), map[string]any{"key": 42}, types.StringValue("d")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := c.current
			stringattr.SetDefault(&s, c.data, "key", "d")
			assert.Equal(t, c.expected, s)
		})
	}
}

func TestSetDefaultEmpty(t *testing.T) {
	s := types.StringNull()
	stringattr.SetDefault(&s, map[string]any{}, "key", "")
	assert.Equal(t, types.StringValue(""), s)
}
