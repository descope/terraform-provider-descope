package settings

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func provisioningHandler() *helpers.Handler {
	return helpers.NewHandler(context.Background(), &diag.Diagnostics{})
}

// require_invitation is the inverse of the server's projectSelfProvisioning flag. The flag
// must always be written so it is never dropped from the payload (an absent value is treated
// as false / self sign-up disabled by the server).
func TestInviteSelfProvisioningWrite(t *testing.T) {
	cases := []struct {
		name    string
		require types.Bool
		want    bool
	}{
		{"require_invitation=false", boolattr.Value(false), true},
		{"require_invitation=true", boolattr.Value(true), false},
		{"require_invitation=null", types.BoolNull(), true},
		{"require_invitation=unknown", types.BoolUnknown(), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := *InviteSettingsDefault
			m.RequireInvitation = c.require
			data := m.Values(provisioningHandler())
			got, present := data["projectSelfProvisioning"]
			if !present {
				t.Fatalf("projectSelfProvisioning must always be present in the payload")
			}
			if got != c.want {
				t.Fatalf("projectSelfProvisioning = %v, want %v", got, c.want)
			}
		})
	}
}

// The server may return projectSelfProvisioning as a bool or as a "true"/"false" string.
// Both encodings must round-trip; a string must not be dropped (which would force
// require_invitation to true and disable self sign-up on the next apply).
func TestInviteSelfProvisioningRead(t *testing.T) {
	cases := []struct {
		name string
		data map[string]any
		want bool // expected require_invitation
	}{
		{"bool true", map[string]any{"projectSelfProvisioning": true}, false},
		{"bool false", map[string]any{"projectSelfProvisioning": false}, true},
		{"string true", map[string]any{"projectSelfProvisioning": "true"}, false},
		{"string false", map[string]any{"projectSelfProvisioning": "false"}, true},
		{"absent", map[string]any{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := *InviteSettingsDefault
			m.RequireInvitation = types.BoolUnknown()
			m.SetValues(provisioningHandler(), c.data)
			if got := m.RequireInvitation.ValueBool(); got != c.want {
				t.Fatalf("require_invitation = %v, want %v", got, c.want)
			}
		})
	}
}
