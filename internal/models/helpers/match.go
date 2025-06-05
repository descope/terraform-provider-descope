package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MatchableModel[T any] interface {
	Model[T]
	GetName() types.String
	GetID() types.String
	SetID(id types.String)
}
