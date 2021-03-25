package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type UserRepository struct {
	mock.Mock
}

func (u *UserRepository) Count(ctx context.Context, filters ...*user.User) (int64, error) {
	arg := []interface{}{ctx}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := u.Called(arg...)
	return args.Get(0).(int64), args.Error(1)
}

func (u *UserRepository) AddRole(ctx context.Context, user *user.User, role *user.Role) error {
	args := u.Called(ctx, user, role)
	return args.Error(0)
}

func (u *UserRepository) RemoveRole(ctx context.Context, user *user.User, role *user.Role) error {
	args := u.Called(ctx, user, role)
	return args.Error(0)
}

func (u *UserRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*user.User) ([]*user.User, *shared.CommonLimiter, error) {
	arg := []interface{}{ctx, limiter}
	for _, f := range filters {
		arg = append(arg, f)
	}
	args := u.Called(arg...)
	return args.Get(0).([]*user.User), args.Get(1).(*shared.CommonLimiter), args.Error(2)
}

func (u *UserRepository) Create(ctx context.Context, user *user.User) error {
	args := u.Called(ctx, user)
	return args.Error(0)
}

func (u *UserRepository) Update(ctx context.Context, user *user.User) error {
	panic("implement me")
}

func (u *UserRepository) Delete(ctx context.Context, user *user.User) error {
	panic("implement me")
}

func NewUserMockRepository() *UserRepository  {
	return &UserRepository{}
}
