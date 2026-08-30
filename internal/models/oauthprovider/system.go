package oauthprovider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var systemProviderNames = []string{"apple", "discord", "facebook", "github", "gitlab", "google", "linkedin", "microsoft", "slack"}

var systemClaimMapping = []string{"loginId", "username", "name", "email", "phoneNumber", "verifiedEmail", "verifiedPhone", "picture", "givenName", "middleName", "familyName"}

type systemCasingValidator struct {
}

func (v systemCasingValidator) Description(_ context.Context) string {
	return "system provider ids must use the lowercase identifier"
}

func (v systemCasingValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v systemCasingValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	value := req.ConfigValue.ValueString()
	if value == "" || slices.Contains(systemProviderNames, value) {
		return
	}
	for _, system := range systemProviderNames {
		if strings.EqualFold(system, value) {
			resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid Provider ID",
				fmt.Sprintf("Use %q instead of %q to configure the built-in system provider", system, value)))
			return
		}
	}
}
