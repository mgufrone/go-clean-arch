package user

import (
	"context"
	"recipes/domains/shared"
	"time"
)

type Session struct {
	shared.Model
	User *User
	SessionID string
	Agent string
	ExpiresAt time.Time
	RefreshExpiresAt time.Time
}

type ISessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	ListSession(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Session) ([]*Session, *shared.CommonLimiter, error)
}

type ISessionUseCase interface {
	CreateSession(ctx context.Context, usr *User, session *Session) error
	EndSession(ctx context.Context, usr *User, session *Session) error
	UpdateSession(ctx context.Context, usr *User, session *Session) error
	ListSession(ctx context.Context, limiter *shared.CommonLimiter, filters ...*Session) ([]*Session, *shared.CommonLimiter, error)
	IsSessionExpired(ctx context.Context, usr *User, session *Session) bool
	IsRefreshExpired(ctx context.Context, usr *User, session *Session) bool
}