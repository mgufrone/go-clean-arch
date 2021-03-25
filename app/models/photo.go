package models

import (
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
	"recipes/domains/recipe"
	"recipes/domains/shared"
)

type PhotoModel struct {
	gorm.Model
	RecipeID uint64 `gorm:"index" json:"recipe_id"`
	Sequence uint `gorm:"index" json:"sequence"`
	URL string `json:"url"`
}

func (m *PhotoModel) Transform() *recipe.Photo {
	if m == nil {
		return nil
	}
	return &recipe.Photo{
		Model:    shared.Model{
			ID: m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Sequence: m.Sequence,
		URL: m.URL,
	}
}

func ParsePhoto(m *recipe.Photo) *PhotoModel {
	if m == nil {
		return nil
	}
	var res *PhotoModel
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, res)
	return res
}

func (m *PhotoModel) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, m,
		validation.Field(&m.URL, validation.Required, validation.NotNil, is.URL),
	)
}


