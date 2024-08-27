package testacc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var AttributeIsSet = struct{ int }{1}
var AttributeIsNotSet = struct{ int }{-1}

func AttributeMatchesPattern(pattern string) func(string) error {
	re := regexp.MustCompile(pattern)
	return func(s string) error {
		if !re.MatchString(s) {
			return fmt.Errorf("expected value to match pattern '%s', got '%s'", pattern, s)
		}
		return nil
	}
}

func AttributeMatchesJSON(object string) func(string) error {
	expected, err := convertJSON(object)
	if err != nil {
		panic(err.Error())
	}
	return func(s string) error {
		actual, err := convertJSON(s)
		if err != nil {
			return err
		}
		if actual != expected {
			return fmt.Errorf("expected value to match JSON '%s', got '%s'", expected, actual)
		}
		return nil
	}
}

func AttributeHasPrefix(prefix string) func(string) error {
	return func(s string) error {
		if !strings.HasPrefix(s, prefix) {
			return fmt.Errorf("expected value with prefix %s, got '%s'", prefix, s)
		}
		return nil
	}
}

func convertJSON(s string) (string, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON attribute value: %s", err.Error())
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON attribute value: %s", err.Error())
	}
	return string(b), nil
}
