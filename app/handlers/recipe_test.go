package handlers

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"recipes/domains/recipe"
	mocks2 "recipes/domains/recipe/mocks"
	"recipes/domains/shared"
	mocks3 "recipes/domains/stats/mocks"
	"recipes/domains/user"
	"recipes/domains/user/mocks"
	"strings"
	"syreclabs.com/go/faker"
	"testing"
)

type RecipeHandlerSuite struct {
	suite.Suite
	userRepo *mocks.UserRepository
	recipe   *mocks2.MockRecipeRepository
	like  *mocks3.MockLikeRepository
	view *mocks3.MockViewRepository

	handler       recipe.IRecipeUseCase
	sharedContext context.Context
}

func (u *RecipeHandlerSuite) SetupSuite() {
	u.sharedContext = context.Background()
	u.userRepo = mocks.NewUserMockRepository()
	u.recipe = mocks2.NewMockRecipeRepository()
	u.view = mocks3.NewMockViewRepository()
	u.like = mocks3.NewMockLikeRepository()
	u.handler = NewRecipeHandler(
		u.recipe,
		u.userRepo,
		u.like,
		u.view,
	)
}
func (u *RecipeHandlerSuite) BeforeTest(suite string, testName string) {
	u.recipe.Calls = []mock.Call{}
	u.recipe.ExpectedCalls = []*mock.Call{}
	u.userRepo.Calls = []mock.Call{}
	u.userRepo.ExpectedCalls = []*mock.Call{}
	u.like.Calls = []mock.Call{}
	u.like.ExpectedCalls = []*mock.Call{}
	u.view.Calls = []mock.Call{}
	u.view.ExpectedCalls = []*mock.Call{}
}

func (u *RecipeHandlerSuite) TearDownSuite() {
}

func (u *RecipeHandlerSuite) TestGetAllEmpty() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 0
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil).
		Times(1)
	qLimit := shared.DefaultLimiter(1)

	res, l, err := u.handler.GetAll(u.sharedContext, qLimit)
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Len(u.T(), res, 0)
		assert.Len(u.T(), u.recipe.Calls[0].Arguments, 2)
		assert.Equal(u.T(), u.recipe.Calls[0].Arguments[1], qLimit)
		assert.Equal(u.T(), l.Offset.Total, uint64(0))
	}
}
func (u *RecipeHandlerSuite) TestGetAllIgnore() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 0
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil).
		Times(1)
	qLimit := shared.DefaultLimiter(1)
	qLimit.Offset.PerPage = 0

	res, _, err := u.handler.GetAll(u.sharedContext, qLimit)
	if assert.NotNil(u.T(), err) {
		assert.Len(u.T(), res, 0)
		assert.Len(u.T(), u.recipe.Calls, 0)
	}
}
func (u *RecipeHandlerSuite) TestGetAllFilter() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 5
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil)

	_, l, err := u.handler.GetAll(u.sharedContext, shared.DefaultLimiter(5), &recipe.Recipe{
		Title: "chi",
	})
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Len(u.T(), u.recipe.Calls, 1)
		assert.Equal(u.T(), u.recipe.Calls[0].Arguments[2], &recipe.Recipe{
			Title: "chi",
		})
		assert.Equal(u.T(), l.Offset.Total, uint64(5))
	}
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil)

	_, _, err = u.handler.GetAll(u.sharedContext, shared.DefaultLimiter(5), &recipe.Recipe{
		Title: "chi",
	}, &recipe.Recipe{
		Summary: "nug",
	})
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Equal(u.T(), u.recipe.Calls[1].Arguments[2], &recipe.Recipe{
			Title: "chi",
		})
		assert.Equal(u.T(), u.recipe.Calls[1].Arguments[3], &recipe.Recipe{
			Summary: "nug",
		})
		assert.Equal(u.T(), l.Offset.Total, uint64(5))
	}
}
func (u *RecipeHandlerSuite) TestGetByUser() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 0
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil).
		Times(1)
	qLimit := shared.DefaultLimiter(1)

	res, _, err := u.handler.GetByUser(u.sharedContext, qLimit, &user.User{
		Model: shared.Model{
			ID: 1,
		},
	})
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Len(u.T(), res, 0)
		assert.Len(u.T(), u.recipe.Calls, 1)
	}
}

func (u *RecipeHandlerSuite) TestGetAllDeep() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 5
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil)

	_, l, err := u.handler.GetAll(u.sharedContext, shared.DefaultLimiter(5), &recipe.Recipe{
		Title: "onion",
		Ingredients: []*recipe.Ingredient{{
			Type: "spice",
		}},
	})
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Len(u.T(), u.recipe.Calls, 1)
		assert.Equal(u.T(), u.recipe.Calls[0].Arguments[2], &recipe.Recipe{
			Title: "onion",
			Ingredients: []*recipe.Ingredient{{
				Type: "spice",
			}},
		})
		assert.Equal(u.T(), l.Offset.Total, uint64(5))
	}
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, nil)

	_, _, err = u.handler.GetAll(u.sharedContext, shared.DefaultLimiter(5), &recipe.Recipe{
		Title: "chi",
	}, &recipe.Recipe{
		Summary: "nug",
	})
	if assert.Nil(u.T(), err) {
		u.recipe.AssertExpectations(u.T())
		assert.Equal(u.T(), u.recipe.Calls[1].Arguments[2], &recipe.Recipe{
			Title: "chi",
		})
		assert.Equal(u.T(), u.recipe.Calls[1].Arguments[3], &recipe.Recipe{
			Summary: "nug",
		})
		assert.Equal(u.T(), l.Offset.Total, uint64(5))
	}
}
func (u *RecipeHandlerSuite) TestCreateRecipeSuccess() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = uint64(1)
	u.userRepo.On("Count", mock.Anything, mock.Anything).
		Return(int64(1), nil)
	mockLimiter2 := shared.DefaultLimiter(1)
	mockLimiter2.Offset.Total = 0
	u.recipe.On("Count", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)
	u.recipe.On("Create", mock.Anything, mock.Anything).Return(nil)
	rcp := &recipe.Recipe{
		Title:       faker.RandomString(100),
		Summary:     faker.RandomString(200),
		Description: faker.RandomString(500),
		Steps: []*recipe.Step{
			{Sequence: 1, Step: "What are you doing step-recipe?"},
		},
		Ingredients: []*recipe.Ingredient{
			{
				Measurement: "tbsp",
				Weight: "1",
				Sequence: 1,
				Name: "Onion",
				Type: "spices",
			},
		},
		Photos: []*recipe.Photo{
			{
				Sequence: 1,
				URL:      "http://somedomain.com/somefolder/someimage.png",
			},
		},
	}
	usr := &user.User{
		Model: shared.Model{
			ID: uint64(1),
		},
	}
	err := u.handler.Create(u.sharedContext, usr, rcp)
	if assert.Nil(u.T(), err) {
		rcp.User = usr
		u.recipe.AssertExpectations(u.T())
		u.userRepo.AssertExpectations(u.T())
		assert.Len(u.T(), u.userRepo.Calls, 1)
		assert.Len(u.T(), u.recipe.Calls, 2)
		assert.Equal(u.T(), u.recipe.Calls[1].Arguments[1], rcp)
	}
}
func (u *RecipeHandlerSuite) TestCreateRecipeInvalidUser() {
	u.recipe.On("Count", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)
	u.recipe.On("Create", mock.Anything, mock.Anything).Return(nil)
	rcp := &recipe.Recipe{
		Title:       faker.RandomString(100),
		Summary:     faker.RandomString(200),
		Description: faker.RandomString(500),
		Steps: []*recipe.Step{
			{Sequence: 1, Step: "What are you doing step-recipe?"},
		},
		Ingredients: []*recipe.Ingredient{
			{
				Measurement: "tbsp",
				Weight: "1",
				Sequence: 1,
				Name: "Onion",
				Type: "spices",
			},
		},
		Photos: []*recipe.Photo{
			{
				Sequence: 1,
				URL:      "http://somedomain.com/somefolder/someimage.png",
			},
		},
	}
	usr := &user.User{
		Model: shared.Model{
			ID: nil,
		},
	}
	err := u.handler.Create(u.sharedContext, usr, rcp)
	assert.NotNil(u.T(), err)
	assert.Error(u.T(), err, "invalid user")
	usr = &user.User{
		Model: shared.Model{
			ID: uint64(0),
		},
	}
	err = u.handler.Create(u.sharedContext, usr, rcp)
	assert.NotNil(u.T(), err)
	assert.Error(u.T(), err, "invalid user")
}
func (u *RecipeHandlerSuite) TestCreateRecipeInvalidUser02() {
	u.userRepo.On("Count", mock.Anything, mock.Anything).
		Return(int64(0), nil)
	u.recipe.On("Count", mock.Anything, mock.Anything, mock.Anything).Return(int64(0), nil)
	u.recipe.On("Create", mock.Anything, mock.Anything).Return(nil)
	rcp := &recipe.Recipe{
		Title:       faker.RandomString(100),
		Summary:     faker.RandomString(200),
		Description: faker.RandomString(500),
		Steps: []*recipe.Step{
			{Sequence: 1, Step: "What are you doing step-recipe?"},
		},
		Ingredients: []*recipe.Ingredient{
			{
				Measurement: "tbsp",
				Weight: "1",
				Sequence: 1,
				Name: "Onion",
				Type: "spices",
			},
		},
		Photos: []*recipe.Photo{
			{
				Sequence: 1,
				URL:      "http://somedomain.com/somefolder/someimage.png",
			},
		},
	}
	usr := &user.User{
		Model: shared.Model{
			ID: uint64(1),
		},
	}
	err := u.handler.Create(u.sharedContext, usr, rcp)
	assert.NotNil(u.T(), err)
	assert.Error(u.T(), err, "user not found")
}
func (u *RecipeHandlerSuite) TestRecipeInvalid() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 1
	u.userRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return([]*user.User{{
			Model: shared.Model{ID: 1},
		}}, mockLimiter, nil)
	mockLimiter2 := shared.DefaultLimiter(1)
	mockLimiter2.Offset.Total = 0
	u.recipe.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return([]*recipe.Recipe{}, mockLimiter2, nil)
	u.recipe.On("Create", mock.Anything, mock.Anything).Return(nil)
	rcp := &recipe.Recipe{
		Summary:     faker.RandomString(100),
		Description: faker.RandomString(500),
		Steps: []*recipe.Step{
			{Sequence: 1, Step: "What are you doing step-recipe?"},
		},
		Ingredients: []*recipe.Ingredient{
			{
				Measurement: "tbsp",
				Weight: "1",
				Sequence: 1,
				Name: "Onion",
				Type: "spices",
			},
		},
		Photos: []*recipe.Photo{
			{
				Sequence: 1,
				URL:      "http://somedomain.com/somefolder/someimage.png",
			},
		},
	}
	usr := &user.User{
		Model: shared.Model{
			ID: uint64(1),
		},
	}
	err := u.handler.Create(u.sharedContext, usr, rcp)
	if assert.NotNil(u.T(), err) {
		assert.Len(u.T(), u.userRepo.Calls, 0)
		assert.Len(u.T(), u.recipe.Calls, 0)
	}
}
func (u *RecipeHandlerSuite) TestUpdateRecipeSuccess() {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 1
	u.view.On("CountByReference", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(0), nil)
	u.like.On("CountByReference", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(0), nil)
	u.userRepo.On("Count", mock.Anything, mock.Anything, mock.Anything).
		Return(int64(1), nil)
	u.userRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return([]*user.User{{
			Model: shared.Model{ID: uint64(1)},
		}}, mockLimiter, nil)
	mockLimiter2 := shared.DefaultLimiter(1)
	mockLimiter2.Offset.Total = 1
	exRecp := &recipe.Recipe{
		Model: shared.Model{ID: uint(1)},
		Summary:     faker.Lorem().Paragraph(5),
		Description: strings.Join(faker.Lorem().Paragraphs(10), "\n"),
		Steps: []*recipe.Step{
			{Sequence: 1, Step: "What are you doing step-recipe?"},
		},
		User: &user.User{
			Model: shared.Model{
				ID: uint64(1),
			},
		},
		Ingredients: []*recipe.Ingredient{
			{
				Measurement: "tbsp",
				Weight: "1",
				Sequence: 1,
				Name: "Onion",
				Type: "spices",
			},
		},
		Photos: []*recipe.Photo{
			{
				Sequence: 1,
				URL:      "http://somedomain.com/somefolder/someimage.png",
			},
		},
	}
	u.recipe.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return([]*recipe.Recipe{exRecp}, mockLimiter2, nil)
	u.recipe.On("Update", mock.Anything, mock.Anything).Return(nil)
	rcp := &recipe.Recipe{
		Model: shared.Model{ID: uint(1)},
		Summary:     faker.Lorem().Paragraph(5),
	}
	usr := &user.User{
		Model: shared.Model{
			ID: uint64(1),
		},
	}
	err := u.handler.Update(u.sharedContext, usr, rcp)
	if assert.Nil(u.T(), err) {
		u.userRepo.AssertExpectations(u.T())
		u.recipe.AssertExpectations(u.T())
		assert.Len(u.T(), u.userRepo.Calls, 2)
		assert.Len(u.T(), u.recipe.Calls, 2)
	}
}
func (u *RecipeHandlerSuite) TestGetAllError()  {
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 0
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{}, mockLimiter, errors.New("server error"))
	u.userRepo.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*user.User{{
			Model: shared.Model{ID: uint64(1)},
	}}, mockLimiter, nil)
	qLimit := shared.DefaultLimiter(1)

	_, _, err := u.handler.GetAll(u.sharedContext, qLimit)
	assert.NotNil(u.T(), err)
	// error on count view
	mockLimiter = shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = 0
	rcp := &recipe.Recipe{
		Model: shared.Model{
			ID: uint(1),
		},
		User: &user.User{
			Model:     shared.Model{
				ID: uint64(1),
			},
		},
	}
	u.recipe.
		On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return([]*recipe.Recipe{rcp}, mockLimiter, nil).
		Times(1)
	u.view.On("CountByReference", mock.Anything, "recipes", uint64(1)).
		Return(int64(0), errors.New("something went wrong"))
	u.like.On("CountByReference", mock.Anything, "recipes", uint64(1)).
		Return(int64(0), errors.New("something went wrong"))
	qLimit = shared.DefaultLimiter(1)
	_, _, err = u.handler.GetAll(u.sharedContext, qLimit)
	assert.Nil(u.T(), err)
}
func (u *RecipeHandlerSuite) TestLikeSuccess() {
	usr := &user.User{
		Model:     shared.Model{
			ID: uint64(1),
		},
	}
	rcp := &recipe.Recipe{
		Model: shared.Model{
			ID: uint(1),
		},
	}
	u.recipe.On("Count", u.sharedContext, rcp).
		Return(int64(1), nil)
	u.userRepo.On("Count", u.sharedContext, usr).
		Return(int64(1), nil)
	u.like.On("Put", u.sharedContext, mock.AnythingOfType("*stats.Like")).
		Return(nil)
	err := u.handler.Like(u.sharedContext, usr, rcp)
	if assert.Nil(u.T(), err) {
		u.like.AssertExpectations(u.T())
		u.userRepo.AssertExpectations(u.T())
		u.recipe.AssertExpectations(u.T())
	}
}

func TestRecipeHandler(t *testing.T) {
	suite.Run(t, new(RecipeHandlerSuite))
}
