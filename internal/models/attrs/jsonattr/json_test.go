package jsonattr

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestEquivalent(t *testing.T) {
	cases := []struct {
		name       string
		a, b       string
		equivalent bool
	}{
		{"identical", `{"a":1}`, `{"a":1}`, true},
		{"indented", `{"a":1,"b":[2,3]}`, "{\n  \"a\": 1,\n  \"b\": [\n    2,\n    3\n  ]\n}\n", true},
		{"key order", `{"a":1,"b":2}`, `{"b":2,"a":1}`, true},
		{"nested key order", `{"a":{"x":1,"y":2}}`, `{"a":{"y":2,"x":1}}`, true},
		{"arrays", `["alpha","beta"]`, "[\n  \"alpha\",\n  \"beta\"\n]", true},
		{"array order", `[1,2]`, `[2,1]`, false},
		{"different value", `{"a":1}`, `{"a":2}`, false},
		{"missing field", `{"a":1,"b":2}`, `{"a":1}`, false},
		{"large identifiers", `{"id":9007199254740993}`, `{"id":9007199254740992}`, false},
		{"integer and float", `{"a":1}`, `{"a":1.0}`, false},
		{"scalars", `"same"`, `"same"`, true},
		{"both invalid", `not json`, `not json`, true},
		{"one invalid", `{"a":1}`, `not json`, false},
		{"empty", ``, ``, true},
		{"one empty", `{}`, ``, false},
		{"trailing content", `{"a":1}`, `{"a":1} {"b":2}`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.equivalent, equivalent(c.a, c.b))
			assert.Equal(t, c.equivalent, equivalent(c.b, c.a))
		})
	}
}

func TestUseStateWhenEquivalent(t *testing.T) {
	indented := "{\n  \"a\": 1\n}"
	cases := []struct {
		name     string
		state    types.String
		plan     types.String
		expected types.String
	}{
		{"formatting only", types.StringValue(`{"a":1}`), types.StringValue(indented), types.StringValue(`{"a":1}`)},
		{"real change", types.StringValue(`{"a":1}`), types.StringValue(`{"a":2}`), types.StringValue(`{"a":2}`)},
		{"no state", types.StringNull(), types.StringValue(indented), types.StringValue(indented)},
		{"unknown state", types.StringUnknown(), types.StringValue(indented), types.StringValue(indented)},
		{"unknown plan", types.StringValue(`{"a":1}`), types.StringUnknown(), types.StringUnknown()},
		{"null plan", types.StringValue(`{"a":1}`), types.StringNull(), types.StringNull()},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := planmodifier.StringRequest{Path: path.Root("data"), StateValue: c.state, PlanValue: c.plan}
			resp := &planmodifier.StringResponse{PlanValue: c.plan}
			UseStateWhenEquivalent().PlanModifyString(context.Background(), req, resp)
			assert.Equal(t, c.expected, resp.PlanValue)
		})
	}
}

func TestGet(t *testing.T) {
	data := map[string]any{}
	Get(Value(`{"a":1}`), data, "obj")
	Get(Value(`["x","y"]`), data, "arr")
	assert.Equal(t, map[string]any{"a": json.Number("1")}, data["obj"])
	assert.Equal(t, []any{"x", "y"}, data["arr"])
	assert.Panics(t, func() { Get(Value("not json"), data, "bad") })
}

func TestGetRootKey(t *testing.T) {
	data := map[string]any{"existing": true}
	Get(Value(`{"a":1,"b":"x"}`), data, helpers.RootKey)
	assert.Equal(t, map[string]any{"existing": true, "a": json.Number("1"), "b": "x"}, data)
	assert.Panics(t, func() { Get(Value(`[1,2]`), map[string]any{}, helpers.RootKey) })
}

func TestSet(t *testing.T) {
	indented := "{\n  \"a\": 1\n}"
	cases := []struct {
		name     string
		existing string
		data     map[string]any
		options  []SetOption
		expected string
	}{
		{"overwrites", `{"a":1}`, map[string]any{"k": map[string]any{"a": 2}}, nil, `{"a":2}`},
		{"keeps equivalent formatting", indented, map[string]any{"k": map[string]any{"a": 1}}, nil, indented},
		{"skips when already set", `{"a":1}`, map[string]any{"k": map[string]any{"a": 2}}, []SetOption{SkipIfAlreadySet}, `{"a":1}`},
		{"fills when empty despite skip", "", map[string]any{"k": map[string]any{"a": 2}}, []SetOption{SkipIfAlreadySet}, `{"a":2}`},
		{"absent key on empty", "", map[string]any{}, nil, "{}"},
		{"absent key keeps existing", `{"a":1}`, map[string]any{}, nil, `{"a":1}`},
		{"null value treated as absent", "", map[string]any{"k": nil}, nil, "{}"},
		{"arrays", "", map[string]any{"k": []any{"x"}}, nil, `["x"]`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			value := Value(c.existing)
			Set(&value, c.data, "k", c.options...)
			assert.Equal(t, c.expected, value.ValueString())
		})
	}
}

func TestSetRootKey(t *testing.T) {
	indented := "{\n  \"a\": 1\n}"
	cases := []struct {
		name     string
		existing string
		data     map[string]any
		options  []SetOption
		expected string
	}{
		{"stores the whole map", "", map[string]any{"a": 1}, nil, `{"a":1}`},
		{"keeps equivalent formatting", indented, map[string]any{"a": 1}, nil, indented},
		{"overwrites on a real change", `{"a":1}`, map[string]any{"a": 2}, nil, `{"a":2}`},
		{"skips when already set", `{"a":1}`, map[string]any{"a": 2}, []SetOption{SkipIfAlreadySet}, `{"a":1}`},
		{"unserializable keeps existing", `{"a":1}`, map[string]any{"a": make(chan int)}, nil, `{"a":1}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			value := Value(c.existing)
			Set(&value, c.data, helpers.RootKey, c.options...)
			assert.Equal(t, c.expected, value.ValueString())
		})
	}
}
