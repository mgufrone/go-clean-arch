package common

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"recipes/app/config"
	"recipes/domains/shared"
)

type GormDB struct {
	DB     *gorm.DB
	Config *config.Config
}

func (db *GormDB) CommonFilter(tx *gorm.DB, filter *shared.CommonLimiter) {
	if filter.Offset != nil && filter.Offset.Page > 0 {
		if filter.Offset.PerPage == 0 {
			filter.Offset.PerPage = 10
		}
		tx = tx.
			Limit(int(filter.GetPerPage())).
			Offset(int((filter.Offset.Page - 1) * filter.Offset.PerPage))
	}
	if filter.Sort != nil && filter.Sort.Field != "" {
		if filter.Sort.Direction == "" {
			filter.Sort.Direction = "DESC"
		}
		tx = tx.Order(fmt.Sprintf("%s %s", filter.Sort.Field, filter.Sort.Direction))
	}
	if len(filter.Fields) > 0 {
		tx = tx.Select(filter.Fields)
	}
}

func (db *GormDB) FindAll(ctx context.Context, filter *shared.CommonLimiter, out interface{}) (err error) {
	tx := db.DB.WithContext(ctx)
	db.CommonFilter(tx, filter)
	err = tx.
		Find(out).Error
	return
}

func (db *GormDB) FindById(ctx context.Context, id uint, out interface{}) (err error) {
	tx := db.DB.WithContext(ctx)
	err = tx.Model(out).Where(id).First(out).Error
	return
}

func (db *GormDB) Create(ctx context.Context, value interface{}) (err error) {
	tx := db.DB.WithContext(ctx)
	err = tx.Create(value).Error
	return
}
func (db *GormDB) Update(ctx context.Context, id uint, value interface{}) (err error) {
	tx := db.DB.WithContext(ctx)
	err = tx.Model(value).Where(id).Updates(value).Error
	return
}
func (db *GormDB) Delete(ctx context.Context, id uint) (err error) {
	tx := db.DB.WithContext(ctx)
	err = tx.Delete(id).Error
	return
}
func (db *GormDB) WhereLike(db1 *gorm.DB, field string, value interface{}) *gorm.DB {
	return db1.Where(
		fmt.Sprintf("%s like ?", field),
		fmt.Sprintf("%%%s%%", value),
	)
}
func (db *GormDB) OrWhere(scopes ...func(db *gorm.DB) *gorm.DB) func(db1 *gorm.DB) *gorm.DB {
	return func(db1 *gorm.DB) *gorm.DB {
		for _, sc := range scopes {
			db1 = db1.Or(db1.Scopes(sc))
		}
		return db1
	}
}
func (db *GormDB) Migrate(ctx context.Context, model interface{}) (err error) {
	tx := db.DB.WithContext(ctx)
	err = tx.AutoMigrate(model)
	return
}

func InitializeDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&allowNativePasswords=true&allowOldPasswords=true",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_DB"),
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields: true,
	})
}
