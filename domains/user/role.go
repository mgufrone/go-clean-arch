package user

import (
	"context"
	"recipes/domains/shared"
)

type Role struct {
	shared.Model
	Name string
	Users []*User
	IsActive bool
}

type IRoleRepository interface {
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Role) ([]*Role, *shared.CommonLimiter, error)
	Create(ctx context.Context, role *Role) error
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, role *Role) error
}
