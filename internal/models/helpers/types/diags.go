package types

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func Must[T any](x T, diags diag.Diagnostics) T {
	if errs := diags.Errors(); len(errs) > 0 {
		panic(fmt.Sprintf("%s: %s", errs[0].Summary(), errs[0].Detail()))
	}
	return x
}
