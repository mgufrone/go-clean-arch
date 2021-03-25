package worker

import (
	"context"
	"github.com/gocraft/work"
	"recipes/app/common"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type RecipeWorker interface {
	Logger(job *work.Job, next work.NextMiddlewareFunc)
	ProcessLike(job *work.Job) error
	ProcessView(job *work.Job) error
	ProcessDislike(job *work.Job) error
	BroadcastRecipe(job *work.Job) error
}

type recipeWorker struct {
	handler recipe.IRecipeUseCase
	ctx context.Context
}

func (r *recipeWorker) Logger(job *work.Job, next work.NextMiddlewareFunc) {
	if err := next(); err != nil {
	}
}

func (r *recipeWorker) ProcessLike(job *work.Job) error {
	rcp, usr := job.ArgInt64(common.RecipeKey), job.ArgInt64(common.UserKey)
	return r.handler.Like(r.ctx,
		&user.User{Model: shared.Model{ID: uint(usr)}},
		&recipe.Recipe{Model: shared.Model{ID: uint(rcp)}},
	)
}

func (r *recipeWorker) ProcessView(job *work.Job) error {
	rcp, usr := job.ArgInt64(common.RecipeKey), job.ArgInt64(common.UserKey)
	return r.handler.View(r.ctx,
		&user.User{Model: shared.Model{ID: uint(usr)}},
		&recipe.Recipe{Model: shared.Model{ID: uint(rcp)}},
	)
}

func (r *recipeWorker) ProcessDislike(job *work.Job) error {
	rcp, usr := job.ArgInt64(common.RecipeKey), job.ArgInt64(common.UserKey)
	return r.handler.Dislike(r.ctx,
		&user.User{Model: shared.Model{ID: uint(usr)}},
		&recipe.Recipe{Model: shared.Model{ID: uint(rcp)}},
	)
}

func (r *recipeWorker) BroadcastRecipe(job *work.Job) error {
	panic("implement me")
}

func NewRecipeWorker(ctx context.Context, handler recipe.IRecipeUseCase) RecipeWorker {
	return &recipeWorker{handler, ctx}
}
