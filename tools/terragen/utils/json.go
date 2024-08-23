package utils

import (
	"encoding/json"
	"os"
)

func ReadJSON[T any](path string, target *T) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}

func WriteJSON[T any](path string, target *T) error {
	b, err := json.MarshalIndent(target, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0644)
}
