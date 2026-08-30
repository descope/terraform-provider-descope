package jsonattr

import (
	"bytes"
	"encoding/json"
	"errors"
	"maps"
	"reflect"
	"slices"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This package is for string attributes that hold a JSON document. Terraform compares attribute values as strings, so a document that
// differs only in whitespace or key order shows up as a change, e.g. an indented configuration file against the compact server response.

type Type = types.String

func Value(value string) Type {
	return types.StringValue(value)
}

func Required(fields ...string) schema.StringAttribute {
	return stringattr.Required(stringattr.JSONValidator(fields...), UseStateWhenEquivalent())
}

func Default(value string, fields ...string) schema.StringAttribute {
	return stringattr.Default(value, stringattr.JSONValidator(fields...), UseStateWhenEquivalent())
}

// Get unmarshals the attribute's document into data[key], or into data itself for helpers.RootKey when the document is the entire request.
// The schema already validated the attribute, so a value that doesn't parse here is a programming error rather than bad user input.
func Get(s Type, data map[string]any, key string) {
	value, err := decode(s.ValueString())
	if err != nil {
		panic("Invalid JSON attribute after validation: " + err.Error())
	}
	if key == helpers.RootKey {
		object, ok := value.(map[string]any)
		if !ok {
			panic("JSON attribute at the request root is not an object")
		}
		maps.Copy(data, object)
		return
	}
	data[key] = value
}

type SetOption int

const (
	// SkipIfAlreadySet keeps a document the configuration already holds, for resources whose apply returns a document differing in content.
	SkipIfAlreadySet SetOption = iota
)

// Set overwrites the attribute from data[key], or from data itself for helpers.RootKey, using an empty object when the key holds nothing.
// An equivalent document is left alone, so a refresh doesn't reformat what the configuration wrote.
func Set(s *Type, data map[string]any, key string, options ...SetOption) {
	if slices.Contains(options, SkipIfAlreadySet) && s.ValueString() != "" {
		return
	}
	if key == helpers.RootKey {
		store(s, data)
	} else if value := data[key]; value != nil {
		store(s, value)
	} else if s.ValueString() == "" {
		*s = Value("{}")
	}
}

// store writes the document only when it differs from the existing one, so it keeps its formatting when the server returns the same one back.
func store(s *Type, value any) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return // keep the existing value rather than break state
	}
	if !equivalent(s.ValueString(), string(encoded)) {
		*s = Value(string(encoded))
	}
}

// equivalent reports whether two strings hold the same JSON document, ignoring whitespace and key order. Numbers are compared as literals
// rather than as float64, so identifiers too large for a float aren't reported as equal; 1 and 1.0 compare as different in return.
func equivalent(a, b string) bool {
	if a == b {
		return true
	}
	av, err := decode(a)
	if err != nil {
		return false
	}
	bv, err := decode(b)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(av, bv)
}

func decode(s string) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(s)))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	if decoder.More() {
		return nil, errors.New("unexpected content after the JSON document")
	}
	return value, nil
}
