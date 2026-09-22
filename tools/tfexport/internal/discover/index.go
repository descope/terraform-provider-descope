package discover

import (
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
)

var templateKinds = []struct {
	key      string
	resource string
	methods  []string
}{
	{"emailTemplates", "descope_email_template", []string{"enchantedlink", "magiclink", "otp", "password", "sso"}},
	{"textTemplates", "descope_text_template", []string{"magiclink", "otp"}},
	{"voiceTemplates", "descope_voice_template", []string{"otp"}},
}

type Instance struct {
	Resource    string // the Terraform resource type, e.g. "descope_flow"
	ID          string // the entity ID used to read the resource
	Scope       string // the method or app_id for scoped entities, empty otherwise
	Name        string // a friendly name for deriving the resource label
	MayBeAbsent bool   // the snapshot lists this entity before it exists, so a not-found read is expected
}

func Instances(files map[string]any, connectorTypes map[string]string) (instances []Instance, warnings []warn.Warning) {
	lost := func(format string, args ...any) {
		warnings = append(warnings, warn.Lost(format, args...))
	}
	note := func(format string, args ...any) {
		warnings = append(warnings, warn.Note(format, args...))
	}

	for _, id := range stringList(files, "flows/flows.json", "flows") {
		instances = append(instances, Instance{Resource: "descope_flow", ID: id, Name: id})
	}

	for _, id := range stringList(files, "widgets/widgets.json", "widgets") {
		// every project lists the built-in widgets, but one only exists as an entity once it has been customized
		instances = append(instances, Instance{Resource: "descope_widget", ID: id, Name: id, MayBeAbsent: true})
	}

	if flows := stringList(files, "widgets/flows/flows.json", "flows"); len(flows) > 0 {
		note("Skipped %d widget flows which are not supported yet", len(flows))
	}

	for _, slug := range stringList(files, "connectors/connectors.json", "connectors") {
		connector, ok := files["connectors/"+slug+".json"].(map[string]any)
		if !ok {
			lost("Skipped connector %s: no snapshot file", slug)
			continue
		}
		wireType, _ := connector["type"].(string)
		id, _ := connector["id"].(string)
		name, _ := connector["name"].(string)
		resource, ok := connectorTypes[wireType]
		if !ok {
			note("Skipped connector %s: unsupported connector type %q", name, wireType)
			continue
		}
		if id == "" {
			lost("Skipped connector %s: no id in snapshot", name)
			continue
		}
		instances = append(instances, Instance{Resource: resource, ID: id, Name: name})
	}

	for _, app := range objectList(files, "applications.json", "applications") {
		id, _ := app["id"].(string)
		name, _ := app["name"].(string)
		var resource string
		for key, value := range map[string]string{"oidc": "descope_oidc_app", "saml": "descope_saml_app", "wsfed": "descope_wsfed_app"} {
			if _, ok := app[key]; ok {
				resource = value
				break
			}
		}
		if resource == "" {
			lost("Skipped application %s: unrecognized application type", name)
			continue
		}
		instances = append(instances, Instance{Resource: resource, ID: id, Name: name})

		// application roles and permissions are nested in the application entries
		for key, scoped := range map[string]string{"roles": "descope_app_role", "permissions": "descope_app_permission"} {
			for _, entry := range anyObjectList(app[key]) {
				entryID, _ := entry["id"].(string)
				entryName, _ := entry["name"].(string)
				if entryID != "" {
					instances = append(instances, Instance{Resource: scoped, ID: entryID, Scope: id, Name: entryName})
				}
			}
		}
	}

	for _, app := range objectList(files, "outboundapps.json", "outboundApps") {
		id, _ := app["id"].(string)
		name, _ := app["name"].(string)
		instances = append(instances, Instance{Resource: "descope_outbound_app", ID: id, Name: name})
	}

	for _, template := range objectList(files, "jwttemplates.json", "templates") {
		id, _ := template["id"].(string)
		name, _ := template["name"].(string)
		instances = append(instances, Instance{Resource: "descope_jwt_template", ID: id, Name: name})
	}

	for _, list := range objectList(files, "lists.json", "lists") {
		id, _ := list["id"].(string)
		name, _ := list["name"].(string)
		instances = append(instances, Instance{Resource: "descope_list", ID: id, Name: name})
	}

	for _, role := range objectList(files, "roles.json", "roles") {
		instances = append(instances, roleOrPermission(role, "descope_role", "descope_app_role"))
	}

	for _, permission := range objectList(files, "permissions.json", "permissions") {
		instances = append(instances, roleOrPermission(permission, "descope_permission", "descope_app_permission"))
	}

	if attributes, ok := files["attributes.json"].(map[string]any); ok {
		for key, resource := range map[string]string{"user": "descope_user_attribute", "tenant": "descope_tenant_attribute", "accessKey": "descope_access_key_attribute"} {
			for _, attribute := range anyObjectList(attributes[key]) {
				name, _ := attribute["name"].(string)
				display, _ := attribute["displayName"].(string)
				instances = append(instances, Instance{Resource: resource, ID: name, Name: display})
			}
		}
	}

	for path, value := range files {
		method, ok := strings.CutPrefix(path, "auth/")
		if !ok {
			continue
		}
		method = strings.TrimSuffix(method, ".json")
		settings, ok := value.(map[string]any)
		if !ok {
			continue
		}
		for _, kind := range templateKinds {
			for _, template := range anyObjectList(settings[kind.key]) {
				id, _ := template["id"].(string)
				name, _ := template["name"].(string)
				if id == "" || id == "System" {
					continue // the built-in System template is not an exportable entity and cannot be read
				}
				if !slices.Contains(kind.methods, method) {
					note("Skipped %s %s: the %s method is not supported by the resource", kind.resource, name, method)
					continue
				}
				instances = append(instances, Instance{Resource: kind.resource, ID: id, Scope: method, Name: name})
			}
		}
		if providers, ok := settings["providerSettings"].(map[string]any); ok {
			// the provider identity is the map key: custom entries also carry it as their name field
			for key, provider := range providers {
				settings, ok := provider.(map[string]any)
				if !ok {
					continue
				}
				if custom, _ := settings["custom"].(bool); custom {
					instances = append(instances, Instance{Resource: "descope_oauth_provider", ID: key, Name: key})
				}
			}
		}
	}

	return dedupe(instances), warnings
}

// the same entity can be listed twice by the snapshot, and two resources importing one entity is invisible to the integrity checks
func dedupe(instances []Instance) []Instance {
	seen := map[Instance]bool{}
	result := make([]Instance, 0, len(instances))
	for _, instance := range instances {
		if instance.ID == "" {
			result = append(result, instance) // no id is no identity to compare on, so it is left for the read to report
			continue
		}
		key := Instance{Resource: instance.Resource, ID: instance.ID, Scope: instance.Scope}
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, instance)
	}
	return result
}

func roleOrPermission(entry map[string]any, resource, appResource string) Instance {
	id, _ := entry["id"].(string)
	name, _ := entry["name"].(string)
	if appID, _ := entry["appId"].(string); appID != "" {
		return Instance{Resource: appResource, ID: id, Scope: appID, Name: name}
	}
	return Instance{Resource: resource, ID: id, Name: name}
}

func stringList(files map[string]any, path, key string) []string {
	object, _ := files[path].(map[string]any)
	values, _ := object[key].([]any)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if s, ok := value.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func objectList(files map[string]any, path, key string) []map[string]any {
	object, _ := files[path].(map[string]any)
	return anyObjectList(object[key])
}

func anyObjectList(value any) []map[string]any {
	values, _ := value.([]any)
	result := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if object, ok := value.(map[string]any); ok {
			result = append(result, object)
		}
	}
	return result
}
