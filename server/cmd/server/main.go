package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d0minikt/starboy/server/pkg/domain/service"
	database "github.com/d0minikt/starboy/server/pkg/interface/db"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/d0minikt/starboy/server/pkg/interface/github_api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := env.Load()
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to load the database: ", err)
	}
	defer db.Close()

	r := echo.New()
	r.Use(middleware.CORS())

	r.GET("/auth/github", service.GithubAuthHandler)
	r.GET("/auth/github/callback", service.GithubAuthCallback)

	r.GET("/api/repo/:user/:repo", f)

	if cfg.Production {
		r.File("/", "index.html")
		r.Static("/", "./")
	}

	r.Logger.Fatal(r.Start(":" + cfg.Port))
}

func f(c echo.Context) error {
	user := c.Param("user")
	repo := c.Param("repo")
	token := c.QueryParam("token")
	url := user + "/" + repo

	fmt.Println("req")
	entries := github_api.GetStarHistory(url, token)
	c.JSON(http.StatusOK, map[string]interface{}{
		url: entries,
	})

	return nil
}
