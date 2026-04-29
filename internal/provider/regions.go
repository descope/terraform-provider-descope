package provider

import "strings"

var friendlyToURL = map[string]string{
	"us": "https://api.descope.com",
	"eu": "https://api.euc1.descope.com",
	"ca": "https://api.cac1.descope.com",
	"ap": "https://api.aps2.descope.com",
}

var symbolToURL = map[string]string{
	"use1": "https://api.descope.com",
	"euc1": "https://api.euc1.descope.com",
	"cac1": "https://api.cac1.descope.com",
	"aps2": "https://api.aps2.descope.com",
}

func supportedRegionIDs() []string {
	return []string{"us", "eu", "ca", "ap"}
}

func regionAttributeDescription() string {
	ids := supportedRegionIDs()
	quoted := make([]string, len(ids))
	for i, id := range ids {
		quoted[i] = "`" + id + "`"
	}
	return strings.Join(quoted, ", ")
}

func resolveRegionURL(friendly string) (string, bool) {
	url, ok := friendlyToURL[friendly]
	return url, ok
}

func resolveSymbolURL(symbol string) (string, bool) {
	url, ok := symbolToURL[symbol]
	return url, ok
}

// resolveBaseURL picks a base URL in this order:
//   - explicit base URL
//   - explicit region
//   - auto-detected from management key
//
// The error string is non-empty when region is set to an
// unsupported value.
func resolveBaseURL(baseURL, region, managementKey string) (string, string) {
	if baseURL != "" {
		return baseURL, ""
	}
	if region != "" {
		url, ok := resolveRegionURL(region)
		if !ok {
			return "", "The Descope region '" + region + "' is not supported. Valid values: " + regionAttributeDescription() + "."
		}
		return url, ""
	}
	if symbol := regionSymbolFromMgmtKey(managementKey); symbol != "" {
		if url, ok := resolveSymbolURL(symbol); ok {
			return url, ""
		}
	}
	return "", ""
}

func detectRegionMismatch(configRegion, managementKey string) string {
	if configRegion == "" {
		return ""
	}
	keySymbol := regionSymbolFromMgmtKey(managementKey)
	if keySymbol == "" {
		return ""
	}
	configURL, _ := resolveRegionURL(configRegion)
	keyURL, keyKnown := resolveSymbolURL(keySymbol)
	if !keyKnown || configURL == keyURL {
		return ""
	}
	return "The provider region '" + configRegion + "' does not match the region encoded in the management key and API requests will likely fail."
}

func regionSymbolFromMgmtKey(cleartext string) string {
	const regionEncodedLen = 75
	if len(cleartext) != regionEncodedLen || cleartext[0] != 'K' {
		return ""
	}
	return cleartext[1:5]
}
