package resources

import (
	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/descope/terraform-provider-descope/internal/models/customlanguage"
	"github.com/descope/terraform-provider-descope/internal/models/descoper"
	"github.com/descope/terraform-provider-descope/internal/models/engine"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/descope/terraform-provider-descope/internal/models/managementkey"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewAccessKeyResource() resource.Resource {
	return newResource[accesskey.AccessKeyModel]("access_key", accesskey.Schema)
}

func NewDescoperResource() resource.Resource {
	return newResource[descoper.DescoperModel]("descoper", descoper.Schema)
}

func NewManagementKeyResource() resource.Resource {
	return newResource[managementkey.ManagementKeyModel]("management_key", managementkey.Schema)
}

func NewInboundAppResource() resource.Resource {
	return newResource[inboundapp.InboundAppModel]("inbound_app", inboundapp.Schema)
}

func NewEngineResource() resource.Resource {
	return newResource[engine.EngineModel]("engine", engine.Schema)
}

func NewCustomLanguageResource() resource.Resource {
	// user-facing resource stays descope_custom_language (product vocabulary); the backend infra
	// entity is custom_locale (consistent with projectservice).
	return newResourceWithEntity[customlanguage.CustomLanguageModel]("custom_language", "custom_locale", customlanguage.Schema)
}
