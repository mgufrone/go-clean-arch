package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*user.Role) ([]*user.Role, *shared.CommonLimiter, error) {
	panic("implement me")
}

func (m *MockRoleRepo) Create(ctx context.Context, role *user.Role) error {
	panic("implement me")
}

func (m *MockRoleRepo) Update(ctx context.Context, role *user.Role) error {
	panic("implement me")
}

func (m *MockRoleRepo) Delete(ctx context.Context, role *user.Role) error {
	panic("implement me")
}

func NewMockRoleRepository() *MockRoleRepo {
	return &MockRoleRepo{}
}
