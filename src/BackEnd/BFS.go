package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var (
	url_queue   []string
	mutex       sync.Mutex
	wg          sync.WaitGroup
	current_url string
)

func randomTime() time.Duration {
	return 50 + time.Duration(rand.Intn(250))
}

func bfsHandler(url string, end string) *ResponseAPI {
	var resp ResponseAPI

	defer timeTrack(time.Now(), "scrapWeb", &resp.Time)
	fmt.Println("oawjmdkam")
	var current_url = []string{url}

	BFS(current_url, end, &resp, 0)

	return &resp
}

func BFS(current_url_list []string, end string, resp *ResponseAPI, depth int) bool {

	// fmt.Println(current_url)
	var temp_url_list []string
	fmt.Println(current_url_list)
	for _, elmt := range current_url_list {
		link_res := scrapWeb(elmt)
		fmt.Println(elmt)
		for _, elmt2 := range link_res {
			_, err := visited_bfs[elmt2]
			if !err {
				// not exist
				if elmt2 == end {
					//termminates when elmt2 equals end
					return true
				}
				visited_bfs[elmt2] = true
				temp_url_list = append(temp_url_list, elmt2)
			}
		}
	}
	return BFS(temp_url_list, end, resp, depth+1)
}
