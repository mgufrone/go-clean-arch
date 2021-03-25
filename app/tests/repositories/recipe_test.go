package repositories

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"recipes/app/models"
	"recipes/app/repositories"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/user"
	"strings"
	"syreclabs.com/go/faker"
	"testing"
)

type RecipeRepositorySuite struct {
	suite.Suite
	repository recipe.IRecipeRepository
	ctx        context.Context
	db *gorm.DB
}

func (u *RecipeRepositorySuite) SetupSuite() {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&allowNativePasswords=true&allowOldPasswords=true",
		"root",
		"root",
		"localhost",
		"3306",
		"testing",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields: true,
	})
	if err != nil {
		panic(err)
	}
	u.db = db.Debug()
	err = db.AutoMigrate(&models.RecipeModel{}, models.StepModel{}, models.IngredientModel{}, models.PhotoModel{})
	if err != nil {
		panic(err)
	}
	u.repository = repositories.NewRecipeRepository(db)
	u.ctx = context.Background()
}

func (u *RecipeRepositorySuite) TearDownSuite() {
	db, err := u.db.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
}

func (u *RecipeRepositorySuite) Test03GetAll() {
	r, l, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1))
	if assert.Nil(u.T(), err) {
		assert.NotNil(u.T(), l.Offset)
		assert.Equal(u.T(), uint64(1), l.Offset.Total)
		assert.Equal(u.T(), int(l.Offset.Total), len(r))
		assert.Greater(u.T(), int(r[0].ID.(uint)), 0)
	}
}
func (u *RecipeRepositorySuite) Test04GetAllWithFilter() {
	_, l, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1), &recipe.Recipe{
		Title:       "sauce",
		Summary:     "sauce",
		Description: "sauce",
	})
	if assert.Nil(u.T(), err) {
		assert.GreaterOrEqual(u.T(), uint64(0), l.Offset.Total)
		//assert.GreaterOrEqual(u.T(), uint64(0), len(res))
	}
	_, l, err = u.repository.GetAll(u.ctx, shared.DefaultLimiter(1), &recipe.Recipe{
		Ingredients: []*recipe.Ingredient{{
			Name: "onion",
			Type: "spice",
		}},
	})
	if assert.Nil(u.T(), err) {
		assert.GreaterOrEqual(u.T(), uint64(0), l.Offset.Total)
		//assert.GreaterOrEqual(u.T(), uint64(0), len(res))
	}
}
func (u *RecipeRepositorySuite) Test01Create() {
	u.db.Unscoped().Delete(&models.RecipeModel{}, "id != 0")
	u.db.Unscoped().Delete(&models.PhotoModel{}, "id != 0")
	u.db.Unscoped().Delete(&models.IngredientModel{}, "id != 0")
	u.db.Unscoped().Delete(&models.StepModel{}, "id != 0")
	p := &recipe.Recipe{
		Title:       faker.Lorem().Sentence(10),
		Summary:     faker.Lorem().Paragraph(10),
		Description: strings.Join(faker.Lorem().Paragraphs(10), "\n"),
		User: &user.User{
			Model: shared.Model{
				ID: uint64(1),
			},
		},
		Ingredients: []*recipe.Ingredient{{
			Sequence:    0,
			Type:        "spice",
			Weight:      "1/2",
			Measurement: "cup",
			Name:        "sugar",
		}},
		Photos: []*recipe.Photo{{
			Sequence: 0,
			URL:      faker.Internet().Url(),
		}},
		Steps: []*recipe.Step{{
			Sequence: 0,
			Step:      faker.Lorem().Sentence(10),
		}},
	}
	err := u.repository.Create(u.ctx, p)
	_, l, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(10))
	if err != nil {
		panic(err)
	}

	if assert.Nil(u.T(), err) {
		assert.Equal(u.T(), uint64(1), l.Offset.Total)
		assert.NotNil(u.T(), p.ID)
		assert.NotNil(u.T(), p.CreatedAt)
		assert.NotNil(u.T(), p.UpdatedAt)
		assert.NotNil(u.T(), p.Ingredients[0].ID)
		assert.NotNil(u.T(), p.Steps[0].ID)
		assert.NotNil(u.T(), p.Photos[0].ID)
	}
}
func (u *RecipeRepositorySuite) Test02Update()  {
	res, _, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1))
	assert.Nil(u.T(), err)
	p := &recipe.Recipe{}
	p.ID = res[0].ID
	p.Summary = "change it to another thing"
	err = u.repository.Update(u.ctx, p)
	check := &recipe.Recipe{}
	check.ID = res[0].ID

	if assert.Nil(u.T(), err) {
		res2, l, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1), check)
		if err != nil {
			panic(err)
		}
		assert.Equal(u.T(), uint64(1), l.Offset.Total)
		assert.NotNil(u.T(), p.ID)
		assert.Equal(u.T(), res2[0].Title, res[0].Title)
		assert.Equal(u.T(), res2[0].Summary, p.Summary)
		assert.Equal(u.T(), res2[0].Title, p.Title)
	}
}
func (u *RecipeRepositorySuite) Test05Delete() {
	res, _, err := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1))
	assert.Nil(u.T(), err)
	p := &recipe.Recipe{}
	p.ID = res[0].ID
	err = u.repository.Delete(u.ctx, p)
	if assert.Nil(u.T(), err) {
		_, l, _ := u.repository.GetAll(u.ctx, shared.DefaultLimiter(1))
		assert.Nil(u.T(), p.ID)
		assert.NotNil(u.T(), l)
		assert.Equal(u.T(), uint64(0), l.Offset.Total)
	}
}

func TestRecipeRepository(t *testing.T) {
	suite.Run(t, new(RecipeRepositorySuite))
}
