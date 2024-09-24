package schema

type File struct {
	Name   string
	Dirs   []string
	Models []*Model
}

func (f *File) SkipDocs() bool {
	for _, m := range f.Models {
		if !m.Generated {
			return false
		}
	}
	return true
}
