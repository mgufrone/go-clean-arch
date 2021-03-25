package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"recipes/app"
	"recipes/app/common"
	"recipes/app/models"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/user"
	"strings"
	"syreclabs.com/go/faker"
	"time"
)

func migrator(db *gorm.DB) error {
	return db.AutoMigrate(&models.UserModel{},
		&models.IngredientModel{},
		&models.UserModel{},
		&models.RecipeModel{},
		&models.StepModel{},
		&models.PhotoModel{},
		&models.ViewModel{},
		&models.LikeModel{},
		&models.RoleModel{},
	)
}

func generateIngredients(max int) []*recipe.Ingredient {
	caps := rand.Intn(max) + 1
	res := make([]*recipe.Ingredient, 0, caps)
	for i := 0; i < caps; i += 1 {
		res = append(res, &recipe.Ingredient{
			Sequence:    i,
			Type:        faker.RandomChoice([]string{"spices", "main", "additional", "color", "flavor"}),
			Measurement:      faker.RandomChoice([]string{"g", "mg", "ml", "l", "cc", "oz"}),
			Weight: faker.RandomChoice([]string{"1/2", "1/4", "1", "2", "4", "50"}),
			Name:        faker.Lorem().Word(),
		})
	}
	return res
}
func generatePhotos(max int) []*recipe.Photo {
	caps := rand.Intn(max) + 1
	res := make([]*recipe.Photo, 0, caps)
	for i := 0; i < caps; i += 1 {
		res = append(res, &recipe.Photo{
			Sequence:    uint(i),
			URL: faker.Internet().Url(),
		})
	}
	return res
}
func generateSteps(max int) []*recipe.Step {
	caps := rand.Intn(max) + 1
	res := make([]*recipe.Step, 0, caps)
	for i := 0; i < caps; i += 1 {
		res = append(res, &recipe.Step{
			Sequence:    uint(i),
			Step: faker.Lorem().Sentence(10),
		})
	}
	return res
}

func seeder(db *gorm.DB, role user.IRoleRepository, h1 user.IUserUseCase, h2 recipe.IRecipeUseCase) {
	const (
		MaxUser int = 100
		MaxRecipes int = 1000
		MaxPhotoPerRecipe int = 10
		MaxIngredientPerRecipe int = 10
		MaxStepPerRecipe int = 20
	)
	db.Exec("TRUNCATE user_roles")
	db.Exec("TRUNCATE role_models")
	db.Exec("TRUNCATE user_models")
	db.Exec("TRUNCATE recipe_models")
	db.Exec("TRUNCATE ingredient_models")
	db.Exec("TRUNCATE photo_models")
	db.Exec("TRUNCATE step_models")
	ctx := context.Background()
	admin, member := &user.Role{IsActive: true, Name: "admin"}, &user.Role{IsActive: true, Name: "member"}
	usrs := make([]*user.User, 0, MaxUser)
	rcps := make([]*recipe.Recipe, 0, MaxRecipes)
	err := common.Try(func() error {
		return role.Create(ctx, member)
	}, func() error {
		return role.Create(ctx, admin)
	}, func() error {
		for i := 0; i < MaxUser; i++ {
			rRole := rand.Intn(1)
			r := member
			if rRole == 0 {
				r = admin
			}
			usr := &user.User{
				LastName:  faker.Name().LastName(),
				FirstName: faker.Name().FirstName(),
				Email: faker.Internet().Email(),
				Roles: []*user.Role{
					r,
				},
			}
			if err := h1.Create(ctx, nil, usr); err != nil {
				return err
			}
			usrs = append(usrs, usr)
		}
		return nil
	}, func() error {
		for i := 0; i < MaxRecipes; i++ {
			randUser := rand.Intn(MaxUser)
			u := usrs[randUser]
			rcp := &recipe.Recipe{
				Title:       faker.Lorem().Sentence(10),
				Summary:     faker.Lorem().Sentence(15),
				Description: strings.Join(faker.Lorem().Paragraphs(5), "\n"),
				Photos:      generatePhotos(MaxPhotoPerRecipe),
				Steps:       generateSteps(MaxStepPerRecipe),
				Ingredients: generateIngredients(MaxIngredientPerRecipe),
			}
			if err := h2.Create(ctx, u, rcp); err != nil {
				return err
			}
			rcps = append(rcps, rcp)
		}
		return nil
	}, func() error {
		// push random actions. (view, like)
		for j := 0; j < MaxUser * MaxRecipes; j ++ {
			act := rand.Intn(1)
			randUser := rand.Intn(len(usrs) - 1)
			rIdx := rand.Intn(len(rcps) - 1)
			r := rcps[rIdx]
			usr := usrs[randUser]
			fmt.Println("is liking", act)
			if act == 1 {
				if err := h2.Like(ctx, usr, r); err != nil {
					log.Println("like action failed", err)
				}
			} else {
				if err := h2.View(ctx, usr, r); err != nil {
					log.Println("view action failed", err)
				}
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func formulate(h1 recipe.IRecipeUseCase, r1 user.IUserRepository) error {
	start := time.Now()
	ctx := context.Background()
	usr, _, _ := r1.GetAll(ctx, shared.DefaultLimiter(1))
	pops, _, err := h1.Popular(ctx, usr[0], shared.DefaultLimiter(10))
	for _, pop := range pops {
		fmt.Println(pop.ID, pop.ViewCount, pop.LikeCount)
	}
	end := time.Since(start)
	fmt.Println("execution time", end)
	return err
}

func main() {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	flag.Parse()
	a := app.App(fx.Invoke(
		migrator,
		formulate))
	ctx := context.Background()
	err := a.Start(ctx)
	if err != nil {
		panic(err)
	}
	a.Stop(ctx)
}
