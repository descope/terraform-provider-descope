package export

import (
	"slices"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/read"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// The backend fills a built-in OAuth provider's redirect_url from an environment default the API doesn't expose, so the value
// most built-in providers share stands in for it.
func builtinRedirectURL(results []read.Result) string {
	counts := map[string]int{}
	for _, result := range results {
		if !result.Instance.Builtin {
			continue
		}
		if value, ok := result.Object.Attributes()["redirect_url"].(basetypes.StringValue); ok && value.ValueString() != "" {
			counts[value.ValueString()]++
		}
	}
	best, tied := "", false
	for value, count := range counts {
		switch {
		case count > counts[best]:
			best, tied = value, false
		case count == counts[best] && value != best:
			tied = true
		}
	}
	if tied || counts[best] < 2 {
		return ""
	}
	return best
}

func withoutBuiltinRedirectURL(attrs []prune.Attr, redirectURL string) []prune.Attr {
	return slices.DeleteFunc(slices.Clone(attrs), func(attr prune.Attr) bool {
		value, ok := attr.Value.(basetypes.StringValue)
		return ok && attr.Name == "redirect_url" && redirectURL != "" && value.ValueString() == redirectURL
	})
}

func uncustomized(attrs []prune.Attr) bool {
	return !slices.ContainsFunc(attrs, func(attr prune.Attr) bool { return attr.Name != "id" })
}
