package routes

import (
	"github.com/gin-gonic/gin"
	"recipes/api/controllers"
	"recipes/api/http"
	"recipes/api/requests"
)

type authHandler struct {
	handler *controllers.AuthController
}

func (a *authHandler) Protected(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/logout", a.handler.Login)
	}
}

func (a *authHandler) Unprotected(rg *gin.RouterGroup)  {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", http.Sanitize(&requests.Login{}), a.handler.Login)
	}
}
func AuthRoutes(handler *controllers.AuthController) http.Routes {
	h := &authHandler{handler}
	return http.Routes{
		Unprotected: h.Unprotected,
		Protected:   h.Protected,
	}
}
