package registry

import (
	"context"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
)

func Load(ctx context.Context) map[string]resources.ExportableResource {
	result := map[string]resources.ExportableResource{}
	for _, constructor := range provider.NewDescopeProvider("tfexport")().Resources(ctx) {
		if exportable, ok := constructor().(resources.ExportableResource); ok {
			result["descope_"+exportable.ExportName()] = exportable
		}
	}
	return result
}

func ConnectorTypes(ctx context.Context, loaded map[string]resources.ExportableResource) map[string]string {
	result := map[string]string{}
	for name, exportable := range loaded {
		if !strings.HasSuffix(name, "_connector") {
			continue
		}
		if wireType := exportable.ExportWireType(ctx); wireType != "" {
			result[wireType] = name
		}
	}
	return result
}
