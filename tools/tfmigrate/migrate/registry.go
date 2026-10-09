package migrate

import (
	"context"
	"fmt"
	"sort"
	"strings"

	descope "github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ProviderTypeName is the provider's Terraform type prefix, matching the descope provider's Metadata response.
const ProviderTypeName = "descope"

// ProviderSource is the only provider source the migration is allowed to initialize or verify.
const ProviderSource = "descope/descope"

// ProjectResourceType is the legacy resource whose ownership is being split.
const ProjectResourceType = ProviderTypeName + "_project"

// ScopeKind is the addressing scheme a destination resource uses, which also determines its import id grammar.
type ScopeKind string

const (
	// ScopeSingleton addresses a per-project singleton, imported with the project id alone.
	ScopeSingleton ScopeKind = "singleton"
	// ScopeProject addresses a project entity with a server assigned id, imported with project_id/entity_id.
	ScopeProject ScopeKind = "project"
	// ScopeApp addresses an entity owned by an SSO application, imported with project_id/app_id/entity_id.
	ScopeApp ScopeKind = "app"
	// ScopeMethod addresses a messaging template of an auth method, imported with project_id/method/template_id.
	ScopeMethod ScopeKind = "method"
	// ScopeCompany addresses a company level entity, imported with entity_id alone.
	ScopeCompany ScopeKind = "company"
)

// Destination is a resource the provider registers, described only by what a migration needs: its Terraform type, its
// schema, and how its import id is spelled. Everything here is read from the live provider so the migration can never
// drift from the resource implementations.
type Destination struct {
	Type      string
	Name      string
	Scope     ScopeKind
	Schema    schema.Schema
	WireType  string
	Connector bool
}

// Registry indexes every destination resource by Terraform type.
type Registry struct {
	byType map[string]*Destination
}

// NewRegistry discovers the destination resources by instantiating the provider's resource list and reading each
// resource's exported schema. The legacy project resource is deliberately absent: it is the migration source.
func NewRegistry(ctx context.Context) (*Registry, error) {
	r := &Registry{byType: map[string]*Destination{}}
	for _, constructor := range descope.NewDescopeProvider("tfmigrate")().Resources(ctx) {
		res := constructor()
		metadata := &resource.MetadataResponse{}
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: ProviderTypeName}, metadata)
		if metadata.TypeName == ProjectResourceType {
			continue
		}
		exportable, ok := res.(resources.ExportableResource)
		if !ok {
			return nil, fmt.Errorf("resource %s does not implement ExportableResource", metadata.TypeName)
		}
		sc := exportable.ExportSchema()
		dst := &Destination{
			Type:      metadata.TypeName,
			Name:      exportable.ExportName(),
			Scope:     scopeOf(sc, exportable.ExportSingleton()),
			Schema:    sc,
			WireType:  exportable.ExportWireType(ctx),
			Connector: strings.HasSuffix(exportable.ExportName(), "_connector"),
		}
		if existing, dup := r.byType[dst.Type]; dup {
			return nil, fmt.Errorf("resource %s registered twice (%s and %s)", dst.Type, existing.Name, dst.Name)
		}
		r.byType[dst.Type] = dst
	}
	if len(r.byType) == 0 {
		return nil, fmt.Errorf("no destination resources discovered")
	}
	return r, nil
}

// scopeOf mirrors the import id dispatch in internal/resources/base.go ImportState.
func scopeOf(sc schema.Schema, singleton bool) ScopeKind {
	if _, ok := sc.Attributes["project_id"]; !ok {
		return ScopeCompany
	}
	if singleton {
		return ScopeSingleton
	}
	if _, ok := sc.Attributes["method"]; ok {
		return ScopeMethod
	}
	if _, ok := sc.Attributes["app_id"]; ok {
		return ScopeApp
	}
	return ScopeProject
}

// Lookup returns the destination for a Terraform type.
func (r *Registry) Lookup(resourceType string) (*Destination, bool) {
	dst, ok := r.byType[resourceType]
	return dst, ok
}

// Types returns every discovered Terraform type, sorted.
func (r *Registry) Types() []string {
	types := make([]string, 0, len(r.byType))
	for t := range r.byType {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// ConnectorSuffix separates a connector resource name from its family name.
const ConnectorSuffix = "_connector"

// ConnectorFamilies maps each connector's legacy `connectors` attribute name to its destination Terraform type. The
// legacy attribute name is the destination resource name without the connector suffix, so the set is derived from the
// provider rather than restated here.
func (r *Registry) ConnectorFamilies() map[string]string {
	families := map[string]string{}
	for _, dst := range r.byType {
		if !dst.Connector {
			continue
		}
		families[strings.TrimSuffix(dst.Name, ConnectorSuffix)] = dst.Type
	}
	return families
}

// ImportID spells the import id for a destination instance, using the same grammar the resource's ImportState parses.
func (d *Destination) ImportID(scope Scope, entityID string) (string, error) {
	switch d.Scope {
	case ScopeSingleton:
		if scope.ProjectID == "" {
			return "", fmt.Errorf("%s: missing project id", d.Type)
		}
		return scope.ProjectID, nil
	case ScopeProject:
		if scope.ProjectID == "" || entityID == "" {
			return "", fmt.Errorf("%s: missing project id or entity id", d.Type)
		}
		return scope.ProjectID + "/" + entityID, nil
	case ScopeApp:
		if scope.ProjectID == "" || scope.AppID == "" || entityID == "" {
			return "", fmt.Errorf("%s: missing project id, app id or entity id", d.Type)
		}
		return scope.ProjectID + "/" + scope.AppID + "/" + entityID, nil
	case ScopeMethod:
		if scope.ProjectID == "" || scope.Method == "" || entityID == "" {
			return "", fmt.Errorf("%s: missing project id, method or template id", d.Type)
		}
		return scope.ProjectID + "/" + scope.Method + "/" + entityID, nil
	case ScopeCompany:
		if entityID == "" {
			return "", fmt.Errorf("%s: missing entity id", d.Type)
		}
		return entityID, nil
	}
	return "", fmt.Errorf("%s: unknown scope %q", d.Type, d.Scope)
}
