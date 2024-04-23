package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var bfs_slave int = 1000
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
	var current_url = []string{url}

	BFS(current_url, end, &resp, 0)

	return &resp
}

func BFS(current_url_list []string, end string, resp *ResponseAPI, depth int) {
	semaphore := make(chan struct{}, bfs_slave)
	// fmt.Println(current_url)
	var temp_url_list []string
	for _, elmt := range current_url_list {
		wg.Add(1)

		semaphore <- struct{}{}
		time.Sleep(100 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-semaphore }()
			defer wg.Done()

			link_res := scrapWeb(elmt_conc)
			fmt.Println(elmt_conc)

			for _, elmt2 := range link_res {
				// mutex.Lock()
				_, err := visited_bfs[elmt2]
				// mutex.Unlock()
				if !err {
					// not exist
					if elmt2 == end {
						fmt.Println(depth)
						return
					}
					// mutex.Lock()
					visited_bfs[elmt2] = true
					// mutex.Unlock()
					temp_url_list = append(temp_url_list, elmt2)
				}
			}
		}(elmt)

	}
	wg.Wait()

	BFS(temp_url_list, end, resp, depth+1)
}
