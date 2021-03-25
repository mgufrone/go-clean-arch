package repositories

import (
	"context"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"recipes/app/config"
	"recipes/app/models"
	"recipes/app/repositories"
	"recipes/domains/shared"
	"recipes/domains/user"
	"syreclabs.com/go/faker"
	"testing"
)

type UserRepositorySuite struct {
	suite.Suite
	repository user.IUserRepository
	ctx        context.Context
}

func (u *UserRepositorySuite) SetupSuite() {
	cfg := &config.Config{AppName: "testing"}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.UserModel{}, &models.RoleModel{})
	if err != nil {
		panic(err)
	}
	u.repository = repositories.NewUserRepository(db, cfg)
	u.ctx = context.Background()
}

func (u *UserRepositorySuite) TearDownSuite() {
	os.Remove("test.db")
}

func (u *UserRepositorySuite) TestGetAll() {
	u.repository.GetAll(u.ctx, shared.DefaultLimiter(1))
}
func (u *UserRepositorySuite) TestGetAllWithFilter() {
	u.repository.GetAll(u.ctx, shared.DefaultLimiter(1), &user.User{
		FirstName: "hero",
		LastName:  "ola",
	}, &user.User{
		LastName:  "hora",
		FirstName: "comoesta",
	}, &user.User{
		LastName:  faker.Name().LastName(),
		Email: faker.Internet().Email(),
	},
	)
}

func TestInitializeRepository(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}
