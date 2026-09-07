package stringattr_test

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestURLValidator(t *testing.T) {
	cases := []struct {
		name  string
		value types.String
		valid bool
	}{
		{"empty", types.StringValue(""), true},
		{"null", types.StringNull(), true},
		{"unknown", types.StringUnknown(), true},
		{"https", types.StringValue("https://example.com"), true},
		{"path and query", types.StringValue("http://localhost:8443/p?q=1"), true},
		{"garbage", types.StringValue("not a url"), false},
		{"scheme only", types.StringValue("https://"), false},
		{"no scheme", types.StringValue("example.com"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := validator.StringRequest{Path: path.Root("url"), ConfigValue: c.value}
			resp := &validator.StringResponse{}
			stringattr.URLValidator.ValidateString(context.Background(), req, resp)
			assert.Equal(t, c.valid, !resp.Diagnostics.HasError())
		})
	}
}
