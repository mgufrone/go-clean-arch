package repositories

import (
	"context"
	"gorm.io/gorm"
	"recipes/app/common"
	config2 "recipes/app/config"
	"recipes/app/models"
	"recipes/domains/shared"
	"recipes/domains/user"
)

type userRepo struct {
	common.GormDB
}

func (u *userRepo) Count(ctx context.Context, filters ...*user.User) (int64, error) {
	_, total, err := u.count(ctx, filters...)
	return total, err
}

func (u *userRepo) AddRole(ctx context.Context, user *user.User, role *user.Role) error {
	tx := u.DB.WithContext(ctx)
	rl := models.ParseRole(role)
	err := tx.Model(&models.UserModel{}).Where("id = ?", user.ID).Association("Roles").Append(rl)
	if err != nil {
		return err
	}
	user.Roles = append(user.Roles, rl.Transform())
	return nil
}

func (u *userRepo) RemoveRole(ctx context.Context, usr *user.User, role *user.Role) error {
	tx := u.DB.WithContext(ctx)
	err := tx.Model(usr).Association("Roles").Delete(role)
	if err != nil {
		return err
	}
	var rls []*user.Role
	for _, r := range usr.Roles {
		if r.ID != role.ID {
			rls = append(rls, r)
		}
	}
	usr.Roles = rls
	return nil
}

func (u *userRepo) count(ctx context.Context, filters ...*user.User) (tx *gorm.DB, total int64, err error) {
	tx = u.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}
func (u *userRepo) findAll(ctx context.Context, filters ...*user.User) (tx *gorm.DB) {
	tx = u.DB.WithContext(ctx).Model(&models.UserModel{})
	copied := u.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseUser(fltr)
		if mdl.FirstName != "" {
			ses = u.WhereLike(ses, "first_name", mdl.FirstName)
		}
		if mdl.LastName != "" {
			ses = u.WhereLike(ses, "last_name", mdl.LastName)
		}
		if mdl.Email != "" {
			ses = ses.Where("email = ?", mdl.Email)
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

func (u *userRepo) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*user.User) (res []*user.User, ltr *shared.CommonLimiter, err error) {
	tx, total, err := u.count(ctx, filters...)
	if err != nil {
		return
	}
	ltr = limiter
	ltr.SetTotal(uint64(total))
	if total == 0 {
		return
	}
	caps := uint64(total)
	if caps > ltr.GetPerPage() {
		caps = ltr.GetPerPage()
	}
	usrs := make([]*models.UserModel, 0, caps)
	u.GormDB.CommonFilter(tx, limiter)
	result := tx.Preload("Roles", "is_active = ?", true).Find(&usrs)
	if result.Error != nil {
		err = result.Error
		return
	}
	for _, usr := range usrs {
		res = append(res, usr.Transform())
	}
	return
}

func (u *userRepo) Create(ctx context.Context, usr *user.User) error {
	mdl := models.ParseUser(usr)
	err := u.GormDB.Create(ctx, mdl)
	if err != nil {
		return err
	}
	*usr = *mdl.Transform()
	return nil
}

func (u *userRepo) Update(ctx context.Context, usr *user.User) error {
	mdl := models.ParseUser(usr)
	err := u.GormDB.Update(ctx, mdl.ID, mdl)
	if err != nil {
		return err
	}
	*usr = *mdl.Transform()
	return nil
}

func (u *userRepo) Delete(ctx context.Context, usr *user.User) error {
	return u.GormDB.Delete(ctx, usr.ID.(uint))
}

func NewUserRepository(db *gorm.DB, cfg *config2.Config) user.IUserRepository {
	return &userRepo{common.GormDB{DB: db, Config: cfg}}
}