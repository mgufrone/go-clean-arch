package app

import (
	"context"
	firebase "firebase.google.com/go"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"os"
	"recipes/app/common"
	"recipes/app/config"
	"recipes/app/externals"
	"recipes/app/handlers"
	"recipes/app/repositories"
)

func fb() (*firebase.App, error)  {
	var opt = option.WithCredentialsFile(os.Getenv("FIREBASE_CONFIG"))
	FirebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	return FirebaseApp, err
}
func auth(app *firebase.App) (externals.IFirebaseAuth, error) {
	return app.Auth(context.Background())
}
func App(opts ...fx.Option) *fx.App {
	opts = append([]fx.Option{
		fx.Provide(
			fb,
			auth,
			config.Configure,
			common.InitializeDB,
			repositories.NewUserRepository,
			repositories.NewRecipeRepository,
			repositories.NewViewRepository,
			repositories.NewLikeRepository,
			repositories.NewSessionRepository,
			repositories.NewRoleRepository,
			handlers.NewRecipeHandler,
			handlers.NewUserHandler,
		),
	}, opts...)
	app := fx.New(
		opts...
	)
	return app
}
