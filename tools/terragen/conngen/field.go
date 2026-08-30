package conngen

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

const (
	FieldTypeString       = "string"
	FieldTypeSecret       = "secret"
	FieldTypeBool         = "boolean"
	FieldTypeNumber       = "number"
	FieldTypeHTTPAuth     = "httpAuth"
	FieldTypeObject       = "object"
	FieldTypeAuditFilters = "auditFilters"
)

// Generated

var UseStaticIPsField = &Field{
	Name:        "useStaticIps",
	Description: "Whether the connector should send all requests from specific static IPs.",
	Type:        FieldTypeBool,
}

// The empty initial value doubles as the attribute default and pins the generated test values, so tests don't reference a real engine.
var EngineIDField = &Field{
	Name:        "engineId",
	Description: "The ID of the Descope Engine that runs this connector's actions inside your private network. Leave empty to run the connector in the Descope backend.",
	Type:        FieldTypeString,
	Initial:     "",
}

// Field
//
// Console-only field keys are knowingly ignored: displayName (attribute names come from name and naming.json), validation (a format
// hint whose validators would reject configs that apply today) and _initialValueComment. Dynamic is parsed but unused: it marks fields
// accepting {{...}} flow expressions, which any Terraform string already allows.
type Field struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Required    bool             `json:"required"`
	Dynamic     bool             `json:"dynamic"`
	Initial     any              `json:"initialValue"`
	Hidden      bool             `json:"hidden"`
	Options     []*FieldOption   `json:"options"`
	Dependency  *FieldDependency `json:"dependsOn"`

	naming *Naming

	// set on fields that other fields depend on, so test updates don't flip them
	hasDependents bool
}

func (f *Field) StructName() string {
	return f.naming.GetName("field", f.Name, "struct", f.defaultStructName())
}

func (f *Field) defaultStructName() string {
	return utils.CapitalCase(f.Name)
}

func (f *Field) StructType() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return `stringattr.Type`
	case FieldTypeBool:
		return `boolattr.Type`
	case FieldTypeNumber:
		return `floatattr.Type`
	case FieldTypeObject:
		return `strmapattr.Type`
	case FieldTypeAuditFilters:
		return `listattr.Type[AuditFilterFieldModel]`
	case FieldTypeHTTPAuth:
		return `objattr.Type[HTTPAuthFieldModel]`
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) AttributeName() string {
	return f.naming.GetName("field", f.Name, "attribute", f.defaultAttributeName())
}

func (f *Field) defaultAttributeName() string {
	return utils.SnakeCase(f.Name)
}

// the boilerplate attributes of standalone resources reserve several attribute names
var reservedResourceAttributes = []string{"id", "project_id", "name", "description"}

// Like AttributeName, but configuration fields colliding with a reserved boilerplate attribute name (e.g. a vendor-side `project_id`)
// get a `config_` prefix.
func (f *Field) ResourceAttributeName() string {
	if name := f.AttributeName(); !slices.Contains(reservedResourceAttributes, name) {
		return name
	}
	return "config_" + f.AttributeName()
}

// Like StructName, matching the `config_` prefix that ResourceAttributeName adds to colliding configuration fields.
func (f *Field) ResourceFieldName() string {
	if f.ResourceAttributeName() != f.AttributeName() {
		return "Config" + f.StructName()
	}
	return f.StructName()
}

func (f *Field) AttributeType() string {
	switch f.Type {
	case FieldTypeString:
		validator := ""

		if len(f.Options) > 0 {
			values := []string{}
			if !f.Required {
				values = append(values, `""`)
			}
			for _, option := range f.Options {
				values = append(values, fmt.Sprintf("%q", option.Value))
			}
			validator = fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(values, ", "))

			if v, ok := f.Initial.(string); ok {
				return fmt.Sprintf(`stringattr.Default(%q, %s)`, v, validator)
			}
			if f.Required {
				return fmt.Sprintf(`stringattr.Required(%s)`, validator)
			}
			return fmt.Sprintf(`stringattr.Default("", %s)`, validator)
		}

		if f.Required && f.Dependency == nil {
			return fmt.Sprintf(`stringattr.Required(%s)`, validator)
		}

		if validator != "" {
			validator = ", " + validator
		}

		defValue := ""
		if v, ok := f.Initial.(string); ok {
			defValue = v
		}

		return fmt.Sprintf(`stringattr.Default(%q, %s)`, defValue, validator)
	case FieldTypeSecret:
		if f.Required && f.Dependency == nil {
			return `stringattr.SecretRequired()`
		}
		return `stringattr.SecretOptional()`
	case FieldTypeBool:
		if f.Required && f.Dependency == nil {
			return `boolattr.Required()`
		}
		if f.Initial == true {
			return `boolattr.Default(true)`
		}
		return `boolattr.Default(false)`
	case FieldTypeNumber:
		if f.Required && f.Dependency == nil {
			return `floatattr.Required()`
		}
		if v, ok := f.Initial.(float64); ok {
			return fmt.Sprintf(`floatattr.Default(%g)`, v)
		}
		return `floatattr.Default(0)`
	case FieldTypeObject:
		return `strmapattr.Default()`
	case FieldTypeAuditFilters:
		return `listattr.Default[AuditFilterFieldModel](AuditFilterFieldAttributes)`
	case FieldTypeHTTPAuth:
		if f.Required && f.Dependency == nil {
			return `objattr.Required[HTTPAuthFieldModel](HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
		}
		return `objattr.Default(HTTPAuthFieldDefault, HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) GetValueStatement() string {
	if f.Hidden {
		switch f.Type {
		case FieldTypeString:
			return fmt.Sprintf(`c[%q] = %q`, f.Name, f.Initial.(string)) // nolint:forcetypeassert
		case FieldTypeBool:
			return fmt.Sprintf(`c[%q] = %t`, f.Name, f.Initial.(bool)) // nolint:forcetypeassert
		default:
			panic("unexpected hidden field type: " + f.Type)
		}
	}

	accessor := fmt.Sprintf(`m.%s`, f.ResourceFieldName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`stringattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`getHeaders(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`listattr.Get(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objattr.Get(%s, c, %q, h)`, accessor, f.Name)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) SetValueStatement() string {
	accessor := fmt.Sprintf(`&m.%s`, f.ResourceFieldName())
	switch f.Type {
	case FieldTypeString:
		return fmt.Sprintf(`stringattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeSecret:
		return fmt.Sprintf(`stringattr.Nil(%s)`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`setHeaders(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`listattr.Set(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objattr.Set(%s, c, %q, h)`, accessor, f.Name)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) IsZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.ResourceFieldName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`%s.ValueString() == ""`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`!%s.ValueBool()`, accessor)
	case FieldTypeNumber:
		return fmt.Sprintf(`%s.ValueFloat64() == 0`, accessor)
	case FieldTypeObject:
		return fmt.Sprintf(`%s.IsEmpty()`, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`%s.IsEmpty()`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`!%s.IsSet()`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

// HasDependencyChecks reports whether the dependency renders any check into Validate: a boolean dependency skips the conflict check
// when the field defaults to a non-zero value, leaving nothing to generate unless the field is also required.
func (f *Field) HasDependencyChecks() bool {
	if f.Dependency == nil {
		return false
	}
	if f.Dependency.Field.Type == FieldTypeBool {
		return f.Required || !f.HasNonZeroInitial()
	}
	return true
}

// HasNonZeroInitial reports whether the attribute defaults to a non-zero value, making an unset field indistinguishable from that value.
func (f *Field) HasNonZeroInitial() bool {
	switch v := f.Initial.(type) {
	case string:
		return v != ""
	case bool:
		return v
	case float64:
		return v != 0
	default:
		return false
	}
}

func (f *Field) IsNonZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.ResourceFieldName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`%s.ValueString() != ""`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`%s.ValueBool()`, accessor)
	case FieldTypeNumber:
		return fmt.Sprintf(`%s.ValueFloat64() != 0`, accessor)
	case FieldTypeObject:
		return fmt.Sprintf(`!%s.IsEmpty()`, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`!%s.IsEmpty()`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`%s.IsSet()`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

// Tests

func (f *Field) GetTestAssignment() string {
	return f.testAssignment(false)
}

// GetTestUpdateAssignment produces a changed value for the update step, keeping the create value where a change isn't valid: pinned
// initials, single-option fields, unsatisfied dependencies, and bools that other fields depend on.
func (f *Field) GetTestUpdateAssignment() string {
	return f.testAssignment(true)
}

func (f *Field) testAssignment(update bool) string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		if v, ok := f.Initial.(string); ok {
			return fmt.Sprintf(`%q`, v)
		}
		if !f.testDependencySatisfied() {
			return `null`
		}
		if len(f.Options) > 0 {
			return fmt.Sprintf(`%q`, f.testOption(update))
		}
		return fmt.Sprintf(`%q`, f.testString(update))
	case FieldTypeBool:
		return fmt.Sprintf(`%t`, f.testBool(update))
	case FieldTypeNumber:
		if !f.testDependencySatisfied() {
			return `null`
		}
		return fmt.Sprintf(`%d`, f.testNumber(update))
	case FieldTypeObject:
		return fmt.Sprintf(`{
    							"key" = %q
    						}`, f.testString(update))
	case FieldTypeAuditFilters:
		if !f.testDependencySatisfied() {
			return `[]`
		}
		return fmt.Sprintf(`[{ key = "actions", operator = "includes", values = [%q] }]`, f.testString(update))
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`{
    							bearer_token = %q
    						}`, f.testString(update))
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) GetTestCheck() string {
	return f.testCheck(false)
}

func (f *Field) GetTestUpdateCheck() string {
	return f.testCheck(true)
}

func (f *Field) testCheck(update bool) string {
	attribute := f.ResourceAttributeName()
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		if v, ok := f.Initial.(string); ok {
			return fmt.Sprintf(`"%s": %q`, attribute, v)
		}
		if !f.testDependencySatisfied() {
			if f.Type == FieldTypeSecret {
				return fmt.Sprintf(`"%s": testacc.AttributeIsNotSet`, attribute)
			}
			return fmt.Sprintf(`"%s": ""`, attribute)
		}
		if len(f.Options) > 0 {
			return fmt.Sprintf(`"%s": %q`, attribute, f.testOption(update))
		}
		return fmt.Sprintf(`"%s": %q`, attribute, f.testString(update))
	case FieldTypeBool:
		return fmt.Sprintf(`"%s": %t`, attribute, f.testBool(update))
	case FieldTypeNumber:
		if !f.testDependencySatisfied() {
			v, _ := f.Initial.(float64)
			return fmt.Sprintf(`"%s": %q`, attribute, fmt.Sprintf("%g", v))
		}
		return fmt.Sprintf(`"%s": %d`, attribute, f.testNumber(update))
	case FieldTypeObject:
		return fmt.Sprintf(`"%s.key": %q`, attribute, f.testString(update))
	case FieldTypeAuditFilters:
		if !f.testDependencySatisfied() {
			return fmt.Sprintf(`"%s.#": 0`, attribute)
		}
		return fmt.Sprintf(`"%s.0.values": []string{%q}`, attribute, f.testString(update))
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`"%s.bearer_token": %q`, attribute, f.testString(update))
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) testBool(update bool) bool {
	b := sha256.Sum256([]byte(f.Name))
	v := b[0]%2 == 0
	if update && !f.hasDependents {
		return !v
	}
	return v
}

func (f *Field) testDependencySatisfied() bool {
	d := f.Dependency
	if d == nil {
		return true
	}
	if d.Field.Type == FieldTypeBool {
		return d.Field.testBool(false) == d.Value
	}
	return d.Value == d.Field.Initial
}

func (f *Field) testString(update bool) string {
	name := f.Name
	if update {
		name += "+"
	}
	b := sha256.Sum256([]byte(name))
	s := base32.StdEncoding.EncodeToString(b[:])
	return strings.ToLower(s[:min(len(s), len(f.Name))])
}

func (f *Field) testNumber(update bool) int {
	if update {
		return len(f.Name) + 1
	}
	return len(f.Name)
}

// testOption returns the second option on update when there is one, since a single-option field has no other valid value.
func (f *Field) testOption(update bool) string {
	if update && len(f.Options) > 1 {
		return f.Options[1].Value
	}
	return f.Options[0].Value
}

// Dependency

type FieldDependency struct {
	Name   string   `json:"name"`
	Value  any      `json:"value"`
	Values []string `json:"values"`
	*Field
}

func (d *FieldDependency) DefaultValue() any {
	switch d.Field.Type {
	case FieldTypeString, FieldTypeSecret:
		v, _ := d.Field.Initial.(string)
		return v
	case FieldTypeBool:
		v, _ := d.Field.Initial.(bool)
		return v
	default:
		return d.Field.Initial
	}
}

func (d *FieldDependency) ValuesSlice() string {
	return fmt.Sprintf("%#v", d.Values)
}

// Options

type FieldOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
