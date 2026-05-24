package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// MTLSCredential

var MTLSCredentialAttributes = map[string]schema.Attribute{
	"id":         stringattr.Optional(),
	"name":       stringattr.Required(stringattr.StandardLenValidator),
	"subject_dn": stringattr.Default(""),
	"pem":        stringattr.Default(""),
	"thumbprint": stringattr.Optional(),
}

type MTLSCredentialModel struct {
	ID         stringattr.Type `tfsdk:"id"`
	Name       stringattr.Type `tfsdk:"name"`
	SubjectDN  stringattr.Type `tfsdk:"subject_dn"`
	PEM        stringattr.Type `tfsdk:"pem"`
	Thumbprint stringattr.Type `tfsdk:"thumbprint"`
}

func (m *MTLSCredentialModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.SubjectDN, data, "subjectDn")
	stringattr.Get(m.PEM, data, "pem")
	return data
}

func (m *MTLSCredentialModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.SubjectDN, data, "subjectDn")
	stringattr.Set(&m.PEM, data, "pem")
	stringattr.Set(&m.Thumbprint, data, "thumbprint")
}

func (m *MTLSCredentialModel) GetName() stringattr.Type { return m.Name }
func (m *MTLSCredentialModel) GetID() stringattr.Type   { return m.ID }
func (m *MTLSCredentialModel) SetID(id stringattr.Type) { m.ID = id }

// MTLSTLSClientAuth

var MTLSTLSClientAuthAttributes = map[string]schema.Attribute{
	"credentials": listattr.Default[MTLSCredentialModel](MTLSCredentialAttributes),
}

type MTLSTLSClientAuthModel struct {
	Credentials listattr.Type[MTLSCredentialModel] `tfsdk:"credentials"`
}

func (m *MTLSTLSClientAuthModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Credentials, data, "credentials", h)
	return data
}

func (m *MTLSTLSClientAuthModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.Credentials, data, "credentials", "name", h)
}

// MTLSSelfSignedTLSClientAuth

var MTLSSelfSignedTLSClientAuthAttributes = map[string]schema.Attribute{
	"credentials": listattr.Default[MTLSCredentialModel](MTLSCredentialAttributes),
}

type MTLSSelfSignedTLSClientAuthModel struct {
	Credentials listattr.Type[MTLSCredentialModel] `tfsdk:"credentials"`
}

func (m *MTLSSelfSignedTLSClientAuthModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.Credentials, data, "credentials", h)
	return data
}

func (m *MTLSSelfSignedTLSClientAuthModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.Credentials, data, "credentials", "name", h)
}

// MTLSClientAuthMethods

var MTLSClientAuthMethodsAttributes = map[string]schema.Attribute{
	"tls_client_auth":             objattr.Default[MTLSTLSClientAuthModel](nil, MTLSTLSClientAuthAttributes),
	"self_signed_tls_client_auth": objattr.Default[MTLSSelfSignedTLSClientAuthModel](nil, MTLSSelfSignedTLSClientAuthAttributes),
}

type MTLSClientAuthMethodsModel struct {
	TLSClientAuth           objattr.Type[MTLSTLSClientAuthModel]           `tfsdk:"tls_client_auth"`
	SelfSignedTLSClientAuth objattr.Type[MTLSSelfSignedTLSClientAuthModel] `tfsdk:"self_signed_tls_client_auth"`
}

func (m *MTLSClientAuthMethodsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	objattr.Get(m.TLSClientAuth, data, "tlsClientAuth", h)
	objattr.Get(m.SelfSignedTLSClientAuth, data, "selfSignedTlsClientAuth", h)
	return data
}

func (m *MTLSClientAuthMethodsModel) SetValues(h *helpers.Handler, data map[string]any) {
	objattr.Set(&m.TLSClientAuth, data, "tlsClientAuth", h)
	objattr.Set(&m.SelfSignedTLSClientAuth, data, "selfSignedTlsClientAuth", h)
}
