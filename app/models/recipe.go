package models

import (
	"context"
	"encoding/json"
	"github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type RecipeModel struct {
	gorm.Model
	Title       string `gorm:"index:,class:FULLTEXT" json:"title"`
	Summary     string `gorm:"index,class:FULLTEXT" json:"summary"`
	Description string `gorm:"index,class:FULLTEXT" json:"description"`
	UserID      uint `gorm:"index,priority:1" json:"user_id"`
	Photos      []*PhotoModel `gorm:"foreignKey:RecipeID" json:"photos"`
	Steps       []*StepModel `gorm:"foreignKey:RecipeID" json:"steps"`
	Ingredients []*IngredientModel `gorm:"foreignKey:RecipeID" json:"ingredients"`
}

func (m *RecipeModel) Transform() *recipe.Recipe {
	if m == nil {
		return nil
	}
	var (
		photos      []*recipe.Photo
		steps       []*recipe.Step
		ingredients []*recipe.Ingredient
	)
	for _, step := range m.Steps {
		steps = append(steps, step.Transform())
	}
	for _, photo := range m.Photos {
		photos = append(photos, photo.Transform())
	}
	for _, ingredient := range m.Ingredients {
		ingredients = append(ingredients, ingredient.Transform())
	}
	return &recipe.Recipe{
		Model: shared.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Title:       m.Title,
		Summary:     m.Summary,
		Description: m.Description,
		User: &user.User{Model: shared.Model{
			ID: m.UserID,
		}},
		Photos:     photos,
		Steps:      steps,
		Ingredients: ingredients,
	}
}
func ParseRecipe(m *recipe.Recipe) *RecipeModel {
	if m == nil {
		return nil
	}
	var (
		id    uint
		usrID uint
	)
	if m.ID != nil {
		id = m.ID.(uint)
	}
	if m.User != nil && m.User.ID != nil {
		usrID = m.User.ID.(uint)
	}
	var res *RecipeModel
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &res)
	res.ID = id
	res.UserID = usrID
	return res
}
func (m *RecipeModel) ValidateWithContext(ctx context.Context) error {
	action := ctx.Value(common.ActionKey).(common.Action)
	return validation.ValidateStructWithContext(ctx, m,
		validation.Field(&m.ID, validation.When(action == common.ActionUpdate, validation.Required, validation.NotNil)),
		validation.Field(&m.Title,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(20, 150),
		),
		validation.Field(&m.Summary,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(20, 300),
		),
		validation.Field(&m.Description,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(20, 5000),
		),
		validation.Field(&m.UserID,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
		),
		validation.Field(&m.Ingredients,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(1, 1000),
		),
		validation.Field(&m.Photos,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(1, 1000),
		),
		validation.Field(&m.Steps,
			validation.When(action == common.ActionCreate, validation.Required, validation.NotNil),
			validation.Length(1, 1000),
		),
	)
}
