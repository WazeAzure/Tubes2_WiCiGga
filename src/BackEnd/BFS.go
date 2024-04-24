package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var parent_child_bfs = make(map[string]map[string]bool)
var bfs_slave int = 100
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
	time_start := time.Now()
	for k := range parent_child_bfs {
		delete(parent_child_bfs, k)
	}
	var resp ResponseAPI

	var current_url = []string{url}
	n := 0
	x := BFS(semaphore, current_url, end, &resp, &n)

	timeTrack(time_start, "scrapWeb", &resp.Time)

	fmt.Println("degree dari path bernilai := ", n+1, x)

	time.Sleep(1 * time.Second)
	// path cleaning
	if x {
		var parent_child_bfs_temp = make(map[string]map[string]bool)

		for key, val := range parent_child_bfs {
			if key == end {
				fmt.Println(key, " | ", val)
			}
			if len(val) > 0 {
				parent_child_bfs_temp[key] = val
			}
		}

		resp.Nodes, resp.Edges = convertToVisualizerHandler(url, end, parent_child_bfs_temp, n)
	}
	return &resp
}

func BFS(semaphore chan struct{}, current_url_list []string, end string, resp *ResponseAPI, depth *int) bool {
	// semaphore := make(chan struct{}, bfs_slave)
	var temp_url_list []string
	stop := false

	for _, elmt := range current_url_list {
		wg.Add(1)
		semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-semaphore }()
			defer wg.Done()

			link_res := scrapWeb(elmt_conc)
			for _, elmt2 := range link_res {
				mutex.Lock()
				_, err1 := parent_child_bfs[elmt2]
				mutex.Unlock()
				if !err1 {
					mutex.Lock()
					parent_child_bfs[elmt2] = make(map[string]bool)
					parent_child_bfs[elmt2][elmt_conc] = true
					mutex.Unlock()
				} else {
					mutex.Lock()
					parent_child_bfs[elmt2][elmt_conc] = true
					mutex.Unlock()
				}

				mutex.Lock()
				_, err := visited_bfs[elmt2]
				mutex.Unlock()
				if !err {
					// not exist
					if elmt2 == end {
						// stop right here
						stop = true
						fmt.Println(elmt_conc)
					}
					mutex.Lock()
					visited_bfs[elmt2] = true
					mutex.Unlock()
					temp_url_list = append(temp_url_list, elmt2)
				}
			}
		}(elmt)
		fmt.Println(elmt)
	}

	wg.Wait()
	if stop {
		return stop
	}
	*depth = *depth + 1
	return BFS(semaphore, temp_url_list, end, resp, depth)
}
