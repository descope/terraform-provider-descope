package connectors

import (
	"slices"

	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
)

func (m *VonageConnectorModel) Validate(h *helpers.Handler) {
	smsChannelComplete := isVonageChannelComplete(m.APIKey, m.APISecret, m.Sender)
	voiceChannelComplete := isVonageChannelComplete(m.ApplicationID, m.PrivateKey, m.FromPhoneVoice)
	if !smsChannelComplete && !voiceChannelComplete {
		h.Conflict("The Vonage connector requires at least one complete channel: api_key, api_secret and sender for text messages, or application_id, private_key and from_phone_voice for voice calls")
	}
}

func isVonageChannelComplete(channelFields ...stringattr.Type) bool {
	return !slices.ContainsFunc(channelFields, func(field stringattr.Type) bool {
		return !field.IsUnknown() && field.ValueString() == ""
	})
}
