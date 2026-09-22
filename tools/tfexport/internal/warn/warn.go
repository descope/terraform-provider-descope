package warn

import (
	"fmt"
	"slices"
)

// Warning is Lossy when an entity that exists is missing from the export because something failed, as opposed to a
// known limitation or a value the user has to supply after the export.
type Warning struct {
	Text  string
	Lossy bool
}

func Lost(format string, args ...any) Warning {
	return Warning{Text: fmt.Sprintf(format, args...), Lossy: true}
}

func Note(format string, args ...any) Warning {
	return Warning{Text: fmt.Sprintf(format, args...)}
}

func Incomplete(warnings []Warning) bool {
	return slices.ContainsFunc(warnings, func(w Warning) bool { return w.Lossy })
}
