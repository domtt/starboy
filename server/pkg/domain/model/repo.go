package model

type Repo struct {
	User  string          `json:"user" gorm:"primary_key"`
	Repo  string          `json:"repo" gorm:"primary_key"`
	Stars []StargazerPage `json:"stars"`
}
