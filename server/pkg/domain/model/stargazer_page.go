package model

type StargazerPage struct {
	ID       uint  `json:"-" gorm:"primary_key"`
	UnixTime int64 `json:"unixTime"`
	Stars    int   `json:"stars"`
}
