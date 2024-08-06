package helpers

import (
	"fmt"
	"sync/atomic"
)

const (
	ConnectorReferenceKey   = "ref:connector"
	RoleReferenceKey        = "ref:role"
	JWTTemplateReferenceKey = "ref:jwttemplate"
)

type ModelReference struct {
	Type string
	ID   string
	Key  string
}

func (r *ModelReference) ReferenceValue() string {
	if r.ID != "" {
		return r.ID
	}
	return r.Key
}

func (r *ModelReference) ProviderValue() string {
	if r.ID != "" && r.Type != "" {
		return r.Type + ":" + r.ID
	}
	return r.ReferenceValue()
}

// References container

type ReferencesMap map[string]*ModelReference

func (r ReferencesMap) Add(key, typ, id, name string) {
	ref := &ModelReference{Type: typ, ID: id}
	if ref.ID == "" {
		ref.Key = generateReferenceKey(key)
	}
	refName := fmt.Sprintf("%s:%s", key, name)
	r[refName] = ref
}

func (r ReferencesMap) Get(key, name string) *ModelReference {
	if key == ConnectorReferenceKey && name == DescopeConnector {
		return &ModelReference{Key: DescopeConnector}
	}
	refName := fmt.Sprintf("%s:%s", key, name)
	return r[refName]
}

var referencesMapCounter atomic.Int64

func generateReferenceKey(key string) string {
	return fmt.Sprintf("%s:%d", key, referencesMapCounter.Add(1))
}
