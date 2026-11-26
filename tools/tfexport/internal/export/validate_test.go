package export

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
)

func find(attrs []prune.Attr, name string) (prune.Attr, bool) {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr, true
		}
	}
	return prune.Attr{}, false
}

func TestWithPlaceholderTopLevel(t *testing.T) {
	attrs := []prune.Attr{{Name: "name"}, {Name: "client_secret"}}
	result := withPlaceholder(attrs, "client_secret")

	if secret, ok := find(result, "client_secret"); !ok || !secret.Placeholder {
		t.Errorf("expected client_secret to be a placeholder, got %+v", secret)
	}
	if secret, _ := find(attrs, "client_secret"); secret.Placeholder {
		t.Error("the original tree was modified")
	}
}

func TestWithPlaceholderAddsMissingAttribute(t *testing.T) {
	result := withPlaceholder([]prune.Attr{{Name: "name"}}, "api_key")

	if secret, ok := find(result, "api_key"); !ok || !secret.Placeholder {
		t.Errorf("expected api_key to be added as a placeholder, got %+v", secret)
	}
}

func TestWithPlaceholderNested(t *testing.T) {
	attrs := []prune.Attr{{Name: "auth", IsNested: true, Nested: []prune.Attr{{Name: "user"}, {Name: "password"}}}}
	result := withPlaceholder(attrs, "auth.password")

	auth, ok := find(result, "auth")
	if !ok || !auth.IsNested {
		t.Fatalf("expected a nested auth attribute, got %+v", auth)
	}
	if password, ok := find(auth.Nested, "password"); !ok || !password.Placeholder {
		t.Errorf("expected auth.password to be a placeholder, got %+v", password)
	}
	original, _ := find(attrs, "auth")
	if password, _ := find(original.Nested, "password"); password.Placeholder {
		t.Error("the original nested tree was modified")
	}
}

func TestWithPlaceholderCreatesMissingParent(t *testing.T) {
	result := withPlaceholder([]prune.Attr{{Name: "name"}}, "auth.password")

	auth, ok := find(result, "auth")
	if !ok || !auth.IsNested {
		t.Fatalf("expected auth to be created as a nested attribute, got %+v", auth)
	}
	if password, ok := find(auth.Nested, "password"); !ok || !password.Placeholder {
		t.Errorf("expected auth.password to be a placeholder, got %+v", password)
	}
}

func TestWithPlaceholderDeeplyNested(t *testing.T) {
	result := withPlaceholder(nil, "a.b.c")

	a, ok := find(result, "a")
	if !ok {
		t.Fatal("expected a")
	}
	b, ok := find(a.Nested, "b")
	if !ok {
		t.Fatal("expected a.b")
	}
	if c, ok := find(b.Nested, "c"); !ok || !c.Placeholder {
		t.Errorf("expected a.b.c to be a placeholder, got %+v", c)
	}
}

func TestWithPlaceholderElement(t *testing.T) {
	attrs := []prune.Attr{{Name: "headers", IsElements: true, Elements: [][]prune.Attr{
		{{Name: "key"}},
		{{Name: "key"}, {Name: "value"}},
	}}}
	result := withPlaceholder(attrs, "headers[1].value")

	headers, ok := find(result, "headers")
	if !ok || len(headers.Elements) != 2 {
		t.Fatalf("expected two elements, got %+v", headers)
	}
	if value, ok := find(headers.Elements[1], "value"); !ok || !value.Placeholder {
		t.Errorf("expected headers[1].value to be a placeholder, got %+v", value)
	}
	if _, ok := find(headers.Elements[0], "value"); ok {
		t.Error("the other element should be untouched")
	}
	original, _ := find(attrs, "headers")
	if value, _ := find(original.Elements[1], "value"); value.Placeholder {
		t.Error("the original tree was modified")
	}
}

func TestWithPlaceholderElementOutOfRange(t *testing.T) {
	attrs := []prune.Attr{{Name: "headers", IsElements: true, Elements: [][]prune.Attr{{{Name: "key"}}}}}
	if result := withPlaceholder(attrs, "headers[5].value"); len(result) != 1 {
		t.Errorf("expected the tree to be left alone, got %+v", result)
	}
}

func TestCutIndex(t *testing.T) {
	for _, test := range []struct {
		segment string
		name    string
		index   int
		indexed bool
	}{
		{segment: "headers[0]", name: "headers", index: 0, indexed: true},
		{segment: "headers[12]", name: "headers", index: 12, indexed: true},
		{segment: "headers", name: "headers"},
		{segment: "headers[x]", name: "headers[x]"},
		{segment: "headers[-1]", name: "headers[-1]"},
	} {
		name, index, indexed := cutIndex(test.segment)
		if name != test.name || index != test.index || indexed != test.indexed {
			t.Errorf("cutIndex(%q) = %q, %d, %v; expected %q, %d, %v", test.segment, name, index, indexed, test.name, test.index, test.indexed)
		}
	}
}
