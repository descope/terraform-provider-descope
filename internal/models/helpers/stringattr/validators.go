package stringattr

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var TimeUnitValidator = stringvalidator.OneOf("seconds", "minutes", "hours", "days", "weeks")
