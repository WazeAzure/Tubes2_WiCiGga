package main

import (
	"fmt"
	"time"
)

func DLS(start string, end string, maxdepth int, visited_dls map[string]bool, saved_path []string, ans *[][]string) bool {
	if start == end {
		*ans = append(*ans, saved_path)
		for _, elmt := range saved_path {
			fmt.Println(elmt)
		}
		return true
	}

	if maxdepth <= 0 {
		return false
	}

	url_scrap := scrapWeb(start)

	url_list := []string{}

	for _, elmt := range url_scrap {
		_, err := visited_dls[elmt]

		if !err {
			visited_dls[elmt] = true
			url_list = append(url_list, elmt)
		}
	}

	for _, elmt := range url_list {
		saved_path2 := append(saved_path, elmt)
		if DLS(elmt, end, maxdepth-1, visited_dls, saved_path2, ans) {
			return true
		}
	}
	return false
}

func IDS(start string, end string) *ResponseAPI {
	var resp ResponseAPI
	defer timeTrack(time.Now(), "IDS", &resp.Time)

	isFound := false
	saved_path := []string{}

	var i int = 0
	for !isFound {
		var visited_dls = map[string]bool{}
		if DLS(start, end, i, visited_dls, saved_path, &resp.Path) {
			isFound = true
			break
		}
		i++
	}

	resp.Status = isFound
	resp.Path = append(resp.Path, saved_path)
	if isFound {
		resp.Message = "Path Found"
		for _, elmt := range saved_path {
			fmt.Println(elmt)
		}
	} else {
		resp.Message = "Path not found"
	}

	return &resp
}
