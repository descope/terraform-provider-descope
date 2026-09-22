package infra

import (
	"encoding/json"
	"fmt"
	"os"
)

// payload logging is gated behind a dedicated variable rather than TF_LOG so that request and
// response payloads, which may carry secrets, can never end up in logs by accident
var isLogging = os.Getenv("TF_UNSAFE_LOGS") != ""

var isShallow = os.Getenv("TF_LOG") != "DEBUG" && os.Getenv("TF_LOG") != "TRACE"

const (
	shallowDepth = 5
	maxLength    = 80
)

// the debug helpers are evaluated as tflog arguments on every request, so they skip all
// the marshaling work unless logging was actually enabled for the terraform run
func debugRequest(v any) string {
	if !isLogging {
		return ""
	}
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	b, _ = json.MarshalIndent(trimmedMap(m), "", "  ")
	return string(b)
}

func debugResponse(s string) string {
	if !isLogging {
		return ""
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(s), &m)
	b, _ := json.MarshalIndent(trimmedMap(m), "", "  ")
	return string(b)
}

func trimmedMap(m map[string]any) map[string]any {
	result, _ := trimValue(m, 0).(map[string]any)
	return result
}

func trimValue(v any, depth int) any {
	switch src := v.(type) {
	case map[string]any:
		if isShallow && depth > shallowDepth {
			return fmt.Sprintf("Map{len: %d}", len(src))
		}
		dest := map[string]any{}
		for k, e := range src {
			dest[k] = trimValue(e, depth+1)
		}
		return dest
	case string:
		if len(src) > maxLength {
			return src[:maxLength] + "..."
		}
		return src
	case []any:
		if isShallow && depth > shallowDepth {
			return fmt.Sprintf("List[len: %d]", len(src))
		}
		var dest []any
		for _, e := range src {
			dest = append(dest, trimValue(e, depth+1))
		}
		return dest
	}
	return v
}
