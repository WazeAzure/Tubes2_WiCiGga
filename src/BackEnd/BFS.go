package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var bfs_slave int = 20
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
	// stop := make(chan struct{})
	x := BFS(current_url, end, &resp, 0)

	fmt.Println("nilai dari x := ", x)
	return &resp
}

var terminationLock sync.Mutex
var terminationClosed bool

func BFS(current_url_list []string, end string, resp *ResponseAPI, depth int) bool {
	semaphore := make(chan struct{}, bfs_slave)
	var temp_url_list []string
	terminate := make(chan struct{})
	stop := make(chan struct{})
	var once sync.Once

	x := false
	for _, elmt := range current_url_list {
		wg.Add(1)
		semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-semaphore }()
			defer wg.Done()

			link_res := scrapWeb(elmt_conc)
			if !x {
				fmt.Println(elmt_conc)
			}

			for _, elmt2 := range link_res {
				select {
				case <-terminate:
					once.Do(func() {
						close(stop)
						fmt.Println("BFS STOPPED")
						x = true
					})
				default:
					mutex.Lock()
					_, err := visited_bfs[elmt2]
					mutex.Unlock()
					if !err {
						// not exist
						if elmt2 == end {
							// stop right
							fmt.Println(elmt2)
							fmt.Println(depth)
							close(terminate)
							return
						}
						mutex.Lock()
						visited_bfs[elmt2] = true
						mutex.Unlock()
						temp_url_list = append(temp_url_list, elmt2)
					}

				}
			}
		}(elmt)

		if x {
			return true
		}
	}

	wg.Wait()
	select {
	case <-stop:
		return true
	default:
		return BFS(temp_url_list, end, resp, depth+1)
	}
}
