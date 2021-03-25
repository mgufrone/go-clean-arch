package recipe

import (
	"context"
	"recipes/domains/shared"
)

type Ingredient struct {
	shared.Model
	Sequence int `json:"sequence"`
	Type string `json:"type"`
	Weight string `json:"weight"`
	Measurement string `json:"measurement"`
	Name string `json:"name"`
}

type IIngredientRepository interface {
	GetBySequence(ctx context.Context, sequence uint) (*Ingredient, error)
	DeleteSequence(ctx context.Context, sequence uint) error
	Update(ctx context.Context, sequence uint, step *Ingredient) error
	Create(ctx context.Context, ingredient *Ingredient) error
}
