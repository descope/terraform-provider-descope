package policyrule

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func validateConditionValue(operator string, value json.RawMessage) error {
	trimmed := bytes.TrimSpace(value)
	if !json.Valid(trimmed) {
		return fmt.Errorf("condition value must be valid JSON")
	}

	isList := len(trimmed) > 0 && trimmed[0] == '['
	isListOperator := operator == "in" || operator == "notIn"
	if isListOperator && !isList {
		return fmt.Errorf("condition operator %q requires a list JSON value", operator)
	}
	if !isListOperator && isList {
		return fmt.Errorf("condition operator %q requires a scalar JSON value", operator)
	}
	if !isList && len(trimmed) > 0 && trimmed[0] == '{' {
		return fmt.Errorf("condition operator %q requires a scalar JSON value", operator)
	}
	return nil
}

func compactJSON(value json.RawMessage) (string, error) {
	var buffer bytes.Buffer
	if err := json.Compact(&buffer, value); err != nil {
		return "", fmt.Errorf("compact condition JSON: %w", err)
	}
	return buffer.String(), nil
}
