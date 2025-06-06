package schema

type FieldType string

const (
	FieldTypeBool     FieldType = "bool"
	FieldTypeDuration FieldType = "duration"
	FieldTypeFloat    FieldType = "float"
	FieldTypeInt      FieldType = "int"
	FieldTypeList     FieldType = "list"
	FieldTypeObject   FieldType = "object"
	FieldTypeMap      FieldType = "map"
	FieldTypeSet      FieldType = "set"
	FieldTypeString   FieldType = "string"
	FieldTypeSecret   FieldType = "secret"
)

type Field struct {
	Name        string
	Description string
	Type        FieldType
	Required    bool
	Element     string
	Default     string
	Declaration string
}
