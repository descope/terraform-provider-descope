package intattr_test

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/intattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	cases := []struct {
		name     string
		current  types.Int64
		data     map[string]any
		expected types.Int64
	}{
		{"value", types.Int64Null(), map[string]any{"key": float64(12)}, types.Int64Value(12)},
		{"int value", types.Int64Null(), map[string]any{"key": int64(12)}, types.Int64Value(12)},
		{"zero value", types.Int64Value(12), map[string]any{"key": float64(0)}, types.Int64Value(0)},
		{"missing key", types.Int64Null(), map[string]any{}, types.Int64Value(8)},
		{"missing key over value", types.Int64Value(12), map[string]any{}, types.Int64Value(12)},
		{"missing key over unknown", types.Int64Unknown(), map[string]any{}, types.Int64Value(8)},
		{"wrong type", types.Int64Null(), map[string]any{"key": "12"}, types.Int64Value(8)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			n := c.current
			intattr.SetDefault(&n, c.data, "key", 8)
			assert.Equal(t, c.expected, n)
		})
	}
}
