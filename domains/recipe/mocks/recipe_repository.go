package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/recipe"
	"recipes/domains/shared"
)

type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) Count(ctx context.Context, filters ...*recipe.Recipe) (int64, error) {
	arg := []interface{}{ctx}
	if len(filters) > 0 {
		for _, f := range filters {
			arg = append(arg, f)
		}
	}
	args := m.Called(arg...)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecipeRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*recipe.Recipe) ([]*recipe.Recipe, *shared.CommonLimiter, error) {
	arg := []interface{}{ctx, limiter}
	for _, r := range filters {
		arg = append(arg, r)
	}
	args := m.Called(arg...)
	return args.Get(0).([]*recipe.Recipe), args.Get(1).(*shared.CommonLimiter), args.Error(2)
}

func (m *MockRecipeRepository) Create(ctx context.Context, recipe *recipe.Recipe) error {
	args := m.Called(ctx, recipe)
	return args.Error(0)
}

func (m *MockRecipeRepository) Update(ctx context.Context, recipe *recipe.Recipe) error {
	args := m.Called(ctx, recipe)
	return args.Error(0)
}

func (m *MockRecipeRepository) Delete(ctx context.Context, recipe *recipe.Recipe) error {
	args := m.Called(ctx, recipe)
	return args.Error(0)
}

func (m *MockRecipeRepository) CreateBatch(ctx context.Context, recipes ...*recipe.Recipe) error {
	arg := []interface{}{ctx}
	for _, r := range recipes {
		arg = append(arg, r)
	}
	args := m.Called(arg...)
	return args.Error(0)
}

func (m *MockRecipeRepository) DeleteBatch(ctx context.Context, recipes ...*recipe.Recipe) error {
	arg := []interface{}{ctx}
	for _, r := range recipes {
		arg = append(arg, r)
	}
	args := m.Called(arg...)
	return args.Error(0)
}

func NewMockRecipeRepository() *MockRecipeRepository  {
	return &MockRecipeRepository{}
}
