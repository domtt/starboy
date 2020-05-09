package main

import (
	"log"

	database "github.com/d0minikt/starboy/server/pkg/interface/db"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := env.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to load the database: ", err)
	}
	defer db.Close()

	r := echo.New()
	/*
		token := ""
		github_api.GetStarHistory("/facebook/react", token)
	*/

	//r.GET("/auth/github", service.GithubAuthHandler)

	if cfg.Production {
		r.File("/", "index.html")
		r.Static("/", "./")
	}

	r.Logger.Fatal(r.Start(":" + cfg.Port))
}
