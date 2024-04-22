package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	url_queue   []string
	visited_bfs map[string]string
	mutex       sync.Mutex
	wg          sync.WaitGroup
	slave       int
	current_url string
)

func bfsHandler(url string, end string) *ResponseAPI {
	var resp ResponseAPI

	defer timeTrack(time.Now(), "scrapWeb", &resp.Time)

	current_url = url
	BFS(current_url, end, &resp)
	current_url = url_queue[0]
	url_queue = url_queue[1:]

	for current_url != end {
		if len(url_queue) != 0 {
			go BFS(current_url, end, &resp)
			// give 100 ms delay
			time.Sleep(50 * time.Millisecond)
			current_url = url_queue[0]
			url_queue = url_queue[1:]
		}
	}

	if current_url == end {
		resp.Status = true
		resp.Message = "Path Found"
	} else {
		resp.Status = false
		resp.Message = "No Possible Path Found"
	}

	fmt.Println("=======[ END ]=======")
	fmt.Println(current_url)

	return &resp
}

func BFS(current_url string, end string, resp *ResponseAPI) {

	// fmt.Println(current_url)

	link_res := scrapWeb(current_url)
	for _, elmt := range link_res {
		// check if it not exist
		if !isIn(elmt, url_queue) {
			url_queue = append(url_queue, elmt)
		}
	}

	if current_url == end {
		resp.Status = true
		resp.Message = "Path Found"
		fmt.Println("====================================")
		fmt.Println("------------ FOUND -----------------")
		fmt.Println("====================================")
	}

}
