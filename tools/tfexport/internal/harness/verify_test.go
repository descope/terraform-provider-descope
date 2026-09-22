package harness

import "testing"

func TestStripComponentsVersion(t *testing.T) {
	tests := []struct {
		name       string
		first      string
		second     string
		wantEquals bool
	}{
		{
			name:       "styles differing only in the top level stamp",
			first:      `{"componentsVersion":"2.0.0","styles":{"primary":"#fff"}}`,
			second:     `{"componentsVersion":"3.14.10","styles":{"primary":"#fff"}}`,
			wantEquals: true,
		},
		{
			name:       "flow differing only in the metadata stamp",
			first:      `{"flowId":"sign-in","metadata":{"componentsVersion":"2.0.0","name":"Sign In"}}`,
			second:     `{"flowId":"sign-in","metadata":{"componentsVersion":"3.14.10","name":"Sign In"}}`,
			wantEquals: true,
		},
		{
			name:       "a real difference alongside a differing stamp still fails",
			first:      `{"componentsVersion":"2.0.0","styles":{"primary":"#fff"}}`,
			second:     `{"componentsVersion":"3.14.10","styles":{"primary":"#000"}}`,
			wantEquals: false,
		},
		{
			name:       "the same key nested elsewhere is not exempt",
			first:      `{"styles":{"theme":{"componentsVersion":"2.0.0"}}}`,
			second:     `{"styles":{"theme":{"componentsVersion":"3.14.10"}}}`,
			wantEquals: false,
		},
		{
			name:       "unparseable content is compared verbatim",
			first:      "not json",
			second:     "not json either",
			wantEquals: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			first, second := stripComponentsVersion(test.first), stripComponentsVersion(test.second)
			if equals := first == second; equals != test.wantEquals {
				t.Errorf("stripped equality is %v, expected %v\n  first:  %s\n  second: %s", equals, test.wantEquals, first, second)
			}
		})
	}
}
