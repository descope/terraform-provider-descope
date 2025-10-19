package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var systemProviderNames = []string{"apple", "discord", "facebook", "github", "gitlab", "google", "linkedin", "microsoft", "slack"}

var OAuthSystemProviderAttributes = map[string]schema.Attribute{
	"apple":     objattr.Default[OAuthProviderAppleModel](nil, OAuthProviderAppleAttributes, OAuthProviderAppleValidator),
	"discord":   objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"facebook":  objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"github":    objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"gitlab":    objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"google":    objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"linkedin":  objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"microsoft": objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
	"slack":     objattr.Default[OAuthProviderSystemModel](nil, OAuthProviderSystemAttributes, OAuthProviderSystemValidator),
}

type OAuthSystemProvidersModel struct {
	Apple     objattr.Type[OAuthProviderAppleModel]  `tfsdk:"apple"`
	Discord   objattr.Type[OAuthProviderSystemModel] `tfsdk:"discord"`
	Facebook  objattr.Type[OAuthProviderSystemModel] `tfsdk:"facebook"`
	Github    objattr.Type[OAuthProviderSystemModel] `tfsdk:"github"`
	Gitlab    objattr.Type[OAuthProviderSystemModel] `tfsdk:"gitlab"`
	Google    objattr.Type[OAuthProviderSystemModel] `tfsdk:"google"`
	Linkedin  objattr.Type[OAuthProviderSystemModel] `tfsdk:"linkedin"`
	Microsoft objattr.Type[OAuthProviderSystemModel] `tfsdk:"microsoft"`
	Slack     objattr.Type[OAuthProviderSystemModel] `tfsdk:"slack"`
}

func (m *OAuthSystemProvidersModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	// XXX might need to drop if null
	objattr.Get(m.Apple, data, "apple", h)
	objattr.Get(m.Discord, data, "discord", h)
	objattr.Get(m.Facebook, data, "facebook", h)
	objattr.Get(m.Github, data, "github", h)
	objattr.Get(m.Gitlab, data, "gitlab", h)
	objattr.Get(m.Google, data, "google", h)
	objattr.Get(m.Linkedin, data, "linkedin", h)
	objattr.Get(m.Microsoft, data, "microsoft", h)
	objattr.Get(m.Slack, data, "slack", h)
	return data
}

func (m *OAuthSystemProvidersModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.Apple, data, "apple", h)
	objattr.Set(&m.Discord, data, "discord", h)
	objattr.Set(&m.Facebook, data, "facebook", h)
	objattr.Set(&m.Github, data, "github", h)
	objattr.Set(&m.Gitlab, data, "gitlab", h)
	objattr.Set(&m.Google, data, "google", h)
	objattr.Set(&m.Linkedin, data, "linkedin", h)
	objattr.Set(&m.Microsoft, data, "microsoft", h)
	objattr.Set(&m.Slack, data, "slack", h)
}
