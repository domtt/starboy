package db

import (
	"fmt"

	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repo struct {
	gorm.Model
	DisplayName string
	GithubPath  string
}

func Connect() (*gorm.DB, error) {
	cfg := env.Config().DB
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.User, cfg.Pass, cfg.Name))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Repo{})
	return db, nil
}
