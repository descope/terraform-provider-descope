package boolattr_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	cases := []struct {
		name     string
		current  types.Bool
		data     map[string]any
		expected types.Bool
	}{
		{"value", types.BoolNull(), map[string]any{"key": true}, types.BoolValue(true)},
		{"false value", types.BoolNull(), map[string]any{"key": false}, types.BoolValue(false)},
		{"missing key", types.BoolNull(), map[string]any{}, types.BoolNull()},
		{"missing key over value", types.BoolValue(true), map[string]any{}, types.BoolValue(true)},
		{"missing key over unknown", types.BoolUnknown(), map[string]any{}, types.BoolValue(false)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := c.current
			boolattr.Set(&b, c.data, "key")
			assert.Equal(t, c.expected, b)
		})
	}
}

func TestSetDefault(t *testing.T) {
	cases := []struct {
		name     string
		current  types.Bool
		data     map[string]any
		expected types.Bool
	}{
		{"value", types.BoolNull(), map[string]any{"key": false}, types.BoolValue(false)},
		{"missing key", types.BoolNull(), map[string]any{}, types.BoolValue(true)},
		{"missing key over value", types.BoolValue(false), map[string]any{}, types.BoolValue(false)},
		{"missing key over unknown", types.BoolUnknown(), map[string]any{}, types.BoolValue(true)},
		{"wrong type", types.BoolNull(), map[string]any{"key": "yes"}, types.BoolValue(true)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b := c.current
			boolattr.SetDefault(&b, c.data, "key", true)
			assert.Equal(t, c.expected, b)
		})
	}
}
