package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/d0minikt/starboy/server/pkg/domain/model"
)

const sampleCount = 15

type stargazer struct {
	StarredAt string `json:"starred_at"`
}

func GetStarHistory(repo, token string) []model.StargazerPage {
	pageCount := getPageCount(repo, token)
	sampleEvery := 1
	if pageCount >= sampleCount {
		sampleEvery = int(pageCount) / sampleCount
	}
	channel := make(chan model.StargazerPage)
	fetch := func(page int) {
		date, count := getPage(repo, token, page)
		channel <- model.StargazerPage{Stars: count, UnixTime: date.Unix()}
	}
	// fetch all page values in parallel
	for i := 0; i < sampleCount+1; i += 1 {
		go fetch(i * sampleEvery)
	}
	// wait for all pallel fetches to finish and join the array accordingly
	entries := make([]model.StargazerPage, sampleCount+1+1)
	for i := 0; i < len(entries)-1; i += 1 {
		entry := <-channel
		entries[i] = entry
	}
	// regardless of sample size, make sure last value is up-to-date
	entries[len(entries)-1] = model.StargazerPage{
		Stars:    GetStarCount(repo, token),
		UnixTime: time.Now().Unix(),
	}

	return entries
}

// gets star count of a given github repository
func GetStarCount(repo string, token string) int {
	res := get(fmt.Sprintf("https://api.github.com/repos/%s", repo), token)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	body := map[string]interface{}{}
	json.Unmarshal(bodyBytes, &body)
	return int(body["stargazers_count"].(float64))
}

// utility function for getting data from GitHub
func get(path string, token string) *http.Response {
	c := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Accept", "application/vnd.github.v3.star+json")
	res, err := c.Do(req)
	return res
}

// parses GitHub api date format
func parseDate(date string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		panic(err)
	}
	return t
}

// GitHub stargazers are split into pages, each 30 entries
// this function returns the date of this page as well as how many users starred this page at that date
func getPage(repo string, token string, page int) (time.Time, int) {
	res := get(fmt.Sprintf("https://api.github.com/repos/%s/stargazers?page=%d", repo, page), token)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	stargazers := []stargazer{}
	json.Unmarshal(bodyBytes, &stargazers)
	if len(stargazers) == 0 {
		return time.Now(), 0
	}
	totalCount := (page - 1) * 30
	if page == 0 {
		totalCount = 0
	}
	starredAt := parseDate(stargazers[0].StarredAt)
	return starredAt, totalCount
}

// returns number of pages the stargazers are split into
func getPageCount(repo string, token string) int {
	res := get("https://api.github.com/repos/"+repo+"/stargazers", token)
	linkHeader := res.Header.Get("link")
	ex, _ := regexp.Compile("page\\=([0-9]+)")
	matches := ex.FindAllStringSubmatch(linkHeader, -1)
	pageCount := matches[1][1]
	value, _ := strconv.ParseInt(pageCount, 10, 32)
	return int(value)
}
