package requests

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"recipes/api/http"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"regexp"
)

type RecipeGet struct {
	http.CommonFilter
	ID int64 `uri:"recipe_id" json:"id"`
	Title string `form:"title" json:"title"`
	Summary string `form:"summary" json:"summary"`
	Description string `form:"description" json:"description"`
	Keyword string `form:"keyword" json:"keyword"`
}

func (r *RecipeGet) Validate(ctx context.Context) error {
	validString := regexp.MustCompile("w+")
	return validation.ValidateStructWithContext(ctx, r,
		validation.Field(
			&r.Title,
			validation.Match(validString),
		),
		validation.Field(
			&r.Description,
			validation.Match(validString),
		),
		validation.Field(
			&r.Summary,
			validation.Match(validString),
		),
		validation.Field(
			&r.Keyword,
			validation.Match(validString),
		),
	)
}
func (r *RecipeGet) Transform() (*shared.CommonLimiter, *recipe.Recipe) {
	if r.Keyword != "" && r.Title == "" {
		r.Title = r.Keyword
	}
	ltr := r.CommonFilter.Transform()
	if r.Keyword == "" && r.Title == "" && r.Description == "" && r.Summary == "" {
		return ltr, nil
	}
	return ltr, &recipe.Recipe{
		Title:       r.Title,
		Summary:     r.Summary,
		Description: r.Description,
		Photos:      nil,
		Steps:       nil,
		Ingredients: nil,
	}
}
