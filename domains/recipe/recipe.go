package recipe

import (
	"context"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type Recipe struct {
	shared.Model
	Title string `json:"title"`
	Summary string `json:"summary"`
	Description string `json:"description"`
	User *user.User `json:"user"`
	Photos []*Photo `json:"photos"`
	Steps []*Step `json:"steps"`
	Ingredients []*Ingredient `json:"ingredients"`
}

type StatsRecipe struct {
	*Recipe
	ViewCount uint64 `json:"view_count"`
	LikeCount uint64 `json:"like_count"`
}

type IRecipeRepository interface {
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Recipe) ([]*Recipe, *shared.CommonLimiter, error)
	Count(ctx context.Context, filters ...*Recipe) (int64, error)
	Create(ctx context.Context, recipe *Recipe) error
	Update(ctx context.Context, recipe *Recipe) error
	Delete(ctx context.Context, recipe *Recipe) error
	CreateBatch(ctx context.Context, recipes ...*Recipe) error
	DeleteBatch(ctx context.Context, recipes ...*Recipe) error
}

type IRecipeUseCase interface {
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Recipe) ([]*StatsRecipe, *shared.CommonLimiter, error)
	GetByUser(ctx context.Context, limiter *shared.CommonLimiter, usr *user.User, filters ...*Recipe) ([]*StatsRecipe, *shared.CommonLimiter, error)
	Create(ctx context.Context, usr *user.User, recipe *Recipe) error
	Update(ctx context.Context, usr *user.User, recipe *Recipe) error
	Delete(ctx context.Context, usr *user.User, recipe *Recipe) error
	CreateBatch(ctx context.Context, usr *user.User, recipes ...*Recipe) error
	DeleteBatch(ctx context.Context, usr *user.User, recipes ...*Recipe) error

	Popular(ctx context.Context, usr *user.User, limiter *shared.CommonLimiter, filters ...*Recipe) ([]*StatsRecipe, *shared.CommonLimiter, error)
	Like(ctx context.Context, usr *user.User, recipe *Recipe) error
	Dislike(ctx context.Context, usr *user.User, recipe *Recipe) error
	View(ctx context.Context, usr *user.User, recipe *Recipe) error

	AddStep(ctx context.Context, usr *user.User, recipe *Recipe, step string) error
	UpdateStep(ctx context.Context, usr *user.User, recipe *Recipe, sequence int, step string) error
	RemoveStep(ctx context.Context, usr *user.User, recipe *Recipe, sequence int) error

	AddIngredient(ctx context.Context, usr *user.User, recipe *Recipe, ingredient *Ingredient) error
	UpdateIngredient(ctx context.Context, usr *user.User, recipe *Recipe, ingredient *Ingredient) error
	RemoveIngredient(ctx context.Context, usr *user.User, recipe *Recipe, sequence int) error

	AddPhoto(ctx context.Context, usr *user.User, recipe *Recipe, ingredient *Photo) error
	UpdatePhoto(ctx context.Context, usr *user.User, recipe *Recipe, ingredient *Photo) error
	RemovePhoto(ctx context.Context, usr *user.User, recipe *Recipe, sequence int) error
}
