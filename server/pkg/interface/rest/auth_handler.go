package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d0minikt/starboy/server/pkg/domain/service"
	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/labstack/echo"
)

func githubAuthHandler(c echo.Context) error {
	c.Redirect(http.StatusFound, fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		env.Config().GithubClientID,
		env.Config().ServerURL+"/auth/github/callback",
	))
	return nil
}

func githubAuthCallback(c echo.Context) error {
	// 1. Get the query code
	code := c.QueryParam("code")

	// 2. get the access token that works for longer
	accessToken, err := service.CodeToAccessToken(code)
	if err != nil {
		log.Fatal(err)
	}
	// 3. redirect to web app
	c.Redirect(http.StatusFound, env.Config().WebAppURL+"?token="+accessToken)
	return nil
}
