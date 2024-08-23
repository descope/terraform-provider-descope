package utils

import "flag"

var Flags = struct {
	Verbose       bool
	SkipValidate  bool
	SkipTemplates bool
}{}

func ParseFlags() {
	flag.BoolVar(&Flags.Verbose, "verbose", false, "set to true to print verbose logs")
	flag.BoolVar(&Flags.SkipValidate, "skip-validate", false, "set to true to not fail on validation errors")
	flag.BoolVar(&Flags.SkipTemplates, "skip-templates", false, "set to true to skip parsing templates")
	flag.Parse()
}
