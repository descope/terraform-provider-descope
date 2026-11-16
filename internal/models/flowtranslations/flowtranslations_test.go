package flowtranslations_test

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/flowtranslations"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

const translatedFlow = `"translation": {"connectorId": "manual", "enabled": true, "sourceLanguage": "en", "targetLanguages": ["he", "fr"]}`

func TestFlowTranslations(t *testing.T) {
	projectID := testacc.ProjectID(t)
	f := testacc.Flow(t)
	tr := testacc.FlowTranslations(t)
	flow := f.Block(`
		project_id = "` + projectID + `"
		flow_id = "testacc-translated"
		data = ` + testacc.FixtureJSON(t, "../flow/testdata/flow.json", `"translation": null`, translatedFlow),
	)
	translations := func(languages string) string {
		return flow + tr.Block(`
			project_id = "`+projectID+`"
			flow_id = `+f.Path()+`.flow_id
			languages = jsonencode(`+languages+`)
		`)
	}
	testacc.Run(t,
		// translate one text into each of the flow's two languages
		resource.TestStep{
			Config: translations(`{ he = { "HFf8.GLo1.google-button" = "המשך עם גוגל" }, fr = { "HFf8.GLo1.google-button" = "Continuer avec Google" } }`),
			Check: tr.Check(map[string]any{
				"id":         "testacc-translated",
				"project_id": testacc.AttributeIsSet,
				"flow_id":    "testacc-translated",
			}, resource.TestCheckResourceAttrWith(tr.Path(), "languages", languagesAre(map[string]map[string]string{
				"he": {"HFf8.GLo1.google-button": "המשך עם גוגל"},
				"fr": {"HFf8.GLo1.google-button": "Continuer avec Google"},
			}))),
		},
		// an import replaces every translation of the flow, so a language left out is removed
		resource.TestStep{
			Config: translations(`{ he = { "HFf8.GLo1.google-button" = "גוגל" } }`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(tr.Path(), plancheck.ResourceActionUpdate)},
			},
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrWith(tr.Path(), "languages", languagesAre(map[string]map[string]string{
					"he": {"HFf8.GLo1.google-button": "גוגל"},
				})),
				exportedLanguages(t, projectID, "testacc-translated", map[string]map[string]string{
					"he": {"HFf8.GLo1.google-button": "גוגל"},
				}),
			),
		},
		resource.TestStep{
			ResourceName:      tr.Path(),
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: testacc.GenerateImportStateID(tr.Path(), "project_id", "id"),
		},
		// destroying the resource clears the flow's translations and leaves the flow in place
		resource.TestStep{
			Config: flow,
			Check:  exportedLanguages(t, projectID, "testacc-translated", map[string]map[string]string{}),
		},
	)
}

func TestFlowTranslationsInvalidValues(t *testing.T) {
	projectID := testacc.ProjectID(t)
	tr := testacc.FlowTranslations(t)
	testacc.Run(t,
		resource.TestStep{
			Config: tr.Block(`
				project_id = "` + projectID + `"
				flow_id = "has spaces"
				languages = jsonencode({})
			`),
			ExpectError: regexp.MustCompile(`must only contain alphanumeric`),
		},
		resource.TestStep{
			Config: tr.Block(`
				project_id = "` + projectID + `"
				flow_id = "sign-up-or-in"
				languages = jsonencode({ he = "not an object" })
			`),
			ExpectError: regexp.MustCompile(`texts of each language`),
		},
		resource.TestStep{
			Config: tr.Block(`
				project_id = "` + projectID + `"
				flow_id = "sign-up-or-in"
				languages = jsonencode({ he = { "HFf8.GLo1.google-button" = 1 } })
			`),
			ExpectError: regexp.MustCompile(`texts of each language`),
		},
	)
}

func TestFlowTranslationsValues(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)
	m := &flowtranslations.FlowTranslationsModel{
		FlowID:    stringattr.Value("sign-up-or-in"),
		Languages: jsonattr.Value(`{"he": {"a.b.c": "שלום"}}`),
	}
	require.Equal(t, map[string]any{
		"flowId": "sign-up-or-in",
		"overrides": map[string]any{
			"flowId":    "sign-up-or-in",
			"languages": map[string]any{"he": map[string]any{"a.b.c": "שלום"}},
		},
	}, m.Values(h))
	require.False(t, diags.HasError())
}

func TestFlowTranslationsSetValues(t *testing.T) {
	var diags diag.Diagnostics
	h := helpers.NewHandler(context.Background(), &diags)

	m := &flowtranslations.FlowTranslationsModel{ID: stringattr.Value("sign-up-or-in"), FlowID: types.StringNull(), Languages: types.StringNull()}
	m.SetValues(h, map[string]any{"overrides": map[string]any{"flowId": "sign-up-or-in", "languages": map[string]any{"he": map[string]any{"a.b.c": "x"}}}})
	require.Equal(t, "sign-up-or-in", m.FlowID.ValueString())
	require.JSONEq(t, `{"he": {"a.b.c": "x"}}`, m.Languages.ValueString())

	m.SetValues(h, map[string]any{"overrides": map[string]any{"flowId": "sign-up-or-in", "languages": map[string]any{}}})
	require.JSONEq(t, `{}`, m.Languages.ValueString())
	m.Languages = jsonattr.Value(`{"he": {"a.b.c": "x"}}`)
	m.SetValues(h, map[string]any{})
	require.JSONEq(t, `{}`, m.Languages.ValueString())
	require.Equal(t, "sign-up-or-in", m.FlowID.ValueString())
	require.False(t, diags.HasError())
}

func TestFlowTranslationsValidate(t *testing.T) {
	for name, tc := range map[string]struct {
		languages string
		valid     bool
	}{
		"no languages":           {`{}`, true},
		"texts by key":           {`{"he": {"a.b.c": "x", "d.e.f": ""}, "fr": {}}`, true},
		"language not an object": {`{"he": "x"}`, false},
		"text not a string":      {`{"he": {"a.b.c": 1}}`, false},
		"nested deeper":          {`{"he": {"a.b.c": {"x": "y"}}}`, false},
		"not an object":          {`[]`, false},
	} {
		t.Run(name, func(t *testing.T) {
			var diags diag.Diagnostics
			m := &flowtranslations.FlowTranslationsModel{FlowID: stringattr.Value("f"), Languages: jsonattr.Value(tc.languages)}
			m.Validate(helpers.NewHandler(context.Background(), &diags))
			require.Equal(t, !tc.valid, diags.HasError(), "%v", diags)
		})
	}

	var diags diag.Diagnostics
	m := &flowtranslations.FlowTranslationsModel{FlowID: stringattr.Value("f"), Languages: types.StringUnknown()}
	m.Validate(helpers.NewHandler(context.Background(), &diags))
	require.False(t, diags.HasError(), "an unknown document is validated once it is known")
}

func languagesAre(expected map[string]map[string]string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var actual map[string]map[string]string
		if err := json.Unmarshal([]byte(value), &actual); err != nil {
			return err
		}
		if len(actual) != len(expected) {
			return mismatch(expected, actual)
		}
		for language, texts := range expected {
			if len(actual[language]) != len(texts) {
				return mismatch(expected, actual)
			}
			for key, text := range texts {
				if actual[language][key] != text {
					return mismatch(expected, actual)
				}
			}
		}
		return nil
	}
}

func exportedLanguages(t *testing.T, projectID, flowID string, expected map[string]map[string]string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		body := testacc.OutOfBandPostData(t, projectID, "/v1/mgmt/localization/flow/export", map[string]any{"flowId": flowID})
		overrides, _ := body["overrides"].(map[string]any)
		languages, _ := overrides["languages"].(map[string]any)
		if languages == nil {
			languages = map[string]any{}
		}
		encoded, err := json.Marshal(languages)
		if err != nil {
			return err
		}
		return languagesAre(expected)(string(encoded))
	}
}

func mismatch(expected, actual map[string]map[string]string) error {
	e, _ := json.Marshal(expected)
	a, _ := json.Marshal(actual)
	return errors.New("expected translations " + string(e) + ", got " + string(a))
}
