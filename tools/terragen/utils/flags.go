package utils

import "flag"

var Flags = struct {
	Verbose       bool
	SkipValidate  bool
	AddConnectors bool
}{}

func ParseFlags() {
	flag.BoolVar(&Flags.Verbose, "verbose", false, "set to true to print verbose logs")
	flag.BoolVar(&Flags.SkipValidate, "skip-validate", false, "set to true to not fail on validation errors")
	flag.BoolVar(&Flags.AddConnectors, "add-connectors", false, "set to true to add new connectors from the templates")
	flag.Parse()
}
