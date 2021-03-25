package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) Create(ctx context.Context, session *user.Session) error {
	return nil
}

func (m *mockSessionRepo) Delete(ctx context.Context, session *user.Session) error {
	panic("implement me")
}

func (m *mockSessionRepo) Update(ctx context.Context, session *user.Session) error {
	panic("implement me")
}

func (m *mockSessionRepo) ListSession(ctx context.Context, limiter *shared.CommonLimiter, filters ...*user.Session) ([]*user.Session, *shared.CommonLimiter, error) {
	panic("implement me")
}

func NewMockSessionRepository() user.ISessionRepository {
	return &mockSessionRepo{}
}
