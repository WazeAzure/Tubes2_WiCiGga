package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var bfs_slave int = 20
var semaphore = make(chan struct{}, bfs_slave)
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
	n := 0
	BFS(semaphore, current_url, end, &resp, &n)
	fmt.Println("nilai dari depth := ", n)
	return &resp
}

func BFS(semaphore chan struct{}, current_url_list []string, end string, resp *ResponseAPI, depth *int) {
	// semaphore := make(chan struct{}, bfs_slave)
	var temp_url_list []string
	terminate := make(chan struct{})
	stop := false
	var once sync.Once

	for _, elmt := range current_url_list {
		wg.Add(1)
		semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-semaphore }()
			defer wg.Done()

			link_res := scrapWeb(elmt_conc)
		free:
			for _, elmt2 := range link_res {
				select {
				case <-terminate:
					once.Do(func() {
						stop = true
					})
					break free
				default:
					mutex.Lock()
					_, err := visited_bfs[elmt2]
					mutex.Unlock()
					if !err {
						// not exist
						if elmt2 == end {
							// stop right here
							close(terminate)
							fmt.Println(elmt2)
						}
						mutex.Lock()
						visited_bfs[elmt2] = true
						mutex.Unlock()
						temp_url_list = append(temp_url_list, elmt2)
					}

				}
			}
		}(elmt)
		if stop {
			break
		}
		fmt.Println(elmt)
	}

	wg.Wait()
	if stop {
		return
	}
	*depth = *depth + 1
	BFS(semaphore, temp_url_list, end, resp, depth)
}
