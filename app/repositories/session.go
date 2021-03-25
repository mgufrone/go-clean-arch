package repositories

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/app/models"
	shared2 "recipes/domains/shared"
	"recipes/domains/user"
	"time"
)

type sessionRepo struct {
	common.GormDB
}

func (s *sessionRepo) Create(ctx context.Context, session *user.Session) error {
	sessID := fmt.Sprintf("%s:%d:%d", session.Agent, session.User.ID, time.Now().Unix())
	session.SessionID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(sessID)).String()
	mdl := models.ParseSession(session)
	err := s.GormDB.Create(ctx, mdl)
	if err != nil {
		return err
	}
	*session = *mdl.Transform()
	return nil
}

func (s *sessionRepo) Delete(ctx context.Context, session *user.Session) error {
	err := s.GormDB.Delete(ctx, session.ID.(uint))
	if err != nil {
		return err
	}
	session = nil
	return nil
}

func (s *sessionRepo) Update(ctx context.Context, session *user.Session) error {
	mdl := models.ParseSession(session)
	err := s.GormDB.Update(ctx, mdl.ID, mdl)
	if err != nil {
		return err
	}
	*session = *mdl.Transform()
	return nil
}

func (s *sessionRepo) count(ctx context.Context, filters ...*user.Session) (tx *gorm.DB, total int64, err error) {
	tx = s.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}
func (s *sessionRepo) findAll(ctx context.Context, filters ...*user.Session) (tx *gorm.DB) {
	tx = s.DB.WithContext(ctx).Model(&models.UserModel{})
	copied := s.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseSession(fltr)
		if mdl.UserID != 0 {
			ses = ses.Where("user_id = ?", mdl.UserID)
		}
		if mdl.SessionID != "" {
			ses = ses.Where("session_id = ?", mdl.SessionID)
		}
		if mdl.ID > 0 {
			ses = ses.Where("id = ?", mdl.ID)
		}
		if idx == 0 {
			tx.Where(ses)
		} else {
			tx.Or(ses)
		}
	}
	return tx
}
func (s *sessionRepo) ListSession(ctx context.Context, limiter *shared2.CommonLimiter, filters ...*user.Session) (res []*user.Session, ltr *shared2.CommonLimiter, err error) {
	tx, total, err := s.count(ctx, filters...)
	if err != nil {
		return
	}
	ltr = limiter
	ltr.SetTotal(uint64(total))
	if total == 0 {
		return
	}
	s.GormDB.CommonFilter(tx, limiter)
	caps := uint64(total)
	if caps > ltr.GetPerPage() {
		caps = ltr.GetPerPage()
	}
	result := make([]*models.SessionModel, 0, caps)
	tx.Find(&result)
	for _, r := range result {
		res = append(res, r.Transform())
	}
	return
}

func NewSessionRepository(db *gorm.DB) user.ISessionRepository {
	return &sessionRepo{GormDB: common.GormDB{DB: db}}
}
