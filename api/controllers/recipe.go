package controllers

import (
	"github.com/gin-gonic/gin"
	"recipes/api/constants"
	"recipes/api/http"
	"recipes/api/requests"
	"recipes/app/common"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type RecipeController struct {
	handler recipe.IRecipeUseCase
}

func (r *RecipeController) Like(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) Dislike(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) Popular(ctx *gin.Context) {
	req, ok := ctx.Get(constants.RequestKey)
	var filters []*recipe.Recipe
	limiter := shared.DefaultLimiter(10)
	if ok {
		reqSan := req.(*requests.RecipeGet)
		ltr, fltr := reqSan.Transform()
		limiter = ltr
		if fltr != nil {
			filters = []*recipe.Recipe{fltr}
		}
	}
	res, l, err := r.handler.Popular(ctx, nil, limiter, filters...)
	if err != nil {
		http.ServerError(ctx, err)
	}
	http.OkWithTotal(ctx, res, uint(l.Offset.Total))
}

func (r *RecipeController) Create(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) Update(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) View(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) Delete(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) AddPhoto(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) UpdatePhoto(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) RemovePhoto(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) AddIngredient(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) UpdateIngredient(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) RemoveIngredient(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) AddStep(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) UpdateStep(ctx *gin.Context) {
	panic("implement me")
}

func (r *RecipeController) RemoveStep(ctx *gin.Context) {
	panic("implement me")
}

// RecipeList godoc
// @Summary Get available recipes
// @Description list of recipes
// @Security BearerToken
// @Tags recipe
// @Accept  json
// @Produce  json
// @Param q query requests.RecipeGet false "proposal payload"
// @Success 200 {object} http.DataWithTotalResponse{data=recipe.Recipe}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /recipes [get]
func (r *RecipeController) GetAll(ctx *gin.Context) {
	req, ok := ctx.Get(constants.RequestKey)
	limiter, filter := shared.DefaultLimiter(10), &recipe.Recipe{}
	if ok {
		reqSan, ok := req.(requests.RecipeGet)
		if ok {
			limiter, filter = reqSan.Transform()
		}
	}
	res, l, err := r.handler.GetAll(ctx, limiter, filter)
	if err != nil {
		http.ServerError(ctx, err)
	}
	http.OkWithTotal(ctx, res, uint(l.Offset.Total))
}
func (r *RecipeController) GetAllByUser(ctx *gin.Context) {
	req, ok := ctx.Get(constants.RequestKey)
	if !ok {
		http.BadRequest(ctx, nil)
		return
	}
	usr, ok := ctx.Get(common.UserKey)
	reqSan, ok := req.(requests.RecipeGet)
	if !ok {
		http.BadRequest(ctx, nil)
		return
	}
	limiter, filter := reqSan.Transform()
	filter.User = &user.User{
		Model:     shared.Model{
			ID: usr.(uint),
		},
	}
	res, l, err := r.handler.GetAll(ctx, limiter, filter)
	if err != nil {
		http.ServerError(ctx, err)
	}
	http.OkWithTotal(ctx, res, uint(l.Offset.Total))
}

func NewRecipeController(handler recipe.IRecipeUseCase) *RecipeController  {
	return &RecipeController{handler}
}

