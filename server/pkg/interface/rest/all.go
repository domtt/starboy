package rest

import (
	"github.com/d0minikt/starboy/server/pkg/domain/service"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Services struct {
	repos service.ReposService
}

func New(r service.ReposService) *Services {
	return &Services{r}
}

func (s *Services) Router() *echo.Echo {
	e := echo.New()
	cfg := env.Config()

	e.Use(middleware.CORS())

	e.GET("/auth/github", githubAuthHandler)
	e.GET("/auth/github/callback", githubAuthCallback)

	e.GET("/api/repos/:user/:repo/stars", s.repoStarHistory)
	e.GET("/api/repos/:user/:repo", s.repoInfoHandler)

	// static SPA
	if cfg.Production {
		e.File("/", "index.html")
		e.Static("/", "./")
	}
	return e
}
