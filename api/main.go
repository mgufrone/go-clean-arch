package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
	"net/http"
	"os"
	"recipes/api/controllers"
	http2 "recipes/api/http"
	"recipes/api/routes"
	"recipes/app"
	"recipes/app/config"
	"time"
)

// @title Recipe API
// @version 1.0
// @description This is a sample server celler server.

// @contact.name API Support
// @contact.url https://mgufron.com
// @contact.email mgufronefendi@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apiKey BearerToken
// @in header
// @name Authorization
func server(cfg *config.Config) *gin.Engine {
	ginMode := gin.ReleaseMode
	if cfg.AppEnv != "production" && cfg.AppEnv != "test" {
		ginMode = gin.DebugMode
	}
	if cfg.AppEnv == "test" {
		ginMode = gin.TestMode
	}
	gin.SetMode(ginMode)
	ginSrv := gin.Default()
	return ginSrv
}
func mountRouter(engine *gin.Engine, auth *controllers.AuthController, routes http2.RegisteredRoutes)  {
	rg := engine.Group("/")
	{
		for _, r := range routes.Unprotected {
			if r != nil {
				r(rg)
			}
		}
	}
	{
		rg.Use(auth.IsLoggedIn)
		for _, r := range routes.Protected {
			if r != nil {
				r(rg)
			}
		}
	}
}
func startServer(engine *gin.Engine, cfg *config.Config, lc fx.Lifecycle) {
	fmt.Println("assign server to port", cfg.AppPort)
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.AppPort),
		Handler: engine,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("start server at", cfg.AppPort)
			go func() {
				if err := srv.ListenAndServe(); err != nil {
					log.Println(err)
				}
			}()
			return nil
		},
		OnStop: srv.Shutdown,
	})
}
func main() {
	if os.Getenv("APP_ENV") == "" {
		godotenv.Load()
	}
	a := app.App(fx.Provide(
		controllers.NewAuthController,
		controllers.NewRecipeController,
		routes.AuthRoutes,
		routes.RecipeRoutes,
		server,
	), fx.Invoke(mountRouter, startServer))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
	defer cancel()
	if err := a.Start(ctx); err != nil {
		panic(err)
	}
	<- a.Done()
}

