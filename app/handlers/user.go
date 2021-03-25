package handlers

import (
	"context"
	"errors"
	"recipes/app/common"
	config2 "recipes/app/config"
	"recipes/app/externals"
	"recipes/domains/shared"
	"recipes/domains/user"
	"strings"
	"time"
)

type userHandler struct {
	user user.IUserRepository
	role user.IRoleRepository
	session user.ISessionRepository
	fb externals.IFirebaseAuth
	cfg *config2.Config
}

func (u *userHandler) CheckSession(ctx context.Context, usr *user.User, session *user.Session) error {
	panic("implement me")
}

func (u *userHandler) checkRole(ctx context.Context, rl *user.Role) (err error) {
	_, ltr, err := u.role.GetAll(ctx, shared.DefaultLimiter(1), rl)
	if err != nil {
		return err
	}
	if ltr.Total() != 1 {
		return errors.New("role not found")
	}
	return nil
}
func (u *userHandler) checkUser(ctx context.Context, usr *user.User) (err error) {
	total, err := u.user.Count(ctx, usr)
	if err != nil {
		return err
	}
	if total != 1 {
		return errors.New("user not found")
	}
	return nil
}
func (u *userHandler) getOne(ctx context.Context, usr *user.User) error {
	all, _, err := u.user.GetAll(ctx, shared.DefaultLimiter(1), usr)
	if err != nil {
		return err
	}
	*usr = *all[0]
	return nil
}
func (u *userHandler) Refresh(ctx context.Context, sess *user.Session) error {
	panic("implement me")
}

func (u *userHandler) GetAll(ctx context.Context, actor *user.User, limiter *shared.CommonLimiter, filters ...*user.User) ([]*user.User, *shared.CommonLimiter, error) {
	return u.user.GetAll(ctx, limiter, filters...)
}

func (u *userHandler) Create(ctx context.Context, actor *user.User, usr *user.User) (err error) {
	// TODO: add tracking maybe? when actor is creating user
	rls := usr.Roles
	usr.Roles = nil
	err = u.user.Create(ctx, usr)
	if err != nil {
		return
	}
	return u.AddRole(ctx, actor, usr, rls[0])
}

func (u *userHandler) Update(ctx context.Context, actor *user.User, usr *user.User) error {
	// TODO: only the user itself can update its payload
	if actor.ID != usr.ID {
		return errors.New("forbidden")
	}
	all, total, err := u.user.GetAll(ctx, shared.DefaultLimiter(1), usr)
	if err != nil {
		return err
	}
	if total.Offset.Total != 1 {
		return errors.New("invalid user")
	}
	current := all[0]
	err = u.user.Update(ctx, current)
	if err != nil {
		return err
	}
	*usr = *current
	return nil
}

func (u *userHandler) Delete(ctx context.Context, actor *user.User, usr *user.User) error {
	return u.user.Delete(ctx, usr)
}
func (u *userHandler) AddRole(ctx context.Context, actor *user.User, usr *user.User, role *user.Role) error {
	return common.Try(func() error {
		if actor == nil {
			return nil
		}
		return u.checkUser(ctx, actor)
	}, func() error {
		return u.checkUser(ctx, usr)
	}, func() error {
		return u.checkRole(ctx, role)
	}, func() error {
		return u.user.AddRole(ctx, usr, role)
	})
}

func (u *userHandler) RemoveRole(ctx context.Context, actor *user.User, usr *user.User, role *user.Role) error {
	return common.Try(func() error {
		if actor == nil {
			return nil
		}
		return u.checkUser(ctx, actor)
	}, func() error {
		return u.checkUser(ctx, usr)
	}, func() error {
		return u.checkRole(ctx, role)
	}, func() error {
		return u.user.RemoveRole(ctx, usr, role)
	})
}

func (u *userHandler) Login(ctx context.Context, usr *user.User) (session *user.Session, err error) {
	usr.IsActive = true
	usrs, total, err := u.GetAll(ctx, nil, shared.DefaultLimiter(1), usr)
	if err != nil {
		return
	}
	if total.Total() == 0 {
		return
	}
	usr = usrs[0]
	expire := time.Now().Add(time.Hour)
	refresh := expire.Add(time.Hour * 6)
	session = &user.Session{
		User:             usr,
		Agent:            u.cfg.AppName,
		ExpiresAt:        expire,
		RefreshExpiresAt: refresh,
	}
	err = u.session.Create(ctx, session)
	return
}

func (u *userHandler) Logout(ctx context.Context, usr *user.User, session *user.Session) error {
	return u.session.Delete(ctx, session)
}

func (u *userHandler) FindByToken(ctx context.Context, token string) (usr *user.User, err error) {
	t, err := u.fb.VerifyIDTokenAndCheckRevoked(ctx, token)
	if err != nil {
		return
	}
	firebaseUser, _ := u.fb.GetUser(ctx, t.UID)
	splits := strings.Split(firebaseUser.DisplayName, " ")
	firstName := splits[0]
	lastName := ""
	if len(splits) > 1 {
		lastName = strings.Join(splits[1:], " ")
	}
	usr = &user.User{
		Model:     shared.Model{},
		FirstName: firstName,
		LastName:  lastName,
		Email:     firebaseUser.Email,
	}
	all, total, err := u.user.GetAll(ctx, shared.DefaultLimiter(1), &user.User{Email: firebaseUser.Email})
	if err != nil {
		return
	}
	if total.Total() == 0 {
		usr.Roles = []*user.Role{
			{
				Name: "member",
			},
		}
		usr.IsActive = true
		err = u.Create(ctx, nil, usr)
		return
	}
	if !all[0].IsActive {
		err = errors.New("user is not active")
		return
	}
	usr = all[0]
	return
}

func NewUserHandler(user user.IUserRepository, role user.IRoleRepository, session user.ISessionRepository, fb externals.IFirebaseAuth, cfg *config2.Config) user.IUserUseCase {
	return &userHandler{user, role, session, fb, cfg}
}
