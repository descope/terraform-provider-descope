package schema

type Model struct {
	Name      string
	Package   string
	Fields    []*Field
	Generated bool
}
