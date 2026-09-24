package emit

import (
	"fmt"
	"regexp"
	"strings"
)

var labelInvalidChars = regexp.MustCompile(`[^a-z0-9_]+`)

var namePrefixPattern = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

type Labels struct {
	prefix string
	used   map[string]map[string]bool
}

func NewLabels(prefix string) *Labels {
	return &Labels{prefix: prefix, used: map[string]map[string]bool{}}
}

func (l *Labels) Assign(resourceType, name, fallback string) string {
	label := sanitizeLabel(name)
	if label == "" {
		label = sanitizeLabel(fallback)
	}
	if label == "" {
		label = "this"
	}

	used := l.used[resourceType]
	if used == nil {
		used = map[string]bool{}
		l.used[resourceType] = used
	}
	assigned := label
	for suffix := 2; used[assigned]; suffix++ {
		assigned = fmt.Sprintf("%s_%d", label, suffix)
	}
	used[assigned] = true
	return prefixed(l.prefix, assigned)
}

func ValidateNamePrefix(prefix string) error {
	if !namePrefixPattern.MatchString(prefix) {
		return fmt.Errorf("invalid name prefix %q: expected lowercase letters, digits and underscores, not starting with a digit, such as prod", prefix)
	}
	return nil
}

func prefixed(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "_" + name
}

func sanitizeLabel(name string) string {
	label := labelInvalidChars.ReplaceAllString(strings.ToLower(name), "_")
	label = strings.Trim(label, "_")
	if label != "" && label[0] >= '0' && label[0] <= '9' {
		label = "_" + label
	}
	return label
}
