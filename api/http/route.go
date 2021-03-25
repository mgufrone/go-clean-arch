package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)
type RouteMount func(rg *gin.RouterGroup)

type Routes struct {
	fx.Out
	Unprotected RouteMount `group:"protected"`
	Protected RouteMount `group:"unprotected"`
}
type RegisteredRoutes struct {
	fx.In
	Unprotected []RouteMount `group:"protected"`
	Protected []RouteMount `group:"unprotected"`
}

