package outboundscim

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const maskedValue = "MASKED_VALUE"

func (m OutboundSCIMConfigurationModel) ConfigurationForWrite() (map[string]any, error) {
	configuration, err := DecodeConfiguration(m.Configuration.ValueString())
	if err != nil {
		return nil, fmt.Errorf("configuration: %w", err)
	}
	if path := secretFieldPath(configuration); path != "" {
		return nil, fmt.Errorf("configuration contains secret field %q; move it to secrets", path)
	}
	if m.Secrets.IsNull() || m.Secrets.IsUnknown() || m.Secrets.ValueString() == "" {
		return configuration, nil
	}
	secrets, err := DecodeConfiguration(m.Secrets.ValueString())
	if err != nil {
		return nil, fmt.Errorf("secrets: %w", err)
	}
	return mergeMaps(configuration, secrets), nil
}

func DecodeConfiguration(raw string) (map[string]any, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var configuration map[string]any
	if err := decoder.Decode(&configuration); err != nil {
		return nil, fmt.Errorf("must be a JSON object: %w", err)
	}
	if configuration == nil {
		return nil, fmt.Errorf("must be a JSON object")
	}
	if decoder.More() {
		return nil, fmt.Errorf("must contain one JSON object")
	}
	return configuration, nil
}

func NormalizeConfigurationForState(configuration map[string]any) (string, error) {
	clean := cloneMap(configuration)
	delete(clean, "federatedAppId")
	removeMaskedValues(clean)
	encoded, err := json.Marshal(clean)
	if err != nil {
		return "", fmt.Errorf("encode configuration: %w", err)
	}
	return string(encoded), nil
}

type jsonObjectValidator struct {
	rejectSecrets bool
}

func (v jsonObjectValidator) Description(context.Context) string {
	return "must be a JSON object"
}

func (v jsonObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonObjectValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	configuration, err := DecodeConfiguration(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Invalid JSON object", err.Error()))
		return
	}
	if path := secretFieldPath(configuration); v.rejectSecrets && path != "" {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(req.Path, "Secret field in stateful configuration", fmt.Sprintf("Move %q to the write-only secrets attribute.", path)))
	}
}

func mergeMaps(base, overlay map[string]any) map[string]any {
	merged := cloneMap(base)
	for key, value := range overlay {
		if overlayMap, ok := value.(map[string]any); ok {
			if baseMap, ok := merged[key].(map[string]any); ok {
				merged[key] = mergeMaps(baseMap, overlayMap)
				continue
			}
		}
		merged[key] = value
	}
	return merged
}

func cloneMap(source map[string]any) map[string]any {
	clone := make(map[string]any, len(source))
	for key, value := range source {
		clone[key] = cloneValue(value)
	}
	return clone
}

func cloneValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneMap(typed)
	case []any:
		clone := make([]any, len(typed))
		for index, item := range typed {
			clone[index] = cloneValue(item)
		}
		return clone
	default:
		return typed
	}
}

func removeMaskedValues(configuration map[string]any) {
	for key, value := range configuration {
		switch typed := value.(type) {
		case string:
			if typed == maskedValue || typed == "REMOVE_MASKED_VALUE" {
				delete(configuration, key)
			}
		case map[string]any:
			removeMaskedValues(typed)
		case []any:
			for _, item := range typed {
				if nested, ok := item.(map[string]any); ok {
					removeMaskedValues(nested)
				}
			}
		}
	}
}

func secretFieldPath(configuration map[string]any) string {
	for key, value := range configuration {
		switch strings.ToLower(key) {
		case "bearertoken", "clientsecret", "password", "privatekey", "secret", "token":
			return key
		}
		switch typed := value.(type) {
		case map[string]any:
			if path := secretFieldPath(typed); path != "" {
				return key + "." + path
			}
		case []any:
			for index, item := range typed {
				if nested, ok := item.(map[string]any); ok {
					if path := secretFieldPath(nested); path != "" {
						return fmt.Sprintf("%s[%d].%s", key, index, path)
					}
				}
			}
		}
	}
	return ""
}
