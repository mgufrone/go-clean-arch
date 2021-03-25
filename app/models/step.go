package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"recipes/domains/recipe"
	"recipes/domains/shared"
)

type StepModel struct {
	gorm.Model
	RecipeID uint64 `gorm:"index" json:"recipe_id"`
	Sequence uint `gorm:"index" json:"sequence"`
	Step string `gorm:"index,class:FULLTEXT" json:"step"`
}

func (m *StepModel) Transform() *recipe.Step {
	if m == nil {
		return nil
	}
	return &recipe.Step{
		Model:    shared.Model{
			ID: m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Sequence: m.Sequence,
		Step:     m.Step,
	}
}

func ParseStep(m *recipe.Step) *StepModel {
	if m == nil {
		return nil
	}
	var res *StepModel
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, res)
	return res
}


