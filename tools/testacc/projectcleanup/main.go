package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type projectList struct {
	Projects []project `json:"projects"`
}

type project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, input io.Reader, output io.Writer) error {
	if len(args) != 1 || strings.TrimSpace(args[0]) == "" {
		return errors.New("usage: projectcleanup <project-name-prefix>")
	}

	var projects projectList
	if err := json.NewDecoder(input).Decode(&projects); err != nil {
		return fmt.Errorf("decode project list: %w", err)
	}
	for _, project := range projects.Projects {
		if !strings.HasPrefix(project.Name, args[0]) {
			continue
		}
		if _, err := fmt.Fprintln(output, project.ID); err != nil {
			return fmt.Errorf("write project ID: %w", err)
		}
	}
	return nil
}
