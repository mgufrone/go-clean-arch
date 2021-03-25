package repositories

import (
	"context"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/app/models"
	"recipes/domains/shared"
	"recipes/domains/stats"
)

type likeRepository struct {
	common.GormDB
}

func (v *likeRepository) count(ctx context.Context, filters ...*stats.Like) (tx *gorm.DB, total int64, err error) {
	tx = v.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}
func (v *likeRepository) popular(ctx context.Context, filters ...*stats.Like) (tx *gorm.DB, total int64, err error) {
	tx = v.findAll(ctx, filters...).
		Select("like_models.reference", "like_models.reference_id", "COUNT(user_id) as popular").
		Group("like_models.reference").
		Group("like_models.reference_id")
	tx.Count(&total)
	return
}
func (v *likeRepository) findAll(ctx context.Context, filters ...*stats.Like) (tx *gorm.DB) {
	tx = v.DB.WithContext(ctx).Model(&models.LikeModel{})
	copied := v.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseLike(fltr)
		if mdl.Reference != "" {
			ses = v.WhereLike(ses, "reference", mdl.Reference)
		}
		if mdl.ReferenceID > 0 {
			ses = ses.Where( "reference_id like ?", mdl.ReferenceID)
		}
		if mdl.ID > 0 {
			ses = ses.Where("like_models.id = ?", mdl.ID)
		}
		if idx == 0 {
			tx.Where(ses)
		} else {
			tx.Or(ses)
		}
	}
	return tx
}
func (v *likeRepository) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*stats.Like) (res []*stats.Like, ltr *shared.CommonLimiter, err error) {
	var (
		tx *gorm.DB
		total int64
	)
	if limiter.Sort != nil && limiter.Sort.Field == "popular" {
		tx, total, err = v.popular(ctx, filters...)
	} else {
		tx, total, err = v.count(ctx, filters...)
	}
	if err != nil {
		return
	}
	limiter.SetTotal(uint64(total))
	if total == 0 {
		return
	}
	ltr = limiter
	v.CommonFilter(tx, limiter)
	caps := common.SliceMin(limiter.Total(), limiter.GetPerPage())
	rvs := make([]*models.LikeModel, 0, caps)
	tx.Find(&rvs)
	for _, rv := range rvs {
		res = append(res, rv.Transform())
	}
	return
}

func (v *likeRepository) CountByReference(ctx context.Context, reference string, id uint64) (total int64, err error) {
	_, total, err = v.count(ctx, &stats.Like{
		ReferenceID: id,
		Reference: reference,
	})
	return
}

func (v *likeRepository) Count(ctx context.Context, filters ...*stats.Like) (total int64, err error) {
	_, total, err = v.count(ctx, filters...)
	return
}

func (v *likeRepository) Put(ctx context.Context, like *stats.Like) (err error) {
	mdl := models.ParseLike(like)
	err = v.GormDB.Create(ctx, mdl)
	if err != nil {
		return
	}
	*like = *mdl.Transform()
	return
}
func (v *likeRepository) Delete(ctx context.Context, like *stats.Like) (err error) {
	mdl := models.ParseLike(like)
	res := v.DB.WithContext(ctx).
		Model(&models.LikeModel{}).
		Delete("reference = ? and reference_id = ?", like.Reference, like.ReferenceID)
	if res.Error != nil {
		err = res.Error
		return
	}
	*like = *mdl.Transform()
	return
}

func NewLikeRepository(db *gorm.DB) stats.ILikeRepository {
	return &likeRepository{GormDB: common.GormDB{
		DB: db,
	}}
}
