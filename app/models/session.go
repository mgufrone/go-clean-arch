package models

import (
	"gorm.io/gorm"
	"recipes/domains/shared"
	"recipes/domains/user"
	"time"
)

type SessionModel struct {
	gorm.Model
	UserID uint64 `gorm:"index" json:"user_id"`
	SessionID string `gorm:"index" json:"session_id"`
	Agent string `gorm:"index" json:"agent"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	RefreshExpiresAt time.Time `gorm:"index" json:"refresh_expires_at"`
}

func (m *SessionModel) Transform() *user.Session {
	if m == nil {
		return nil
	}
	return &user.Session{
		Model:            shared.Model{},
		User:             &user.User{
			Model: shared.Model{
				ID: m.UserID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
		SessionID:        m.SessionID,
		Agent:            m.Agent,
		ExpiresAt:        m.ExpiresAt,
		RefreshExpiresAt: m.RefreshExpiresAt,
	}
}

func ParseSession(m *user.Session) *SessionModel {
	if m == nil {
		return  nil
	}
	return &SessionModel{
		Model:            gorm.Model{
			ID: m.ID.(uint),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		UserID:           m.User.ID.(uint64),
		SessionID:        m.SessionID,
		Agent:            m.Agent,
		ExpiresAt:        m.ExpiresAt,
		RefreshExpiresAt: m.RefreshExpiresAt,
	}
}
