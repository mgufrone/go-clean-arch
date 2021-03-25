package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/shared"
	"recipes/domains/stats"
)

type MockLikeRepository struct {
	mock.Mock
}

func (m *MockLikeRepository) CountByReference(ctx context.Context, reference string, id uint64) (int64, error) {
	args := m.Called(ctx, reference, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLikeRepository) Count(ctx context.Context, filters ...*stats.Like) (int64, error) {
	arg := []interface{}{ctx}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := m.Called(arg...)
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockLikeRepository) Put(ctx context.Context, like *stats.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func (m *MockLikeRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*stats.Like) ([]*stats.Like, *shared.CommonLimiter, error) {
	arg := []interface{}{ctx, limiter}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := m.Called(arg...)
	return args.Get(0).([]*stats.Like), args.Get(1).(*shared.CommonLimiter), args.Error(1)
}

func (m *MockLikeRepository) Delete(ctx context.Context, like *stats.Like) error {
	args := m.Called(ctx, like)
	return args.Error(0)
}

func NewMockLikeRepository() *MockLikeRepository {
	return &MockLikeRepository{}
}
