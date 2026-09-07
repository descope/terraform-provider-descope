package migrate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
)

// VerifyResult is the outcome of checking an adoption plan against the migration manifest.
type VerifyResult struct {
	OK       bool
	Checked  int
	Imported int
	Findings []Finding
}

// Finding is one reason a plan is not an acceptable adoption plan.
type Finding struct {
	Kind    string
	Address string
	Actions []string
	Detail  string
}

// Finding kinds identify unsafe adoption-plan actions and missing verification prerequisites.
const (
	FindingUnexpectedCreate  = "unexpected_create"
	FindingUnexpectedUpdate  = "unexpected_update"
	FindingUnexpectedDelete  = "unexpected_delete"
	FindingUnexpectedReplace = "unexpected_replace"
	FindingForgottenObject   = "forgotten_object"
	FindingUnexpectedRead    = "unexpected_read"
	FindingMissingImport     = "missing_import"
	FindingUnexpectedImport  = "unexpected_import"
	FindingUnknownAction     = "unknown_action"
	FindingPlanErrored       = "plan_errored"
	FindingIncompletePlan    = "incomplete_plan"
	FindingDeferredChange    = "deferred_change"
	FindingActionInvocation  = "action_invocation"
	FindingProviderMismatch  = "provider_mismatch"
)

type adoptionPlan struct {
	FormatVersion     string                 `json:"format_version"`
	TerraformVersion  string                 `json:"terraform_version"`
	ResourceChanges   []adoptionPlanResource `json:"resource_changes"`
	DeferredChanges   []json.RawMessage      `json:"deferred_changes"`
	ActionInvocations []json.RawMessage      `json:"action_invocations"`
	Complete          *bool                  `json:"complete"`
	Errored           bool                   `json:"errored"`
}

type adoptionPlanResource struct {
	Address       string             `json:"address"`
	ModuleAddress string             `json:"module_address"`
	Mode          string             `json:"mode"`
	Type          string             `json:"type"`
	Name          string             `json:"name"`
	Index         any                `json:"index"`
	ProviderName  string             `json:"provider_name"`
	Change        adoptionPlanChange `json:"change"`
}

type adoptionPlanChange struct {
	Actions      []string       `json:"actions"`
	Importing    *planImport    `json:"importing"`
	Before       map[string]any `json:"before"`
	After        map[string]any `json:"after"`
	AfterUnknown map[string]any `json:"after_unknown"`
}

type planImport struct {
	ID string `json:"id"`
}

// VerifyPlan checks a `terraform show -json` adoption plan against the manifest. It never mutates either input.
func VerifyPlan(planJSON []byte, manifest *Manifest) (*VerifyResult, error) {
	if manifest == nil {
		return nil, fmt.Errorf("cannot verify an adoption plan without a manifest")
	}
	if !manifest.Ready {
		return nil, fmt.Errorf("cannot verify an adoption plan against a manifest that is not ready")
	}
	if manifest.TargetProvider.Source != ProviderSource {
		return nil, fmt.Errorf("manifest target provider source must be %q", ProviderSource)
	}

	plan := &adoptionPlan{}
	if err := json.Unmarshal(planJSON, plan); err != nil {
		return nil, fmt.Errorf("decoding Terraform plan JSON: %w", err)
	}
	major, _, _ := strings.Cut(plan.FormatVersion, ".")
	if major != "1" {
		return nil, fmt.Errorf("terraform plan format version %q is not supported, expected major version 1", plan.FormatVersion)
	}

	entries := make(map[string]Entry, len(manifest.Entries)+1)
	for _, entry := range manifest.Entries {
		entries[entry.Address] = entry
	}
	if manifest.LegacyAddress != "" {
		entries[manifest.LegacyAddress] = Entry{Address: manifest.LegacyAddress, Type: ProjectResourceType, ImportID: manifest.ProjectID}
	}
	expectedProvider := officialProviderSource

	result := &VerifyResult{}
	imported := make(map[string]bool, len(entries))
	for _, resource := range plan.ResourceChanges {
		if resource.Mode != "managed" {
			continue
		}
		result.Checked++
		if resource.Change.Importing != nil {
			result.Imported++
			imported[resource.Address] = true
		}
		entry, expected := entries[resource.Address]
		if resource.Change.Importing != nil && !expected {
			result.Findings = append(result.Findings, Finding{
				Kind:    FindingUnexpectedImport,
				Address: resource.Address,
				Actions: resource.Change.Actions,
				Detail:  "plan imports an address that is not in the migration manifest",
			})
			continue
		}
		if expected && resource.ProviderName != expectedProvider {
			result.Findings = append(result.Findings, Finding{
				Kind:    FindingProviderMismatch,
				Address: resource.Address,
				Actions: resource.Change.Actions,
				Detail:  fmt.Sprintf("plan uses provider %q, manifest requires %q", resource.ProviderName, expectedProvider),
			})
			continue
		}
		if resource.Change.Importing != nil && resource.Change.Importing.ID != entry.ImportID {
			result.Findings = append(result.Findings, Finding{
				Kind:    FindingUnexpectedImport,
				Address: resource.Address,
				Actions: resource.Change.Actions,
				Detail:  fmt.Sprintf("plan imports id %q, manifest requires %q", resource.Change.Importing.ID, entry.ImportID),
			})
			continue
		}
		if finding := verifyResourceChange(resource, entry, manifest.ProjectCore.DeletionProtection); finding != nil {
			result.Findings = append(result.Findings, *finding)
		}
	}

	for address := range entries {
		if !imported[address] {
			result.Findings = append(result.Findings, Finding{
				Kind:    FindingMissingImport,
				Address: address,
				Detail:  "manifest entry is not imported by the plan",
			})
		}
	}
	if plan.Errored {
		result.Findings = append(result.Findings, Finding{Kind: FindingPlanErrored, Detail: "Terraform reported an errored plan"})
	}
	if plan.Complete == nil || !*plan.Complete {
		result.Findings = append(result.Findings, Finding{Kind: FindingIncompletePlan, Detail: "Terraform plan completeness was not confirmed; use Terraform 1.8 or newer and do not use -target"})
	}
	if len(plan.DeferredChanges) > 0 {
		result.Findings = append(result.Findings, Finding{Kind: FindingDeferredChange, Detail: fmt.Sprintf("plan contains %d deferred change(s)", len(plan.DeferredChanges))})
	}
	if len(plan.ActionInvocations) > 0 {
		result.Findings = append(result.Findings, Finding{Kind: FindingActionInvocation, Detail: fmt.Sprintf("plan contains %d action invocation(s)", len(plan.ActionInvocations))})
	}

	sortFindings(result.Findings)
	result.OK = len(result.Findings) == 0
	return result, nil
}

func verifyResourceChange(resource adoptionPlanResource, entry Entry, projectDeletionProtection *bool) *Finding {
	actions := resource.Change.Actions
	if len(actions) == 1 && actions[0] == "no-op" {
		return nil
	}
	if len(actions) == 1 && actions[0] == "update" {
		return verifyImportedUpdate(resource, entry, projectDeletionProtection)
	}
	if isReplacement(actions) {
		return newActionFinding(FindingUnexpectedReplace, resource, "replacement is not allowed during adoption")
	}
	if len(actions) == 1 {
		switch actions[0] {
		case "create":
			return newActionFinding(FindingUnexpectedCreate, resource, "create is not allowed during adoption")
		case "delete":
			return newActionFinding(FindingUnexpectedDelete, resource, "delete is not allowed during adoption")
		case "forget":
			return newActionFinding(FindingForgottenObject, resource, "forget is not allowed during adoption")
		case "read":
			return newActionFinding(FindingUnexpectedRead, resource, "read is not allowed during adoption")
		}
	}

	for _, action := range actions {
		switch action {
		case "no-op", "create", "read", "update", "delete", "forget":
		default:
			return newActionFinding(FindingUnknownAction, resource, fmt.Sprintf("unrecognized action %q", action))
		}
	}
	return newActionFinding(FindingUnknownAction, resource, "unsupported action sequence")
}

func verifyImportedUpdate(resource adoptionPlanResource, entry Entry, projectDeletionProtection *bool) *Finding {
	if resource.Change.Importing == nil {
		return newActionFinding(FindingUnexpectedUpdate, resource, "update is not part of an import")
	}

	changed := changedAttributes(resource.Change)
	if len(changed) == 0 {
		return newActionFinding(FindingUnexpectedUpdate, resource, "update has no changed attributes")
	}
	allowed := make(map[string]struct{}, len(entry.SecretAttributes))
	for _, attribute := range entry.SecretAttributes {
		allowed[attribute] = struct{}{}
	}
	offending := make([]string, 0, len(changed))
	for _, attribute := range changed {
		if !allowedSecretPath(attribute, allowed) && !allowedProjectProtectionUpdate(attribute, resource.Change, entry, projectDeletionProtection) {
			offending = append(offending, attribute)
		}
	}
	if len(offending) == 0 {
		return nil
	}
	return newActionFinding(
		FindingUnexpectedUpdate,
		resource,
		"update changes non-secret attributes: "+strings.Join(offending, ", "),
	)
}

func allowedSecretPath(path string, allowed map[string]struct{}) bool {
	if _, ok := allowed[path]; ok {
		return true
	}
	for parent := range allowed {
		if strings.HasPrefix(path, parent+".") || strings.HasPrefix(path, parent+"[") {
			return true
		}
	}
	return false
}

func allowedProjectProtectionUpdate(attribute string, change adoptionPlanChange, entry Entry, expected *bool) bool {
	if entry.Type != ProjectResourceType || attribute != "deletion_protection" || expected == nil || hasUnknownValue(change.AfterUnknown[attribute]) {
		return false
	}
	actual, ok := change.After[attribute].(bool)
	return ok && actual == *expected
}

func changedAttributes(change adoptionPlanChange) []string {
	keys := make(map[string]struct{}, len(change.Before)+len(change.After)+len(change.AfterUnknown))
	for key := range change.Before {
		keys[key] = struct{}{}
	}
	for key := range change.After {
		keys[key] = struct{}{}
	}
	for key := range change.AfterUnknown {
		keys[key] = struct{}{}
	}

	changed := []string{}
	for key := range keys {
		before, beforeExists := change.Before[key]
		after, afterExists := change.After[key]
		collectChangedPaths(key, before, beforeExists, after, afterExists, change.AfterUnknown[key], &changed)
	}
	sort.Strings(changed)
	return slices.Compact(changed)
}

func collectChangedPaths(prefix string, before any, beforeExists bool, after any, afterExists bool, unknown any, changed *[]string) {
	if hasUnknownValue(unknown) {
		*changed = append(*changed, prefix)
		return
	}
	beforeMap, beforeMapOK := before.(map[string]any)
	afterMap, afterMapOK := after.(map[string]any)
	unknownMap, unknownMapOK := unknown.(map[string]any)
	if beforeMapOK || afterMapOK || unknownMapOK {
		keys := map[string]struct{}{}
		for key := range beforeMap {
			keys[key] = struct{}{}
		}
		for key := range afterMap {
			keys[key] = struct{}{}
		}
		for key := range unknownMap {
			keys[key] = struct{}{}
		}
		if len(keys) > 0 {
			for key := range keys {
				beforeValue, hasBefore := beforeMap[key]
				afterValue, hasAfter := afterMap[key]
				collectChangedPaths(prefix+"."+key, beforeValue, hasBefore, afterValue, hasAfter, unknownMap[key], changed)
			}
			return
		}
	}

	beforeList, beforeListOK := before.([]any)
	afterList, afterListOK := after.([]any)
	unknownList, unknownListOK := unknown.([]any)
	if beforeListOK || afterListOK || unknownListOK {
		length := max(len(beforeList), len(afterList), len(unknownList))
		for index := range length {
			var beforeValue, afterValue, unknownValue any
			hasBefore := index < len(beforeList)
			hasAfter := index < len(afterList)
			if hasBefore {
				beforeValue = beforeList[index]
			}
			if hasAfter {
				afterValue = afterList[index]
			}
			if index < len(unknownList) {
				unknownValue = unknownList[index]
			}
			collectChangedPaths(fmt.Sprintf("%s[%d]", prefix, index), beforeValue, hasBefore, afterValue, hasAfter, unknownValue, changed)
		}
		return
	}

	if beforeExists != afterExists || !reflect.DeepEqual(before, after) || hasUnknownValue(unknown) {
		*changed = append(*changed, prefix)
	}
}

func hasUnknownValue(value any) bool {
	flag, _ := value.(bool)
	return flag
}

func isReplacement(actions []string) bool {
	return len(actions) == 2 && ((actions[0] == "delete" && actions[1] == "create") ||
		(actions[0] == "create" && actions[1] == "delete"))
}

func newActionFinding(kind string, resource adoptionPlanResource, detail string) *Finding {
	return &Finding{
		Kind:    kind,
		Address: resource.Address,
		Actions: append([]string(nil), resource.Change.Actions...),
		Detail:  detail,
	}
}

func sortFindings(findings []Finding) {
	sort.SliceStable(findings, func(i, j int) bool {
		if findings[i].Kind != findings[j].Kind {
			return findings[i].Kind < findings[j].Kind
		}
		if findings[i].Address != findings[j].Address {
			return findings[i].Address < findings[j].Address
		}
		return findings[i].Detail < findings[j].Detail
	})
}

// Report renders a deterministic, human readable verification report.
func (r *VerifyResult) Report() string {
	status := "FAIL"
	if r.OK {
		status = "PASS"
	}

	var report strings.Builder
	fmt.Fprintf(&report, "adoption plan verification: %s\n", status)
	fmt.Fprintf(&report, "checked %d managed changes, %d imports\n", r.Checked, r.Imported)
	findings := append([]Finding(nil), r.Findings...)
	sortFindings(findings)
	for _, finding := range findings {
		fmt.Fprintf(&report, "  %s %s: %s", finding.Kind, finding.Address, finding.Detail)
		if len(finding.Actions) > 0 {
			fmt.Fprintf(&report, " [actions: %s]", strings.Join(finding.Actions, ","))
		}
		report.WriteByte('\n')
	}
	return report.String()
}
