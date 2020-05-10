package main

import (
	"log"

	"github.com/d0minikt/starboy/server/pkg/domain/service"
	database "github.com/d0minikt/starboy/server/pkg/interface/db"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/d0minikt/starboy/server/pkg/interface/rest"
)

func main() {
	cfg := env.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to load the database: ", err)
	}
	defer db.Close()

	dbService := database.New(db)

	r := service.NewRepos(dbService)
	rest := rest.New(r)
	e := rest.Router()

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
