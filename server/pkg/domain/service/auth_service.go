package service

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/d0minikt/starboy/server/pkg/interface/env"
	"github.com/labstack/echo/v4"
)

func codeToAccessToken(code string) (string, error) {
	config := env.Config()
	resp, err := http.Post(
		fmt.Sprintf(
			"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
			config.GithubClientID,
			config.GithubClientSecret,
			code,
		),
		"application/json",
		bytes.NewBufferString(""),
	)
	if err != nil {
		return "", errors.New("failed to fetch access token")
	}

	// read body response as query params
	body, _ := ioutil.ReadAll(resp.Body)
	values, err := url.ParseQuery(string(body))
	if err != nil || len(values["access_token"]) == 0 {
		log.Println(err)
		return "", errors.New("access token fetched, but failed to parse response")
	}
	token := values["access_token"][0]

	return token, nil
}

func GithubAuthHandler(c echo.Context) error {
	// 1. Get the query code
	code := c.QueryParam("code")

	// 2. get the access token that works for longer
	accessToken, err := codeToAccessToken(code)
	if err != nil {
		log.Fatal(err)
	}
	return c.Redirect(http.StatusFound, "http://localhost:3001?token="+accessToken)
}
