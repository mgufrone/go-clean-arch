package models

import (
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/domains/recipe"
	"recipes/domains/shared"
)

type IngredientModel struct {
	gorm.Model
	RecipeID uint64 `gorm:"index" json:"recipe_id"`
	Sequence int `gorm:"index" json:"sequence"`
	Type string `gorm:"index" json:"type"`
	Weight string  `gorm:"index" json:"weight"`
	Measurement string `gorm:"index" json:"measurement"`
	Name string `gorm:"index,class:FULLTEXT" json:"name"`
}

func (m *IngredientModel) Transform() *recipe.Ingredient {
	if m == nil {
		return nil
	}
	return &recipe.Ingredient{
		Model: shared.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Sequence:    m.Sequence,
		Type:        m.Type,
		Weight:        m.Weight,
		Measurement: m.Measurement,
		Name:        m.Name,
	}
}

func ParseIngredient(m *recipe.Ingredient) *IngredientModel {
	if m == nil {
		return nil
	}
	var (
		id uint
	)
	if m.ID != nil {
		id = uint(m.ID.(int))
	}
	var res *IngredientModel
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &res)
	res.ID = id
	return res
}

func (m *IngredientModel) ValidateWithContext(ctx context.Context) error {
	action := ctx.Value(common.ActionKey).(common.Action)
	return validation.ValidateStructWithContext(ctx, m,
		validation.Field(&m.Name,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
		),
		validation.Field(&m.Type,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			is.Alpha,
		),
		validation.Field(&m.Measurement,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			is.Alpha,
			validation.In("tbsp","ml","mg","g","l","tsp","oz", "cc"),
		),
		validation.Field(&m.Weight,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
		),
	)
}


