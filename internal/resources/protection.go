package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	deletionProtectionAttribute = "deletion_protection"
	deletionProtectionError     = "Deletion Protection Enabled"
)

var (
	_ protectionSource = &tfsdk.State{}
	_ protectionSource = &tfsdk.Plan{}
	_ protectionSource = &tfsdk.Config{}
)

type protectionSource interface {
	Get(ctx context.Context, target any) diag.Diagnostics
	GetAttribute(ctx context.Context, attrPath path.Path, target any) diag.Diagnostics
}

func getDeletionProtection(ctx context.Context, source protectionSource, diags *diag.Diagnostics) types.Bool {
	flag := types.BoolNull()
	diags.Append(source.GetAttribute(ctx, path.Root(deletionProtectionAttribute), &flag)...)
	return flag
}

// Resolves the tristate flag: an explicit value wins, null consults the model's dynamic default, and no default means unprotected.
func resolveDeletionProtection(ctx context.Context, flag types.Bool, source protectionSource, model any, diags *diag.Diagnostics) (protected bool, viaDefault bool) {
	if !flag.IsNull() && !flag.IsUnknown() {
		return flag.ValueBool(), false
	}
	defaulter, ok := model.(helpers.DeletionProtectionDefaulter)
	if !ok {
		return false, false
	}
	diags.Append(source.Get(ctx, model)...)
	if diags.HasError() {
		return false, false
	}
	return defaulter.DeletionProtectionDefault(ctx), true
}

// Errors if the state resource is protected; used for destroy plans and in Delete, as forced replacements aren't visible at plan time.
func checkDestroyProtection(ctx context.Context, state protectionSource, model any, name string, diags *diag.Diagnostics) {
	flag := getDeletionProtection(ctx, state, diags)
	if protected, viaDefault := resolveDeletionProtection(ctx, flag, state, model, diags); protected {
		detail := "This " + name + " resource cannot be destroyed because deletion protection is enabled."
		if viaDefault {
			detail += " Deletion protection was enabled by default because the deletion_protection attribute isn't set."
		}
		detail += " To destroy this resource, first set the deletion_protection attribute to false in the resource configuration and run terraform apply, then retry the destroy operation."
		diags.AddError(deletionProtectionError, detail)
	}
}

// Errors if a resource planned for replacement is protected. Flag is read from state since Delete can't tell a replacement from a destroy.
func checkReplaceProtection(ctx context.Context, state protectionSource, model any, name string, diags *diag.Diagnostics) {
	flag := getDeletionProtection(ctx, state, diags)
	if protected, viaDefault := resolveDeletionProtection(ctx, flag, state, model, diags); protected {
		detail := "This plan requires the " + name + " resource to be replaced (destroyed and recreated), but deletion protection is enabled."
		if viaDefault {
			detail += " Deletion protection was enabled by default because the deletion_protection attribute isn't set."
		}
		detail += " To proceed, either revert the change that requires the resource to be replaced, or set the deletion_protection attribute to false in the resource configuration and apply that change first."
		diags.AddError(deletionProtectionError, detail)
	}
}

// Errors if a not-found resource is protected, instead of dropping it from state: the entity might not
// even be deleted, e.g. when a wrong management key or base URL makes it look not found.
func checkRemovalProtection(ctx context.Context, state protectionSource, model any, name string, diags *diag.Diagnostics) {
	flag := getDeletionProtection(ctx, state, diags)
	if protected, viaDefault := resolveDeletionProtection(ctx, flag, state, model, diags); protected {
		detail := "This " + name + " resource was not found on the backend, and it was not removed from the Terraform state because deletion protection is enabled."
		if viaDefault {
			detail += " Deletion protection was enabled by default because the deletion_protection attribute isn't set."
		}
		detail += " If the " + name + " was deleted intentionally, remove it from the state with the terraform state rm command and retry the operation. Otherwise, make sure the provider is configured with the correct management key and base URL, as an unexpectedly missing resource usually means requests are reaching the wrong environment or company."
		diags.AddError("Protected Resource Not Found", detail)
	}
}

// Attributes that scope a resource to its parent and must never change once set.
var immutableAttributes = []string{"project_id", "app_id"}

// Rejects plans that change an immutable scoping attribute. Their RequiresReplace modifiers stay on the attributes as defense in
// depth: without them a regression here would plan an in-place update that silently rescopes the resource.
func checkImmutableAttributes(ctx context.Context, sc schema.Schema, name string, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	for _, attr := range immutableAttributes {
		if _, ok := sc.Attributes[attr]; !ok {
			continue
		}
		attrPath := path.Root(attr)
		var state, plan types.String
		var diags diag.Diagnostics
		diags.Append(req.State.GetAttribute(ctx, attrPath, &state)...)
		diags.Append(req.Plan.GetAttribute(ctx, attrPath, &plan)...)
		if diags.HasError() || state.IsNull() || state.IsUnknown() || plan.IsNull() || plan.IsUnknown() {
			continue
		}
		if plan.ValueString() != state.ValueString() {
			resp.Diagnostics.AddAttributeError(attrPath, "Immutable Attribute Changed", "The "+attr+" attribute of a "+name+" resource cannot be changed after the resource is created. Create a new resource instead of changing the "+attr+" of an existing one.")
		}
	}
}

// Reports whether any top level string or bool attribute's plan modifiers require replacement; see TestRequiresReplaceModifiers.
func isPlannedReplace(ctx context.Context, sc schema.Schema, req resource.ModifyPlanRequest) bool {
	for name, attr := range sc.Attributes {
		attrPath := path.Root(name)
		var diags diag.Diagnostics
		switch a := attr.(type) {
		case schema.StringAttribute:
			if len(a.PlanModifiers) == 0 {
				continue
			}
			mreq := planmodifier.StringRequest{Path: attrPath, Plan: req.Plan, State: req.State, Config: req.Config}
			diags.Append(req.Plan.GetAttribute(ctx, attrPath, &mreq.PlanValue)...)
			diags.Append(req.State.GetAttribute(ctx, attrPath, &mreq.StateValue)...)
			diags.Append(req.Config.GetAttribute(ctx, attrPath, &mreq.ConfigValue)...)
			if diags.HasError() {
				continue
			}
			for _, m := range a.PlanModifiers {
				mresp := planmodifier.StringResponse{PlanValue: mreq.PlanValue}
				m.PlanModifyString(ctx, mreq, &mresp)
				if mresp.RequiresReplace {
					return true
				}
			}
		case schema.BoolAttribute:
			if len(a.PlanModifiers) == 0 {
				continue
			}
			mreq := planmodifier.BoolRequest{Path: attrPath, Plan: req.Plan, State: req.State, Config: req.Config}
			diags.Append(req.Plan.GetAttribute(ctx, attrPath, &mreq.PlanValue)...)
			diags.Append(req.State.GetAttribute(ctx, attrPath, &mreq.StateValue)...)
			diags.Append(req.Config.GetAttribute(ctx, attrPath, &mreq.ConfigValue)...)
			if diags.HasError() {
				continue
			}
			for _, m := range a.PlanModifiers {
				mresp := planmodifier.BoolResponse{PlanValue: mreq.PlanValue}
				m.PlanModifyBool(ctx, mreq, &mresp)
				if mresp.RequiresReplace {
					return true
				}
			}
		}
	}
	return false
}
