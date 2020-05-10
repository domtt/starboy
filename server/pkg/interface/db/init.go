package db

import (
	"fmt"

	"github.com/d0minikt/starboy/server/pkg/domain/model"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect() (*gorm.DB, error) {
	cfg := env.Config().DB
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.User, cfg.Pass, cfg.Name))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Repo{})
	db.AutoMigrate(&model.StargazerPage{})
	return db, nil
}

type DB interface {
	CreateRepo(r *model.Repo)
}

type database struct {
	db *gorm.DB
}

func New(db *gorm.DB) DB {
	return &database{db: db}
}

func (db *database) CreateRepo(r *model.Repo) {
	db.db.Create(r)
}
