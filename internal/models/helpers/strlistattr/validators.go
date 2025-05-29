package strlistattr

import (
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var CommaSeparatedValidator = listvalidator.ValueStringsAre(stringvalidator.All(
	stringattr.NonEmptyValidator,
	stringvalidator.RegexMatches(regexp.MustCompile(`^[^,]*$`), "must not contain commas"),
))
