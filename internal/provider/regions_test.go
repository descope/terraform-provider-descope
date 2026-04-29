package provider

import (
	"net/url"
	"strings"
	"testing"
)

func TestSymbolToURLValues(t *testing.T) {
	for symbol, raw := range symbolToURL {
		u, err := url.Parse(raw)
		if err != nil {
			t.Fatalf("symbol %q: invalid URL %q: %v", symbol, raw, err)
		}
		if u.Scheme != "https" {
			t.Fatalf("symbol %q: expected https scheme, got %q", symbol, u.Scheme)
		}
		if !strings.HasSuffix(u.Host, "descope.com") {
			t.Fatalf("symbol %q: expected host under descope.com, got %q", symbol, u.Host)
		}
	}
}

func TestFriendlyAndSymbolMapsCovered(t *testing.T) {
	for _, id := range supportedRegionIDs() {
		if _, ok := friendlyToURL[id]; !ok {
			t.Fatalf("supportedRegionIDs lists %q but it has no entry in friendlyToURL", id)
		}
	}
	for _, url := range friendlyToURL {
		found := false
		for _, symURL := range symbolToURL {
			if symURL == url {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("friendlyToURL value %q has no corresponding entry in symbolToURL", url)
		}
	}
}

func TestResolveRegionURL(t *testing.T) {
	tests := []struct {
		region   string
		ok       bool
		expected string
	}{
		{"us", true, "https://api.descope.com"},
		{"eu", true, "https://api.euc1.descope.com"},
		{"ca", true, "https://api.cac1.descope.com"},
		{"ap", true, "https://api.aps2.descope.com"},
		{"", false, ""},
		{"unknown", false, ""},
		{"au", false, ""}, // au was renamed to ap
		{"US", false, ""}, // case-sensitive on purpose
	}
	for _, tt := range tests {
		got, ok := resolveRegionURL(tt.region)
		if ok != tt.ok {
			t.Fatalf("resolveRegionURL(%q): got ok=%v, want %v", tt.region, ok, tt.ok)
		}
		if ok && got != tt.expected {
			t.Fatalf("resolveRegionURL(%q): got %q, want %q", tt.region, got, tt.expected)
		}
	}
}

func TestResolveBaseURL(t *testing.T) {
	euURL, _ := resolveRegionURL("eu")
	usURL, _ := resolveRegionURL("us")

	const (
		ksuid27  = "000000000000000000000000000"
		secret43 = "0000000000000000000000000000000000000000000"
	)
	euKey := "K" + "euc1" + ksuid27 + secret43
	caKey := "K" + "cac1" + ksuid27 + secret43
	legacyKey := "K" + ksuid27 + secret43

	tests := []struct {
		name           string
		baseURL        string
		region         string
		mgmtKey        string
		expectedURL    string
		errMustContain string
	}{
		{
			name:        "base_url wins over region and key",
			baseURL:     "https://custom.example.com",
			region:      "eu",
			mgmtKey:     caKey,
			expectedURL: "https://custom.example.com",
		},
		{
			name:        "region wins over key",
			region:      "eu",
			mgmtKey:     caKey,
			expectedURL: euURL,
		},
		{
			name:        "mgmt key region used when nothing else is set",
			mgmtKey:     euKey,
			expectedURL: euURL,
		},
		{
			name:        "use1 mgmt key (edge case) resolves to US default URL",
			mgmtKey:     "K" + "use1" + ksuid27 + secret43,
			expectedURL: usURL,
		},
		{
			name:        "legacy mgmt key falls through",
			mgmtKey:     legacyKey,
			expectedURL: "",
		},
		{
			name:        "no inputs falls through",
			expectedURL: "",
		},
		{
			name:           "unsupported region errors",
			region:         "moon",
			mgmtKey:        euKey,
			errMustContain: "moon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotErr := resolveBaseURL(tt.baseURL, tt.region, tt.mgmtKey)
			if tt.errMustContain != "" {
				if !strings.Contains(gotErr, tt.errMustContain) {
					t.Fatalf("expected error to contain %q, got %q", tt.errMustContain, gotErr)
				}
				return
			}
			if gotErr != "" {
				t.Fatalf("unexpected error: %q", gotErr)
			}
			if gotURL != tt.expectedURL {
				t.Fatalf("got URL %q, want %q", gotURL, tt.expectedURL)
			}
		})
	}
}

func TestDetectRegionMismatch(t *testing.T) {
	const (
		ksuid27  = "000000000000000000000000000"
		secret43 = "0000000000000000000000000000000000000000000"
	)
	euKey := "K" + "euc1" + ksuid27 + secret43
	caKey := "K" + "cac1" + ksuid27 + secret43
	use1Key := "K" + "use1" + ksuid27 + secret43
	unknownKey := "K" + "xxxx" + ksuid27 + secret43
	legacyKey := "K" + ksuid27 + secret43

	tests := []struct {
		name       string
		region     string
		key        string
		shouldWarn bool
	}{
		{"matching eu config + euc1 key", "eu", euKey, false},
		{"matching us config + use1 key (default URL alias)", "us", use1Key, false},
		{"mismatched us config + euc1 key", "us", euKey, true},
		{"mismatched eu config + cac1 key", "eu", caKey, true},
		{"mismatched eu config + use1 key", "eu", use1Key, true},
		{"no region config", "", euKey, false},
		{"legacy key (no embedded region)", "eu", legacyKey, false},
		{"unknown key region", "eu", unknownKey, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectRegionMismatch(tt.region, tt.key)
			if tt.shouldWarn && got == "" {
				t.Fatalf("expected a warning, got none")
			}
			if !tt.shouldWarn && got != "" {
				t.Fatalf("expected no warning, got %q", got)
			}
		})
	}
}

func TestRegionSymbolFromMgmtKey(t *testing.T) {
	// Synthetic fixtures: 27-char ksuid placeholder and 43-char secret placeholder.
	// Real keys are not used here — only length and prefix shape matter.
	const (
		ksuid27  = "000000000000000000000000000"
		secret43 = "0000000000000000000000000000000000000000000"
	)
	if len(ksuid27) != 27 {
		t.Fatalf("test fixture wrong: ksuid len %d", len(ksuid27))
	}
	if len(secret43) != 43 {
		t.Fatalf("test fixture wrong: secret len %d", len(secret43))
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"region-encoded EU", "K" + "euc1" + ksuid27 + secret43, "euc1"},
		{"region-encoded CA", "K" + "cac1" + ksuid27 + secret43, "cac1"},
		{"region-encoded AP", "K" + "aps2" + ksuid27 + secret43, "aps2"},
		{"region-encoded US (edge case)", "K" + "use1" + ksuid27 + secret43, "use1"},
		{"legacy 71-char key", "K" + ksuid27 + secret43, ""},
		{"empty", "", ""},
		{"too short", "K", ""},
		{"wrong prefix", "P" + "euc1" + ksuid27 + secret43, ""},
		{"one byte too long", "K" + "euc1" + ksuid27 + secret43 + "x", ""},
		{"one byte too short", "K" + "euc1" + ksuid27 + secret43[:len(secret43)-1], ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := regionSymbolFromMgmtKey(tt.input)
			if got != tt.expected {
				t.Fatalf("regionSymbolFromMgmtKey(%q): got %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
