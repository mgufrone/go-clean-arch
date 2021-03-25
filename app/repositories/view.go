package repositories

import (
	"context"
	"github.com/imdario/mergo"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/app/models"
	"recipes/domains/shared"
	"recipes/domains/stats"
)

type viewRepository struct {
	common.GormDB
}

func (v *viewRepository) GroupCount(ctx context.Context, limiter *shared.CommonLimiter, filters ...*stats.View) (res []*stats.ViewWithStats, ltr *shared.CommonLimiter, err error) {
	tx, total, _ := v.popular(ctx, filters...)
	ltr = limiter
	ltr.SetTotal(uint64(total))
	v.CommonFilter(tx, limiter)
	var rvs []*models.ViewModelWithStats
	tx.Debug().Find(&rvs)
	for _, rv := range rvs {
		v1 := rv.ViewModel.Transform()
		res = append(res, &stats.ViewWithStats{
			View:  v1,
			Count: rv.ViewCount,
		})
	}
	return
}

func (v *viewRepository) count(ctx context.Context, filters ...*stats.View) (tx *gorm.DB, total int64, err error) {
	tx = v.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}
func (v *viewRepository) popular(ctx context.Context, filters ...*stats.View) (tx *gorm.DB, total int64, err error) {
	tx = v.findAll(ctx, filters...).
		Select("view_models.reference", "view_models.reference_id", "COUNT(user_id) as view_count").
		Group("view_models.reference").
		Group("view_models.reference_id")
	tx.Debug().Count(&total)
	return
}
func (v *viewRepository) findAll(ctx context.Context, filters ...*stats.View) (tx *gorm.DB) {
	tx = v.DB.WithContext(ctx).Model(&models.ViewModel{})
	copied := v.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseView(fltr)
		if mdl.Reference != "" {
			ses = v.WhereLike(ses, "reference", mdl.Reference)
		}
		if mdl.ReferenceID > 0 {
			ses = ses.Where("reference_id like ?", mdl.ReferenceID)
		}
		if mdl.ID > 0 {
			ses = ses.Where("view_models.id = ?", mdl.ID)
		}
		if idx == 0 {
			tx.Where(ses)
		} else {
			tx.Or(ses)
		}
	}
	return tx
}
func (v *viewRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*stats.View) (res []*stats.View, ltr *shared.CommonLimiter, err error) {
	tx, total, err := v.count(ctx, filters...)
	if err != nil {
		return
	}
	limiter.SetTotal(uint64(total))
	if total == 0 {
		return
	}
	ltr = limiter
	v.CommonFilter(tx, limiter)
	var rvs []*models.ViewModel
	tx.Find(&rvs)
	for _, rv := range rvs {
		res = append(res, rv.Transform())
	}
	return
}

func (v *viewRepository) CountByReference(ctx context.Context, reference string, id uint64) (total int64, err error) {
	_, total, err = v.count(ctx, &stats.View{
		ReferenceID: id,
		Reference: reference,
	})
	return
}

func (v *viewRepository) Count(ctx context.Context, filters ...*stats.View) (total int64, err error) {
	_, total, err = v.count(ctx, filters...)
	return
}

func (v *viewRepository) Put(ctx context.Context, view *stats.View) (err error) {
	mdl := models.ParseView(view)
	err = v.GormDB.Create(ctx, mdl)
	if err != nil {
		return
	}
	err = mergo.Merge(view, mdl.Transform())
	return
}

func NewViewRepository(db *gorm.DB) stats.IViewRepository {
	return &viewRepository{GormDB: common.GormDB{
		DB: db,
	}}
}
