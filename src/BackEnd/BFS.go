package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
var bfs_slave int = 10
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
	BFS(current_url, end, &resp, 0)

	return &resp
}

var terminationLock sync.Mutex
var terminationClosed bool

func BFS(current_url_list []string, end string, resp *ResponseAPI, depth int) {
	semaphore := make(chan struct{}, bfs_slave)
	// fmt.Println(current_url)
	var temp_url_list []string
	terminate := make(chan struct{})
	stop := make(chan struct{})
	var once sync.Once

	for _, elmt := range current_url_list {
		wg.Add(1)
		semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-semaphore }()
			defer wg.Done()

			link_res := scrapWeb(elmt_conc)
			fmt.Println(elmt_conc)

			for _, elmt2 := range link_res {
				select {
				case <-terminate:
					once.Do(func() {
						close(stop)
						fmt.Println("BFS STOPPED")
					})
					return
				default:
					mutex.Lock()
					_, err := visited_bfs[elmt2]
					mutex.Unlock()
					if !err {
						// not exist
						if elmt2 == end {
							// stop right
							fmt.Println(elmt2 + "\n\n\n\n\n\n")
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

	}
	wg.Wait()
	select {
	case <-stop:
		fmt.Println("BFS STOPPED FR")
		return
	default:
		BFS(temp_url_list, end, resp, depth+1)
	}
}

// func scrapper(elmt_conc string, end string, temp_url_list *[]string, semaphore <-chan struct{}) {
// 	select {
// 	case <-terminate:
// 		fmt.Println("case finished")
// 		return
// 	default:
// 		defer func() { <-semaphore }()
// 		defer wg.Done()

// 		link_res := scrapWeb(elmt_conc)
// 		fmt.Println(elmt_conc)

// 		for _, elmt2 := range link_res {
// 			mutex.Lock()
// 			_, err := visited_bfs[elmt2]
// 			mutex.Unlock()
// 			if !err {
// 				// not exist
// 				if elmt2 == end {
// 					// stop right here
// 					fmt.Println(elmt2 + "\n\n\n\n\n\n")
// 					stop = true
// 					return
// 				}
// 				mutex.Lock()
// 				visited_bfs[elmt2] = true
// 				mutex.Unlock()
// 				*temp_url_list = append(*temp_url_list, elmt2)
// 			}
// 		}
// 	}
// }
