package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/shared"
	"recipes/domains/stats"
)

type MockViewRepository struct {
	mock.Mock
}

func (m *MockViewRepository) CountByReference(ctx context.Context, reference string, id uint64) (int64, error) {
	args := m.Called(ctx, reference, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockViewRepository) Count(ctx context.Context, filters ...*stats.View) (int64, error) {
	arg := []interface{}{ctx}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := m.Called(arg...)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockViewRepository) Put(ctx context.Context, view *stats.View) error {
	args := m.Called(ctx, view)
	return args.Error(0)
}

func (m *MockViewRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*stats.View) ([]*stats.View, *shared.CommonLimiter, error) {
	arg := []interface{}{ctx, limiter}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := m.Called(arg...)
	return args.Get(0).([]*stats.View), args.Get(1).(*shared.CommonLimiter), args.Error(1)
}

func (m *MockViewRepository) Delete(ctx context.Context, view *stats.View) error {
	args := m.Called(ctx, view)
	return args.Error(0)
}

func NewMockViewRepository() *MockViewRepository {
	return &MockViewRepository{}
}
