package service

import (
	"github.com/d0minikt/starboy/server/pkg/domain/model"
	"github.com/d0minikt/starboy/server/pkg/interface/db"
	"github.com/d0minikt/starboy/server/pkg/interface/github"
)

type ReposService interface {
	GetRepoInfo(string, string, string) []model.StargazerPage
	GetStarHistory(string, string) []model.StargazerPage
}

type reposService struct {
	db db.DB
}

func NewRepos(db db.DB) ReposService {
	return &reposService{db}
}

func (s *reposService) GetRepoInfo(user, repo, token string) []model.StargazerPage {
	stars := github.GetStarHistory(user+"/"+repo, token)
	s.db.CreateRepo(&model.Repo{User: user, Repo: repo, Stars: stars})
	return stars
}

func (s *reposService) GetStarHistory(repo, token string) []model.StargazerPage {
	stars := github.GetStarHistory(repo, token)
	return stars
}
