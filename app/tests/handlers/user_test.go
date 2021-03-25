package handlers

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"recipes/app/config"
	"recipes/app/handlers"
	"recipes/app/tests/externals"
	"recipes/domains/shared"
	"recipes/domains/user"
	"recipes/domains/user/mocks"
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

type UserHandlerSuite struct {
	suite.Suite
	userRepo *mocks.UserRepository
	fb *externals.FbAuth
	roleRepo *mocks.MockRoleRepo
	sessRepo user.ISessionRepository
	handler user.IUserUseCase
	sharedContext context.Context
}

func (u *UserHandlerSuite) SetupSuite() {
	u.sharedContext = context.Background()
	u.userRepo = mocks.NewUserMockRepository()
	u.roleRepo = mocks.NewMockRoleRepository()
	u.sessRepo = mocks.NewMockSessionRepository()
	u.fb = externals.NewMockFirebaseAuth()
	u.handler = handlers.NewUserHandler(
		u.userRepo,
		u.roleRepo,
		u.sessRepo,
		u.fb,
		&config.Config{
			AppName:            "Testing",
		},
	)
}

func (u *UserHandlerSuite) TearDownSuite() {
}

func (u *UserHandlerSuite) TestLoginSuccess() {
	usr := &user.User{
		Email:     "test@test.com",
		FirstName: "test",
		LastName:  "test",
	}
	mockLimiter := shared.DefaultLimiter(1)
	mockLimiter.Offset.Total = uint64(1)
	u.userRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
		Return([]*user.User{usr}, mockLimiter, nil)
	sess, err := u.handler.Login(u.sharedContext, usr)
	if assert.Nil(u.T(), err) {
		assert.NotNil(u.T(), sess.SessionID)
		assert.Greater(u.T(), sess.ExpiresAt.Unix(), time.Now().Unix())
		assert.Greater(u.T(), sess.RefreshExpiresAt.Unix(), time.Now().Unix())
	}
}
func (u *UserHandlerSuite) TestFindByTokenSuccess() {
	limiter := shared.DefaultLimiter(1)
	limiter.Offset = &shared.OffsetPagination{Total: 0}
	limiter.Cursor = &shared.CursorPagination{Total: 0}
	u.userRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return([]*user.User{}, limiter , nil).Times(1)
	u.userRepo.On("AddRole", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	u.userRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(1)
	firstName := faker.Name().FirstName()
	lastName := faker.Name().LastName()
	displayName := fmt.Sprintf("%s %s", firstName, lastName)
	email := faker.Internet().Email()
	uid := uuid.New().String()
	u.fb.On("VerifyIDTokenAndCheckRevoked", mock.Anything, mock.Anything).Return(&auth.Token{}, nil).Times(1)
	u.fb.On("GetUser", mock.Anything, mock.Anything).Return(&auth.UserRecord{
		UserInfo:      &auth.UserInfo{
			DisplayName: displayName,
			Email:       email,
			UID: uid,
		},
		EmailVerified: true,
	}, nil).Times(1)
	usr, err := u.handler.FindByToken(u.sharedContext, "randomtoken")
	u.userRepo.AssertExpectations(u.T())
	u.fb.AssertExpectations(u.T())
	if assert.Nil(u.T(), err) {
		assert.Equal(u.T(), u.userRepo.Calls[1].Arguments[1], usr)
		assert.Equal(u.T(), usr.Email, email)
		assert.Equal(u.T(), usr.FirstName, firstName)
		assert.Equal(u.T(), usr.LastName, lastName)
		assert.True(u.T(), usr.IsActive)
		assert.Equal(u.T(), usr.Roles, []*user.Role{{
			Name: "member",
		}})
	}
}

func TestUserHandler(t *testing.T)  {
	suite.Run(t, new(UserHandlerSuite))
}
