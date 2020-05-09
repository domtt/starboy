package github_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

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

type Stargazer struct {
	StarredAt string `json:"starred_at"`
}

// parses github api date format
func parseDate(date string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		panic(err)
	}
	return t
}

func getCurrentStars(repo string, token string) int {
	res := get(fmt.Sprintf("https://api.github.com/repos/%s", repo), token)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	body := map[string]interface{}{}
	json.Unmarshal(bodyBytes, &body)
	return int(body["stargazers_count"].(float64))
}

func getPage(repo string, token string, page int) (time.Time, int) {
	res := get(fmt.Sprintf("https://api.github.com/repos/%s/stargazers?page=%d", repo, page), token)
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	stargazers := []Stargazer{}
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

func getPageCount(repo string, token string) int {
	res := get("https://api.github.com/repos/"+repo+"/stargazers", token)
	linkHeader := res.Header.Get("link")
	ex, _ := regexp.Compile("page\\=([0-9]+)")
	matches := ex.FindAllStringSubmatch(linkHeader, -1)
	pageCount := matches[1][1]
	value, _ := strconv.ParseInt(pageCount, 10, 32)
	return int(value)
}

const sampleCount = 15

type Entry struct {
	Time  int64 `json:"t"`
	Value int   `json:"v"`
}

type StarHistory map[int64]int

// {repo: Entry[]}
var starHistoryCache sync.Map

func GetStarHistory(repo, token string) []Entry {
	// return cached value if present
	cacheResult, ok := starHistoryCache.Load(repo)
	if ok {
		return cacheResult.([]Entry)
	}

	pageCount := getPageCount(repo, token)
	sampleEvery := 1
	if pageCount >= sampleCount {
		sampleEvery = int(pageCount) / sampleCount
	}
	channel := make(chan Entry)
	fetch := func(page int) {
		date, count := getPage(repo, token, page)
		channel <- Entry{Value: count, Time: date.Unix()}
	}
	for i := 0; i < sampleCount+1; i += 1 {
		go fetch(i * sampleEvery)
	}
	current := getCurrentStars(repo, token)
	entries := make([]Entry, sampleCount+1+1)
	for i := 0; i < len(entries)-1; i += 1 {
		entry := <-channel
		entries[i] = entry
	}
	entries[len(entries)-1] = Entry{Value: current, Time: time.Now().Unix()}

	// cache the result
	starHistoryCache.Store(repo, entries)
	// return value
	return entries
}
