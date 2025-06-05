package infra

import (
	"encoding/json"
	"fmt"
	"os"
)

var isShallow = os.Getenv("TF_LOG") != "DEBUG"

const (
	shallowDepth = 5
	maxLength    = 80
)

func debugRequest(v any) string {
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	b, _ = json.MarshalIndent(trimmedMap(m), "", "  ")
	return string(b)
}

func debugResponse(s string) string {
	var m map[string]any
	_ = json.Unmarshal([]byte(s), &m)
	b, _ := json.MarshalIndent(trimmedMap(m), "", "  ")
	return string(b)
}

func trimmedMap(m map[string]any) map[string]any {
	result := map[string]any{}
	copyMapShallow(result, m, 0)
	return result
}

func copyMapShallow(dest, src map[string]any, depth int) {
	for k, v := range src {
		value := v
		if srcmap, ok := v.(map[string]any); ok {
			if isShallow && depth >= shallowDepth {
				value = fmt.Sprintf("Map{len: %d}", len(srcmap))
			} else {
				destmap := map[string]any{}
				copyMapShallow(destmap, srcmap, depth+1)
				value = destmap
			}
		} else if srcstr, ok := v.(string); ok {
			if len(srcstr) > maxLength {
				value = srcstr[:maxLength] + "..."
			} else {
				value = srcstr
			}
		} else if srclist, ok := v.([]any); ok {
			if isShallow && depth >= shallowDepth {
				value = fmt.Sprintf("List[len: %d]", len(srclist))
			} else {
				var destlist []any
				copySliceShallow(&destlist, &srclist, depth+1)
				value = destlist
			}
		}
		dest[k] = value
	}
}

func copySliceShallow(dest, src *[]any, depth int) {
	for _, v := range *src {
		value := v
		if srcmap, ok := v.(map[string]any); ok {
			if isShallow && depth >= shallowDepth {
				value = fmt.Sprintf("Map{len: %d}", len(srcmap))
			} else {
				destmap := map[string]any{}
				copyMapShallow(destmap, srcmap, depth+1)
				value = destmap
			}
		} else if srcstr, ok := v.(string); ok {
			if len(srcstr) > maxLength {
				value = srcstr[:maxLength] + "..."
			} else {
				value = srcstr
			}
		} else if srclist, ok := v.([]any); ok {
			if isShallow && depth >= shallowDepth {
				value = fmt.Sprintf("List[len: %d]", len(srclist))
			} else {
				var destlist []any
				copySliceShallow(&destlist, &srclist, depth+1)
				value = destlist
			}
		}
		*dest = append(*dest, value)
	}
}
