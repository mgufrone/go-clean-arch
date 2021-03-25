package routes

import (
	"github.com/gin-gonic/gin"
	"recipes/api/controllers"
	"recipes/api/http"
	"recipes/api/requests"
)

type recipeHandler struct {
	handler *controllers.RecipeController
}

func (a *recipeHandler) Protected(rg *gin.RouterGroup) {
	rcp := rg.Group("/me/recipes")
	{
		rcp.GET("/", http.Sanitize(&requests.RecipeGet{}), a.handler.GetAllByUser)
	}
}

func (a *recipeHandler) Unprotected(rg *gin.RouterGroup)  {
	rcp := rg.Group("/recipes")
	{
		rcp.GET("/", http.Sanitize(&requests.RecipeGet{}), a.handler.GetAll)
		rcp.GET("/popular", http.Sanitize(&requests.RecipeGet{}), a.handler.Popular)
	}
}
func RecipeRoutes(handler *controllers.RecipeController) http.Routes {
	h := &recipeHandler{handler}
	return http.Routes{
		Unprotected: h.Unprotected,
		Protected:   h.Protected,
	}
}
