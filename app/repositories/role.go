package repositories

import (
	"context"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/app/models"
	shared2 "recipes/domains/shared"
	"recipes/domains/user"
)

type roleRepo struct {
	common.GormDB
}

func (r *roleRepo) findAll(ctx context.Context, filters ...*user.Role) (tx *gorm.DB) {
	tx = r.DB.WithContext(ctx).Model(&models.RoleModel{})
	copied := r.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseRole(fltr)
		if mdl.Name != "" {
			ses = ses.Where("name like ?", mdl.Name)
		}
		if mdl.IsActive != nil {
			ses = ses.Where("is_active = ?", *mdl.IsActive)
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
func (r *roleRepo) count(ctx context.Context, filters ...*user.Role) (tx *gorm.DB, total int64, err error) {
	tx = r.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}
func (r *roleRepo) GetAll(ctx context.Context, limiter *shared2.CommonLimiter, filters ...*user.Role) (res []*user.Role, ltr *shared2.CommonLimiter, err error) {
	tx, total, err := r.count(ctx, filters...)
	if err != nil {
		return
	}
	ltr = limiter
	ltr.SetTotal(uint64(total))
	if total == 0 {
		return
	}
	r.GormDB.CommonFilter(tx, limiter)
	caps := uint64(total)
	if caps > ltr.GetPerPage() {
		caps = ltr.GetPerPage()
	}
	rls := make([]*models.RoleModel, 0, caps)
	result := tx.Find(&rls)
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

func (r *roleRepo) Create(ctx context.Context, role *user.Role) error {
	rl := models.ParseRole(role)
	err := r.GormDB.Create(ctx, rl)
	if err != nil {
		return err
	}
	*role = *rl.Transform()
	return nil
}

func (r *roleRepo) Update(ctx context.Context, role *user.Role) error {
	rl := models.ParseRole(role)
	err := r.GormDB.Update(ctx, rl.ID, rl)
	if err != nil {
		return err
	}
	*role = *rl.Transform()
	return nil
}

func (r *roleRepo) Delete(ctx context.Context, role *user.Role) error {
	err := r.GormDB.Delete(ctx, role.ID.(uint))
	if err != nil {
		return err
	}
	role = nil
	return nil
}

func NewRoleRepository(db *gorm.DB) user.IRoleRepository {
	return &roleRepo{GormDB: common.GormDB{DB: db}}
}
