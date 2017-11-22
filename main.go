package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Recentbuild struct {
	Branch          string `json:"branch"`
	BuildUrl        string `json:"build_url"`
	StartTime       string `json:"start_time"`
	BuildTimeMillis int    `json:"build_time_millis"`
	Status          string `json:"status"`
	BuildNum        int    `json:"build_num"`
	UserName        string `json:"username"`
	RepoName        string `json:"reponame"`
}

type Items struct {
	Item []Item `json:"items"`
}

type Item struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Arg      string `json:"arg"`
	Icon     icon   `json:"icon"`
}

type icon struct {
	Path string `json:"path"`
}

func circleci() {
	var token *string = flag.String("t", "secret", "CirclCI Token")
	var filter *string = flag.String("f", "reponame, branch, username, status...", "Search Filter")
	flag.Parse()

	url := "https://circleci.com/api/v1.1/recent-builds?circle-token=" + *token

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var r []Recentbuild
	json.Unmarshal(body, &r)

	var items []Item
	for _, v := range r {
		if strings.Contains(v.RepoName+v.Branch+v.Status+v.UserName, *filter) {
			title := "#" + fmt.Sprint(v.BuildNum) +
				" / " + v.RepoName +
				" / " + v.Branch

			sec := v.BuildTimeMillis / 1000
			subtitle := "user: " + v.UserName +
				" / start: " + v.StartTime +
				" / buildtime: " + fmt.Sprint(sec) + "sec"

			var color string
			if v.Status == "no_tests" || v.Status == "not_run" || v.Status == "not_running" {
				color = "gray"
			} else if v.Status == "fixed" || v.Status == "success" {
				color = "green"
			} else if v.Status == "queued" || v.Status == "scheduled" {
				color = "purple"
			} else if v.Status == "canceled" || v.Status == "failed" || v.Status == "infrastructure_fail" || v.Status == "timeout" {
				color = "red"
			} else if v.Status == "retried" || v.Status == "running" {
				color = "blue"
			}

			items = append(items, Item{
				Title:    title,
				Subtitle: subtitle,
				Arg:      v.BuildUrl,
				Icon:     icon{Path: color + ".png"}})
		}
	}

	j, _ := json.Marshal(Items{Item: items})
	fmt.Println(string(j))
}

func main() {
	circleci()
}
