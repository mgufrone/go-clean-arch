package repositories

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"recipes/app/models"
	"recipes/app/repositories"
	"recipes/domains/shared"
	"recipes/domains/stats"
	"testing"
)

type ViewRepositorySuite struct {
	suite.Suite
	repo stats.IViewRepository
	ctx context.Context
	db *gorm.DB
}
func (v *ViewRepositorySuite) SetupSuite() {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&allowNativePasswords=true&allowOldPasswords=true",
		"root",
		"root",
		"localhost",
		"3306",
		"testing",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields: true,
	})
	if err != nil {
		panic(err)
	}
	v.db = db.Debug()
	err = db.AutoMigrate(&models.ViewModel{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.UserModel{}, &models.RoleModel{})
	if err != nil {
		panic(err)
	}
	v.repo = repositories.NewViewRepository(db)
	v.ctx = context.Background()
}
func (v *ViewRepositorySuite) TearDownSuite() {
	
}

func (v *ViewRepositorySuite) Test01GetAll() {
	_, _, err := v.repo.GetAll(v.ctx, shared.DefaultLimiter(1))
	assert.Nil(v.T(), err)
}
func (v *ViewRepositorySuite) Test02GetPopular() {
	limiter := shared.DefaultLimiter(100)
	limiter.Sort = &shared.Sort{
		Field:     "popular",
		Direction: "desc",
	}
	_, _, err := v.repo.GetAll(v.ctx, limiter)
	assert.Nil(v.T(), err)
}

func TestViewRepository(t *testing.T) {
	suite.Run(t, new(ViewRepositorySuite))
}
