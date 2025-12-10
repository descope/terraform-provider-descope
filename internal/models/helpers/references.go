package helpers

import (
	"fmt"
	"sync/atomic"
)

const (
	ConnectorReferenceKey   = "ref:connector"
	RoleReferenceKey        = "ref:role"
	JWTTemplateReferenceKey = "ref:jwttemplate"
	ListReferenceKey        = "ref:list"
)

type ModelReference struct {
	Type string
	ID   string
	Key  string
	Name string
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
	ref := &ModelReference{Type: typ, ID: id, Name: name}
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

func (r ReferencesMap) Name(id string) string {
	if id == DescopeConnector {
		return DescopeConnector
	}
	for _, value := range r {
		if value.ID == id {
			return value.Name
		}
	}
	return ""
}

var referencesMapCounter atomic.Int64

func generateReferenceKey(key string) string {
	return fmt.Sprintf("%s:%d", key, referencesMapCounter.Add(1))
}
