package strsetattr

import (
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var CommaSeparatedValidator = setvalidator.ValueStringsAre(stringvalidator.All(
	stringattr.NonEmptyValidator,
	stringvalidator.RegexMatches(regexp.MustCompile(`^[^,]*$`), "must not contain commas"),
))
