package main

import (
	"fmt"
	"time"
)

var visited_node = make(map[string]map[string]bool)

func DLS(start string, end string, maxdepth int, saved_path []string, ans *[][]string) bool {
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

	// cek apakah ada elemen map yang key nya bernilai start di visited_node
	var url_scrap []string
	var url_list = make(map[string]bool)

	url_visited, exist := visited_node[start]
	if !exist {
		url_scrap = scrapWeb(start)
		visited_node[start] = make(map[string]bool)
	} else {
		url_list = url_visited
	}

	for _, elmt := range url_scrap {
		_, err := visited_node[elmt]

		if !err {
			// kalau ga ada di map visited_node
			visited_node[start][elmt] = true
			url_list[elmt] = true
		}
	}

	for key, _ := range url_list {
		saved_path2 := append(saved_path, key)
		fmt.Println(key)
		if DLS(key, end, maxdepth-1, saved_path2, ans) {
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

	var p [][]string

	var i int = 0
	for !isFound {
		if DLS(start, end, i, saved_path, &p) {
			isFound = true
			break
		}
		i++
	}

	resp.Status = isFound
	// resp.Path = append(resp.Path, saved_path)
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
