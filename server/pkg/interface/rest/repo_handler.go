package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

func (s *Services) repoInfoHandler(c echo.Context) error {
	user := c.Param("user")
	repo := c.Param("repo")
	token := c.QueryParam("token")
	c.JSON(http.StatusOK, s.repos.GetRepoInfo(user, repo, token))
	return nil
}

func (s *Services) repoStarHistory(c echo.Context) error {
	user := c.Param("user")
	repo := c.Param("repo")
	token := c.QueryParam("token")
	url := user + "/" + repo

	if len(token) == 0 {
		return errors.New("No token provided")
	}

	entries := s.repos.GetStarHistory(url, token)
	c.JSON(http.StatusOK, map[string]interface{}{
		url: entries,
	})

	return nil
}
