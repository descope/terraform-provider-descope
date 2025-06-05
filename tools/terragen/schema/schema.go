package schema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

type Schema struct {
	Files    []*File
	Warnings []string
	Packages []string
}

func ParseSources(root string) *Schema {
	utils.Debug(0, "Scheme")
	utils.Debug(0, "======")
	s := &Schema{}
	s.parseDir(root, nil)
	utils.Debug(0, "")
	return s
}

func (s *Schema) parseDir(root string, dirs []string) {
	path := filepath.Join(root, filepath.Join(dirs...))
	info, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to read files from path %s: %s", path, err.Error())
	}

	slices.SortFunc(info, compareDirEntry)
	for _, entry := range info {
		name := entry.Name()
		fullpath := filepath.Join(path, name)
		if entry.IsDir() && !shouldIgnoreDir(fullpath) {
			utils.Debug(len(dirs), "+ %s:", name)
			s.Packages = append(s.Packages, strings.Join(append(dirs, name), "/"))
			s.parseDir(root, append(dirs, name))
		} else if !shouldIgnoreFile(fullpath) {
			s.addFile(root, dirs, name)
		}
	}
}

func (s *Schema) addFile(root string, dirs []string, filename string) {
	path := filepath.Join(root, filepath.Join(dirs...), filename)
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read source file at path %s: %s", path, err.Error())
	}

	fileset := token.NewFileSet()

	f, err := parser.ParseFile(fileset, "", source, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse source file at path %s: %s", path, err.Error())
	}

	// the package used to access the model struct
	pkg := "models"
	if len(dirs) > 0 {
		pkg = dirs[len(dirs)-1]
	}

	// the file will be created lazily on demand if we find any suitable models
	var file *File

	for _, decl := range f.Decls {
		// expect a generic var declaration at the root of the source file
		decl, ok := decl.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR {
			continue
		}

		// we only expect one spec, i.e., no composite `var ( ... )` declarations
		if len(decl.Specs) != 1 {
			continue
		}

		// filter expect a value spec type
		varSpec, ok := decl.Specs[0].(*ast.ValueSpec)
		if !ok {
			log.Fatalf("unexpected spec declaration in %s: wanted *ast.ValueSpec, found %T", path, decl.Specs[0])
		}

		// we expect exactly one name, i.e., no `var a, b = x, y` declarations
		if len(varSpec.Names) != 1 {
			continue
		}

		// we expect the name to have an "...Attributes" suffix
		varName := varSpec.Names[0].Name
		if !strings.HasSuffix(varName, "Attributes") {
			continue
		}

		// we expect exactly one value, i.e., no `var FooAttributes mytype` declarations
		if len(varSpec.Values) != 1 {
			continue
		}

		// we expect the value expression to be a composite literal
		value, ok := varSpec.Values[0].(*ast.CompositeLit)
		if !ok {
			log.Fatalf("unexpected declaration value for %s in %s: wanted *ast.CompositeLit, found %T", varName, path, varSpec.Values[0])
		}

		// the name of the model is simply the prefix of the schema variable
		modelName := strings.TrimSuffix(varName, "Attributes")

		// the model will be created lazily on demand if we find any suitable fields
		var model *Model

		// we only include models that have at least one non-identifier field
		for _, v := range value.Elts {
			pair, ok := v.(*ast.KeyValueExpr)
			if !ok {
				log.Fatalf("unexpected element in %s in %s: wanted *ast.KeyValueExpr, found %T", varName, path, v)
			}

			key, ok := pair.Key.(*ast.BasicLit)
			if !ok {
				log.Fatalf("unexpected key in %s in %s: wanted *ast.BasicLit, found %T", varName, path, pair.Key)
			}

			if key.Kind != token.STRING {
				log.Fatalf("unexpected non-string key in %s in %s", varName, path)
			}

			fieldName := strings.Trim(key.Value, `"`)

			value, ok := pair.Value.(*ast.CallExpr)
			if !ok {
				log.Fatalf("unexpected value for field %s in %s in %s: wanted *ast.CallExpr, found %T", fieldName, varName, path, pair.Value)
			}

			lead := value.Fun
			if i, ok := lead.(*ast.IndexExpr); ok {
				lead = i.X
			}

			f, ok := lead.(*ast.SelectorExpr)
			if !ok {
				log.Fatalf("unexpected function in value for field %s in %s in %s: wanted *ast.SelectorExpr, found %T", fieldName, varName, path, value.Fun)
			}

			field := &Field{
				Name:        fieldName,
				Declaration: string(source[fileset.Position(value.Pos()).Offset:fileset.Position(value.End()).Offset]),
			}

			field.Type, field.Element, ok = fieldTypeFromSelector(f)
			if !ok {
				log.Fatalf("unexpected package type in value for field %s in %s in %s: %v", fieldName, varName, path, f.X)
			}

			var variant string
			if s := f.Sel; s == nil {
				log.Fatalf("unexpected nil selector in value for field %s in %s in %s: %v", fieldName, varName, path, f.X)
			} else {
				variant = s.Name
			}

			if trimmed := strings.TrimPrefix(variant, "Secret"); field.Type == FieldTypeString && trimmed != variant {
				field.Type = FieldTypeSecret
				variant = trimmed
			}

			if trimmed := strings.TrimSuffix(variant, "Required"); trimmed != variant {
				field.Required = true
				variant = trimmed
			}

			if field.Type == FieldTypeObject || field.Type == FieldTypeMap || field.Type == FieldTypeList || field.Type == FieldTypeSet {
				if field.Element == "" {
					if len(value.Args) < 1 {
						log.Fatalf("unexpected empty arguments in %s field %s in %s in %s", field.Type, fieldName, varName, path)
					}
					for _, arg := range value.Args {
						switch a := arg.(type) {
						case *ast.Ident:
							field.Element = fmt.Sprintf("%s.%s", pkg, a.Name)
						case *ast.SelectorExpr:
							ident, _ := a.X.(*ast.Ident)
							field.Element = fmt.Sprintf("%s.%s", ident.Name, a.Sel.Name)
						default:
							log.Fatalf("unexpected element type in %s field %s in %s in %s: %T", field.Type, fieldName, varName, path, a)
						}
						if strings.HasSuffix(field.Element, "Attributes") {
							field.Element = strings.TrimSuffix(field.Element, "Attributes")
							if len(dirs) > 0 && !strings.Contains(field.Element, ".") {
								field.Element = dirs[len(dirs)-1] + "." + field.Element
							}
							break
						}
					}
				}
				if field.Element == "" {
					log.Fatalf("failed to find element type in %s field %s in %s in %s", field.Type, fieldName, varName, path)
				}
			}

			if variant == "Default" && (field.Type == FieldTypeBool || field.Type == FieldTypeFloat || field.Type == FieldTypeInt || field.Type == FieldTypeString) {
				if len(value.Args) < 1 {
					log.Fatalf("unexpected missing default in %s field %s in %s in %s", field.Type, fieldName, varName, path)
				}
				switch a := value.Args[0].(type) {
				case *ast.BasicLit:
					if a.Value != `""` && a.Value != "0" {
						field.Default = a.Value
					}
				case *ast.Ident:
					if a.Name != "false" {
						field.Default = a.Name
					}
				default:
					log.Fatalf("unexpected default type in %s field %s in %s in %s: %T", field.Type, fieldName, varName, path, a)
				}
				variant = ""
			}

			// ignore identifier fields as they're not entered by the user
			if strings.HasPrefix(variant, "Identifier") {
				continue
			}

			// create the file and it to the schema if needed, now that we know for sure we found at least one model
			if file == nil {
				utils.Debug(len(dirs), "- %s", filename)
				file = &File{
					Name: strings.TrimSuffix(filename, ".go"),
					Dirs: dirs,
				}
				s.Files = append(s.Files, file)
			}

			// create the model if needed, now that we know for sure we found at least one field
			if model == nil {
				utils.Debug(len(dirs)+1, "- %s", modelName)
				model = &Model{
					Name:    modelName,
					Package: pkg,
				}
				file.Models = append(file.Models, model)
			}

			// keep the models and fields in the same order as they are in the original source file
			model.Fields = append(model.Fields, field)

			utils.Debug(len(dirs)+2, "- %s", fieldName)
			if field.Required {
				utils.Debug(len(dirs)+3, "- Type:   \t%s (required)", field.Type)
			} else {
				utils.Debug(len(dirs)+3, "- Type:   \t%s (optional)", field.Type)
			}
			if field.Element != "" {
				utils.Debug(len(dirs)+3, "- Element:\t%s", field.Element)
			}
			if field.Default != "" {
				utils.Debug(len(dirs)+3, "- Default:\t%s", field.Default)
			}
			utils.Debug(len(dirs)+3, "- Declare:\t%s", field.Declaration)
		}
	}
}

func (s *Schema) ValidateIfNeeded() {
	utils.Debug(0, "Validation")
	utils.Debug(0, "==========")

	missing := 0
	for _, f := range s.Files {
		for _, m := range f.Models {
			for _, field := range m.Fields {
				if field.Description == "" {
					if !utils.Flags.SkipValidate {
						fmt.Printf("[warning] missing documentation in %s.md: %s.%s\n", f.Name, m.Name, field.Name)
					}
					missing += 1
				}
			}
		}
	}

	label := "warning"
	if !utils.Flags.SkipValidate {
		label = "error"
	}

	if len(s.Warnings) > 0 {
		for _, w := range s.Warnings {
			fmt.Printf("[%s] %s\n", label, w)
		}
	}

	if len(s.Warnings) > 0 || missing > 0 {
		fmt.Printf("[%s] schema missing documentation for %d fields\n", label, missing)
		if !utils.Flags.SkipValidate {
			os.Exit(1)
		}
	}
}

func shouldIgnoreDir(path string) bool {
	return strings.HasSuffix(path, "/models/helpers") || strings.HasSuffix(path, "/tests")
}

func shouldIgnoreFile(path string) bool {
	return !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.HasPrefix(filepath.Base(path), ".")
}

// converts package path in selector to a field type
func fieldTypeFromSelector(selector *ast.SelectorExpr) (FieldType, string, bool) {
	if pkg, ok := selector.X.(*ast.Ident); ok {
		typ := FieldType(strings.TrimSuffix(pkg.Name, "attr"))
		if typ == FieldTypeBool || typ == FieldTypeDuration || typ == FieldTypeFloat || typ == FieldTypeInt || typ == FieldTypeList || typ == FieldTypeSet || typ == FieldTypeMap || typ == FieldTypeString {
			return typ, "", true
		}
		if typ == "strlist" {
			return FieldTypeList, "string", true
		}
		if typ == "strset" {
			return FieldTypeSet, "string", true
		}
		if typ == "strmap" {
			return FieldTypeMap, "string", true
		}
		if typ == "obj" {
			return FieldTypeObject, "", true
		}
	}
	return "", "", false
}

// sorts directory entries by files first, directories later, then in lexical order
func compareDirEntry(a, b fs.DirEntry) int {
	if a.IsDir() == b.IsDir() {
		return strings.Compare(a.Name(), b.Name())
	}
	if a.IsDir() {
		return 1
	}
	return -1
}
