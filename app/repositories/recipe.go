package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"recipes/app/common"
	"recipes/app/models"
	"recipes/domains/recipe"
	"recipes/domains/shared"
)

type recipeRepo struct {
	common.GormDB
}

func (r *recipeRepo) Count(ctx context.Context, filters ...*recipe.Recipe) (total int64, err error) {
	_, total, err = r.count(ctx, filters...)
	return
}

func (r *recipeRepo) count(ctx context.Context, filters ...*recipe.Recipe) (tx *gorm.DB, total int64, err error) {
	tx = r.findAll(ctx, filters...)
	err = tx.Count(&total).Error
	return
}

func (r *recipeRepo) filterByIngredients(tx *gorm.DB, ingredients []*models.IngredientModel) *gorm.DB {
	ses := tx.WithContext(context.Background())
	for _, i := range ingredients {
		if i.Type != "" {
			ses = r.WhereLike(ses, "ingredient_models.type", i.Type)
		}
		if i.Name != "" {
			ses = r.WhereLike(ses, "ingredient_models.name", i.Name)
		}
	}
	return ses
}

func (r *recipeRepo) findAll(ctx context.Context, filters ...*recipe.Recipe) *gorm.DB {
	tx := r.DB.WithContext(ctx).Model(&models.RecipeModel{})
	copied := r.DB.WithContext(ctx)
	for idx, fltr := range filters {
		ses := copied.WithContext(context.Background())
		mdl := models.ParseRecipe(fltr)
		if mdl.Title != "" {
			ses = r.WhereLike(ses, "title", mdl.Title)
		}
		if mdl.Summary != "" {
			ses = r.WhereLike(ses, "summary", mdl.Summary)
		}
		if mdl.Description != "" {
			ses = ses.Where("description", mdl.Description)
		}
		if mdl.UserID > 0 {
			ses = ses.Where("user_id", mdl.UserID)
		}
		if mdl.ID > 0 {
			ses = ses.Where("recipe_models.id", mdl.ID)
		}

		if mdl.Ingredients != nil && len(mdl.Ingredients) > 0 {
			tx = tx.
				Joins("inner join ingredient_models on ingredient_models.recipe_id = recipe_models.id")
			ses = r.filterByIngredients(ses, mdl.Ingredients)
		}
		if idx == 0 {
			tx.Where(ses)
		} else {
			tx.Or(ses)
		}
	}
	return tx
}

func (r *recipeRepo) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*recipe.Recipe) (res []*recipe.Recipe, ltr *shared.CommonLimiter, err error) {
	var rcps []*models.RecipeModel
	tx, total, err := r.count(ctx, filters...)
	if err != nil {
		return
	}
	if limiter.Offset != nil {
		limiter.Offset.Total = uint64(total)
	}
	if limiter.Cursor != nil {
		limiter.Cursor.Total = uint64(total)
	}
	ltr = limiter
	if total == 0 {
		return
	}
	r.GormDB.CommonFilter(tx, limiter)
	result := tx.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredient_models.sequence asc")
	}).Preload("Photos").Preload("Steps").Find(&rcps)
	err = result.Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	for _, rcp := range rcps {
		res = append(res, rcp.Transform())
	}
	return
}

func (r *recipeRepo) Create(ctx context.Context, rcpRes *recipe.Recipe) (err error) {
	rcp := models.ParseRecipe(rcpRes)
	err = r.GormDB.Create(ctx, &rcp)
	if err != nil {
		return
	}
	*rcpRes = *rcp.Transform()
	return
}

func (r *recipeRepo) Update(ctx context.Context, resRcp *recipe.Recipe) (err error) {
	rcp := models.ParseRecipe(resRcp)
	tx := r.DB.WithContext(ctx).Session(&gorm.Session{
		FullSaveAssociations: false,
	}).Model(&models.RecipeModel{})
	cur := tx.Where("id = ?", rcp.ID)
	tx.Model(&models.IngredientModel{}).Unscoped().Delete("recipe_id = ?", rcp.ID)
	tx.Model(&models.PhotoModel{}).Unscoped().Delete("recipe_id = ?", rcp.ID)
	tx.Model(&models.StepModel{}).Unscoped().Delete("recipe_id = ?", rcp.ID)
	res := cur.Updates(rcp)
	if res.Error != nil {
		err = res.Error
		return
	}
	check := &models.RecipeModel{}
	check.ID = rcp.ID
	ex, _, err := r.GetAll(ctx, shared.DefaultLimiter(1), check.Transform())
	if err != nil {
		return
	}
	*resRcp = *ex[0]
	return
}

func (r *recipeRepo) Delete(ctx context.Context, res *recipe.Recipe) error {
	rcp := models.ParseRecipe(res)
	o := r.DB.WithContext(ctx).Delete(&rcp)
	if o.Error == nil {
		res.ID = nil
	}
	return o.Error
}

func (r *recipeRepo) CreateBatch(ctx context.Context, recipes ...*recipe.Recipe) error {
	rps := make([]*models.RecipeModel, 0, len(recipes))
	for _, r2 := range recipes {
		rps = append(rps, models.ParseRecipe(r2))
	}
	res := r.DB.WithContext(ctx).Model(&rps).CreateInBatches(rps, 50)
	return res.Error
}

func (r *recipeRepo) DeleteBatch(ctx context.Context, recipes ...*recipe.Recipe) error {
	rps := make([]int64, 0, len(recipes))
	for idx, r2 := range recipes {
		mdl := models.ParseRecipe(r2)
		if mdl.ID <= 0 {
			return fmt.Errorf("invalid id to delete in batch in %d", idx)
		}
		rps = append(rps, int64(mdl.ID))
	}
	res := r.DB.WithContext(ctx).Delete(&models.RecipeModel{}, rps)
	return res.Error
}

func NewRecipeRepository(db *gorm.DB) recipe.IRecipeRepository {
	return &recipeRepo{GormDB: common.GormDB{DB: db}}
}
