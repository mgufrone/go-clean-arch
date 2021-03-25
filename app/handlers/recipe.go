package handlers

import (
	"context"
	"errors"
	"fmt"
	"recipes/app/common"
	"recipes/app/models"
	"recipes/domains/recipe"
	"recipes/domains/shared"
	"recipes/domains/stats"
	"recipes/domains/user"
	"sync"
)

type recipeHandler struct {
	repo recipe.IRecipeRepository
	user user.IUserRepository
	view stats.IViewRepository
	like stats.ILikeRepository
}

func (r *recipeHandler) checkUser(ctx context.Context, usr *user.User) (err error) {
	total, err := r.user.Count(ctx, usr)
	if err != nil {
		return err
	}
	if total != 1 {
		return errors.New("user not found")
	}
	return nil
}
func (r *recipeHandler) checkRecipe(ctx context.Context, rcp *recipe.Recipe) (err error) {
	total, err := r.repo.Count(ctx, rcp)
	if err != nil {
		return err
	}
	if total != 1 {
		return errors.New("recipe not found")
	}
	return nil
}
func (r *recipeHandler) canManage(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	if err = r.checkExistence(ctx, usr, rcp); err != nil {
		return
	}
	if err = r.getOne(ctx, rcp); err != nil {
		return err
	}
	if usr.ID != rcp.User.ID {
		return errors.New("action denied")
	}
	return
}
func (r *recipeHandler) getOne(ctx context.Context, rcp *recipe.Recipe) error {
	all, _, err := r.repo.GetAll(ctx, shared.DefaultLimiter(1), rcp)
	if err != nil {
		return err
	}
	*rcp = *all[0]
	return nil
}
func (r *recipeHandler) checkExistence(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	if err = r.checkUser(ctx, usr); err != nil {
		return
	}
	rcp2 := &recipe.Recipe{
		Model: shared.Model{
			ID: rcp.ID,
		},
	}
	if err = r.checkRecipe(ctx, rcp2); err != nil {
		return
	}
	return
}

func (r *recipeHandler) getDependencies(ctx context.Context, wg *sync.WaitGroup, res1 *recipe.StatsRecipe) {
	wg.Add(2)
	go func(wg1 *sync.WaitGroup, str *recipe.StatsRecipe) {
		likeCount, _ := r.like.CountByReference(ctx, "recipes", uint64(str.Recipe.ID.(uint)))
		str.LikeCount = uint64(likeCount)
		wg1.Done()
	}(wg, res1)
	go func(wg1 *sync.WaitGroup, str *recipe.StatsRecipe) {
		ltr := shared.DefaultLimiter(1)
		ltr.Fields = []string{"first_name", "last_name"}
		users, _, err := r.user.GetAll(ctx, ltr, &user.User{
			Model:     shared.Model{
				ID: str.User.ID,
			},
		})
		if err == nil {
			str.User = users[0]
		}
		wg1.Done()
	}(wg, res1)
}
func (r *recipeHandler) GetAll(ctx context.Context, limiter *shared.CommonLimiter, filters ...*recipe.Recipe) ([]*recipe.StatsRecipe, *shared.CommonLimiter, error) {
	if (limiter.Offset != nil && limiter.Offset.PerPage == 0) ||
		(limiter.Cursor != nil && limiter.Cursor.PerPage == 0) {
		return nil, nil, errors.New("set pagination appropriately")
	}
	res, l, err := r.repo.GetAll(ctx, limiter, filters...)
	if err != nil {
		return nil, l, err
	}
	result := make([]*recipe.StatsRecipe, 0, len(res))
	if len(res) > 0 {
		var wg sync.WaitGroup
		for _, rcp := range res {
			res1 := &recipe.StatsRecipe{
				Recipe: rcp,
			}
			r.getDependencies(ctx, &wg, res1)
			result = append(result, res1)
		}
		wg.Wait()
	}
	return result, l, err
}

func (r *recipeHandler) GetByUser(ctx context.Context, limiter *shared.CommonLimiter, usr *user.User, filters ...*recipe.Recipe) ([]*recipe.StatsRecipe, *shared.CommonLimiter, error) {
	if err := r.checkUser(ctx, usr); err != nil {
		return nil, limiter, err
	}
	filters = append([]*recipe.Recipe{
		{
			User: usr,
		},
	}, filters...)
	return r.GetAll(ctx, limiter, filters...)
}

func (r *recipeHandler) Create(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (error) {
	rcp.User = usr
	mdl := models.ParseRecipe(rcp)
	return common.Try(func() error {
		ctx2 := context.WithValue(ctx, common.ActionKey, common.ActionCreate)
		return mdl.ValidateWithContext(ctx2)
	}, func() error {
		if usr.ID == nil || usr.ID.(uint) < 1 {
			return errors.New("invalid user")
		}
		return nil
	}, func() error {
		return r.checkUser(ctx, usr)
	}, func() error {
		check := &recipe.Recipe{
			Title: rcp.Title,
			User: rcp.User,
		}
		if err := r.checkRecipe(ctx, check); err == nil {
			return errors.New("recipe existed")
		}
		return nil
	}, func() error {
		return r.repo.Create(ctx, rcp)
	})
}

func (r *recipeHandler) Update(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	check := &recipe.Recipe{
		Model: shared.Model{ID: rcp.ID},
	}
	limiter := shared.DefaultLimiter(1)
	limiter.Fields = []string{"id", "user_id"}
	existing, l, err := r.GetAll(ctx, limiter, check)
	if err != nil {
		return err
	}
	if l.Offset.Total == 0 {
		return errors.New("recipe not found")
	}
	eRcp := existing[0]
	mdl := models.ParseRecipe(rcp)
	ctx2 := context.WithValue(ctx, common.ActionKey, common.ActionUpdate)
	if err = mdl.ValidateWithContext(ctx2); err != nil {
		return err
	}
	if usr.ID == nil || usr.ID.(uint64) < 1 {
		return errors.New("invalid user")
	}
	if eRcp.User.ID != usr.ID {
		return errors.New("only owner allowed to update its recipe")
	}
	if err = r.checkUser(ctx, usr); err != nil {
		return err
	}
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) Delete(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	return r.repo.Delete(ctx, rcp)
}

func (r *recipeHandler) CreateBatch(ctx context.Context, usr *user.User, recipes ...*recipe.Recipe) (err error) {
	panic("implement me")
}

func (r *recipeHandler) DeleteBatch(ctx context.Context, usr *user.User, recipes ...*recipe.Recipe) (err error) {
	panic("implement me")
}

func (r *recipeHandler) Popular(ctx context.Context, usr *user.User, limiter *shared.CommonLimiter, filters ...*recipe.Recipe) ([]*recipe.StatsRecipe, *shared.CommonLimiter, error) {
	// do reverse query.
	var (
		ids []*stats.View
		res []*recipe.StatsRecipe
	)
	cp := &shared.CommonLimiter{}
	cp2 := &shared.CommonLimiter{}

	if len(filters) > 0 {
		*cp = *limiter
		cp.Offset.Total = 0
		cp.Fields = []string{"id"}
		res1, l, err := r.repo.GetAll(ctx, cp, filters...)
		if err != nil {
			return nil, l, err
		}
		for _, r2 := range res1 {
			ids = append(ids, &stats.View{Reference: common.RecipeReference, ReferenceID: uint64(r2.ID.(uint))})
		}
	}
	*cp2 = *limiter
	if len(ids) == 0 {
		ids = []*stats.View{{Reference: common.RecipeReference}}
	}
	cp2.Sort = &shared.Sort{Field: "view_count", Direction: shared.SortDesc}
	cp2.Fields = []string{}
	views, l, err := r.view.GroupCount(ctx, cp2, ids...)
	if err != nil {
		return nil, l, err
	}
	if l.Total() == 0 {
		return nil, l, err
	}
	fmt.Println(len(views), cp2.GetPerPage(), limiter.GetPerPage())
	if len(views) > 0 {
		var wg sync.WaitGroup
		for _, v1 := range views {
			r1, _, err1 := r.repo.GetAll(ctx, shared.DefaultLimiter(1), &recipe.Recipe{
				Model: shared.Model{
					ID: uint(v1.ReferenceID),
				},
			})
			if err1 != nil {
				continue
			}
			stRecipe := &recipe.StatsRecipe{
				Recipe:    r1[0],
				ViewCount: uint64(v1.Count),
			}
			r.getDependencies(ctx, &wg, stRecipe)
			res = append(res, stRecipe)
		}
		wg.Wait()
	}
	return res, l, nil
}

func (r *recipeHandler) Dislike(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	if err = r.checkExistence(ctx, usr, rcp); err != nil {
		return
	}
	like := &stats.Like{
		User: usr,
		ReferenceID: uint64(rcp.ID.(uint)),
		Reference: common.RecipeReference,
	}
	count, err := r.like.Count(ctx, like)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("action already performed")
	}
	return r.like.Delete(ctx, like)
}
func (r *recipeHandler) Like(ctx context.Context, usr *user.User, rcp *recipe.Recipe) (err error) {
	var like *stats.Like
	return common.Try(func() error {
		return r.checkExistence(ctx, usr, rcp)
	}, func() error {
		if rcp.User.ID != usr.ID {
			return nil
		}
		return errors.New("owner cannot like its own recipe")
	}, func() error {
		like = &stats.Like{
			User: usr,
			ReferenceID: uint64(rcp.ID.(uint)),
			Reference: common.RecipeReference,
		}
		count, err1 := r.like.Count(ctx, like)
		if err1 != nil {
			return err1
		}
		if count > 0 {
			return errors.New("action is already performed")
		}
		return nil
	}, func() error {
		return r.like.Put(ctx, like)
	})
}

func (r *recipeHandler) View(ctx context.Context, usr *user.User, rcp *recipe.Recipe) error {
	var view *stats.View
	return common.Try(func() error {
		return r.checkExistence(ctx, usr, rcp)
	}, func() error {
		view = &stats.View{
			User: usr,
			ReferenceID: uint64(rcp.ID.(uint)),
			Reference: common.RecipeReference,
		}
		count, err1 := r.view.Count(ctx, view)
		if err1 != nil {
			return err1
		}
		if count > 0 {
			return errors.New("action already performed")
		}
		return nil
	}, func() error {
		return r.view.Put(ctx, view)
	})
}

func (r *recipeHandler) AddStep(ctx context.Context, usr *user.User, rcp *recipe.Recipe, step string) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	rcp.Steps = append(rcp.Steps, &recipe.Step{
		Sequence: uint(len(rcp.Steps)),
		Step:     step,
	})
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) UpdateStep(ctx context.Context, usr *user.User, rcp *recipe.Recipe, sequence int, step string) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if sequence > len(rcp.Steps) {
		return errors.New("sequence not found")
	}
	if rcp.Steps[sequence].Step == step {
		return errors.New("step is the same, no need to udpate")
	}
	rcp.Steps[sequence].Step = step
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) RemoveStep(ctx context.Context, usr *user.User, rcp *recipe.Recipe, sequence int) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if sequence > len(rcp.Steps) {
		return errors.New("sequence not found")
	}
	lastIndex := 0
	updated := make([]*recipe.Step, 0, len(rcp.Steps) - 1)
	for _, seq := range rcp.Steps {
		if seq.Sequence == uint(sequence) {
			continue
		}
		seq.Sequence = uint(lastIndex)
		updated = append(updated, seq)
		lastIndex++
	}
	rcp.Steps = updated
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) AddIngredient(ctx context.Context, usr *user.User, rcp *recipe.Recipe, ingredient *recipe.Ingredient) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	rcp.Ingredients = append(rcp.Ingredients, ingredient)
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) UpdateIngredient(ctx context.Context, usr *user.User, rcp *recipe.Recipe, ingredient *recipe.Ingredient) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if ingredient.Sequence > len(rcp.Ingredients) {
		return errors.New("sequence not found")
	}
	rcp.Ingredients[ingredient.Sequence] = ingredient
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) RemoveIngredient(ctx context.Context, usr *user.User, rcp *recipe.Recipe, sequence int) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if sequence > len(rcp.Ingredients) {
		return errors.New("sequence not found")
	}
	lastIndex := 0
	updated := make([]*recipe.Ingredient, 0, len(rcp.Ingredients) - 1)
	for _, seq := range rcp.Ingredients {
		if seq.Sequence == sequence {
			continue
		}
		seq.Sequence = lastIndex
		updated = append(updated, seq)
		lastIndex++
	}
	rcp.Ingredients = updated
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) AddPhoto(ctx context.Context, usr *user.User, rcp *recipe.Recipe, photo *recipe.Photo) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	rcp.Photos = append(rcp.Photos, photo)
	return
}

func (r *recipeHandler) UpdatePhoto(ctx context.Context, usr *user.User, rcp *recipe.Recipe, photo *recipe.Photo) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if photo.Sequence > uint(len(rcp.Photos)) {
		return errors.New("sequence not found")
	}
	rcp.Photos[photo.Sequence] = photo
	return r.repo.Update(ctx, rcp)
}

func (r *recipeHandler) RemovePhoto(ctx context.Context, usr *user.User, rcp *recipe.Recipe, sequence int) (err error) {
	if err = r.canManage(ctx, usr, rcp); err != nil {
		return
	}
	_ = r.getOne(ctx, rcp)
	if sequence > len(rcp.Photos) {
		return errors.New("sequence not found")
	}
	lastIndex := 0
	updated := make([]*recipe.Photo, 0, len(rcp.Photos) - 1)
	for _, seq := range rcp.Photos {
		if seq.Sequence == uint(sequence) {
			continue
		}
		seq.Sequence = uint(lastIndex)
		updated = append(updated, seq)
		lastIndex++
	}
	rcp.Photos = updated
	return r.repo.Update(ctx, rcp)
}

func NewRecipeHandler(
	repo recipe.IRecipeRepository,
	userRepository user.IUserRepository,
	likeRepository stats.ILikeRepository,
	viewRepository stats.IViewRepository,
) recipe.IRecipeUseCase  {
	return &recipeHandler{repo, userRepository, viewRepository, likeRepository}
}
