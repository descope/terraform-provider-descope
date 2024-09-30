package docgen

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

const breakBlankLines = 3

func MergeDocs(root string, sc *schema.Schema) {
	utils.Debug(0, "MergeDocs")
	utils.Debug(0, "=========")
	for _, file := range sc.Files {
		path := filepath.Join(root, filepath.Join(file.Dirs...), file.Name+".md")
		utils.Debug(len(file.Dirs), "+ %s.md:", file.Name)

		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				if !file.SkipDocs() {
					sc.Warnings = append(sc.Warnings, fmt.Sprintf("No existing model documentation file was found: %s", filepath.Join(filepath.Join(file.Dirs...), file.Name+".md")))
				}
				continue
			}
			log.Fatalf("failed to read documentation file at path %s: %s", path, err.Error())
		}

		s := bufio.NewScanner(bytes.NewReader(data))

		model, isModel := scanModelOrField(path, s)
		if !isModel {
			log.Fatalf("expected initial model in documentation file at path %s: %s", path, err.Error())
		}
		utils.Debug(len(file.Dirs)+1, "- %s:", model)

		fields := map[string]string{}

		for {
			next, isModel := scanModelOrField(path, s)
			if next == "" || isModel {
				updateModelDocs(file, model, fields)
				if !isModel {
					break
				}
				model = next
				utils.Debug(len(file.Dirs)+1, "- %s:", model)
				continue
			}
			expectFieldNotes(path, s)
			lines := scanUntilBreak(path, s)
			if len(lines) == 1 && strings.HasPrefix(lines[0], utils.PlaceholderDescription) {
				utils.Debug(len(file.Dirs)+2, "- %s: ???", next)
				continue
			}
			fields[next] = strings.Join(lines, "\n")
			utils.Debug(len(file.Dirs)+2, `- %s: "%s"`, next, fields[next])
		}
	}
	utils.Debug(0, "")
}

func updateModelDocs(file *schema.File, modelName string, modelFields map[string]string) {
	for _, model := range file.Models {
		if model.Name == modelName {
			for _, field := range model.Fields {
				if value := modelFields[field.Name]; value != "" {
					field.Description = value
				}
			}
		}
	}
}

func expectFieldNotes(path string, s *bufio.Scanner) {
	state := "initial"
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		switch {
		case len(line) == 0:
			if state == "munch" {
				return
			}
			log.Fatalf("expected field notes when parsing field in documentation file at path %s", path)
		case strings.HasPrefix(line, "- "):
			state = "munch"
		default:
			log.Fatalf("unexpected line when parsing field in documentation file at path %s: `%s`", path, line)
		}
	}
	log.Fatalf("unexpected end of file when when parsing field in documentation file at path %s", path)
}

func scanModelOrField(path string, s *bufio.Scanner) (name string, isModel bool) {
	state := "initial"
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		switch {
		case len(line) == 0:
			if state == "done" {
				return
			}
			if state != "initial" {
				log.Fatalf("unexpected blank line in documentation file at path %s", path)
			}
		case state == "initial":
			name = line
			state = "separator"
		case state == "separator":
			if strings.HasPrefix(line, "===") {
				isModel = true
			} else if !strings.HasPrefix(line, "---") {
				log.Fatalf("expected proper separator line in documentation file at path %s", path)
			}
			state = "done"
		default:
			log.Fatalf("unexpected non-blank line in documentation file at path %s", path)
		}
	}
	if state != "initial" && state != "done" {
		log.Fatalf("unexpected end of file in documentation file at path %s", path)
	}
	return
}

func scanUntilBreak(path string, s *bufio.Scanner) []string {
	lines := []string{}
	blanks := 0
	for blanks < breakBlankLines && s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			blanks += 1
		} else if blanks == 2 {
			log.Fatalf("unexpected number of blank lines in documentation file at path %s", path)
		} else {
			blanks = 0
		}
		lines = append(lines, line)
	}
	if err := s.Err(); err != nil {
		log.Fatalf("failed to read line from documentation file at path %s: %s", path, err.Error())
	}
	return trimBlankLines(lines)
}

func trimBlankLines(v []string) []string {
	for len(v) > 0 && v[0] == "" {
		v = v[1:]
	}
	for len(v) > 0 && v[len(v)-1] == "" {
		v = v[:len(v)-1]
	}
	return v
}
