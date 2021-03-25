package user

import (
	"context"
	"recipes/domains/shared"
)

type User struct {
	shared.Model
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	Email string `json:"email,omitempty"`
	Roles []*Role `json:"roles,omitempty"`
	IsActive bool `json:"is_active,omitempty"`
}

type IUserRepository interface {
	Count(ctx context.Context, filters ...*User) (int64, error)
	GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*User) ([]*User, *shared.CommonLimiter, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error

	AddRole(ctx context.Context, user *User, role *Role) error
	RemoveRole(ctx context.Context, user *User, role *Role) error
}

type IUserUseCase interface {
	FindByToken(ctx context.Context, token string) (*User, error)

	GetAll(ctx context.Context, actor *User, limiter *shared.CommonLimiter, filters ...*User) ([]*User, *shared.CommonLimiter, error)
	Create(ctx context.Context, actor *User, user *User) error
	Update(ctx context.Context, actor *User, user *User) error
	Delete(ctx context.Context, actor *User, user *User) error

	AddRole(ctx context.Context, actor *User, usr *User, role *Role) error
	RemoveRole(ctx context.Context, actor *User, usr *User, role *Role) error

	CheckSession(ctx context.Context, usr *User, session *Session) error
	Login(ctx context.Context, usr *User) (*Session, error)
	Refresh(ctx context.Context, sess *Session) error
	Logout(ctx context.Context, usr *User, session *Session) error
}
