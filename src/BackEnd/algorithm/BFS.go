package algorithm

import (
	"backend/caching"
	"backend/scraper"
	"backend/util"
	"fmt"
	"sync"
	"time"
)

// global variable. visited nodes in BFS
var visited_bfs = make(map[string]bool)

// gloabl variable. MAP of Parent and Child for BFS to VISUALIZE
var parent_child_bfs = make(map[string]map[string]bool)

// global variable. Maximum concurent slave for bfs
var bfs_slave int = 100

// global variable. Channel to limit concurreny
var semaphore = make(chan struct{}, bfs_slave)

var (
	// mutex for locking and unlocking
	mutex sync.Mutex
	// waitgroup for bfs
	wg sync.WaitGroup
)

// Function to reset the global variables.
func BFSReset() {
	// clear global variable
	for k := range parent_child_bfs {
		delete(parent_child_bfs, k)
	}

	for k := range visited_bfs {
		delete(visited_bfs, k)
	}
}

// Function to handle the BFS request
func BFShandler(url string, end string, ans_type string) *util.ResponseAPI {
	time_start := time.Now()

	var resp util.ResponseAPI

	var current_url = []string{url}
	n := 0
	var x bool

	articleHit := 0
	if ans_type == "multi" {
		x = BFS(&semaphore, current_url, end, &resp, &n, &articleHit)
		time.Sleep(1 * time.Second)
	} else if ans_type == "single" {
		// wg.Add(1)
		BFSSingle(&semaphore, current_url, end, &resp, &n, &articleHit)
		x = true
	}

	util.TimeTrack(time_start, "scrapWeb", &resp.Time)
	fmt.Println("degree dari path bernilai := ", n+1, x)
	resp.Degree = n + 1
	resp.Hit = articleHit

	// path cleaning
	if x {
		var parent_child_bfs_temp = make(map[string]map[string]bool)

		for key, val := range parent_child_bfs {
			if len(val) > 0 {
				parent_child_bfs_temp[key] = val
			}
		}

		var zombie [][]string
		resp.Nodes, resp.Edges = util.ConvertToVisualizerHandler(url, end, parent_child_bfs_temp, n, zombie, "BFS")

		resp.Status = true
	}
	return &resp
}

// Function for BFS Recursion
func BFS(semaphore *chan struct{}, current_url_list []string, end string, resp *util.ResponseAPI, depth *int, articleHit *int) bool {

	fmt.Println("Current Depth ----- ", *depth, " ----")
	// semaphore := make(chan struct{}, bfs_slave)
	var temp_url_list []string
	stop := false

	for _, elmt := range current_url_list {
		*articleHit++
		wg.Add(1)
		*semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-*semaphore }()
			defer wg.Done()

			var link_res []string
			if caching.CheckCacheFile(elmt_conc) {
				// link_res = scraper.ScrapWeb(elmt_conc)
				// caching.SetCacheUrl(elmt_conc, link_res)
				link_res = caching.GetCacheUrl(elmt_conc)
			} else {
				link_res = scraper.ScrapWeb(elmt_conc)
				caching.SetCacheUrl(elmt_conc, link_res)
			}

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
	return BFS(semaphore, temp_url_list, end, resp, depth, articleHit)
}

func BFSSingle(semaphore *chan struct{}, current_url_list []string, end string, resp *util.ResponseAPI, depth *int, articleHit *int) {
	// semaphore := make(chan struct{}, bfs_slave)
	var temp_url_list []string
	terminate := make(chan struct{})
	// stop := false
	var once sync.Once
free:
	for _, elmt := range current_url_list {
		*articleHit++
		wg.Add(1)
		*semaphore <- struct{}{}
		// time.Sleep(50 * time.Millisecond)
		go func(elmt_conc string) {
			defer func() { <-*semaphore }()
			defer wg.Done()

			var link_res []string

			if caching.CheckCacheFile(elmt_conc) {
				link_res = caching.GetCacheUrl(elmt_conc)
			} else {
				link_res = scraper.ScrapWeb(elmt_conc)
				caching.SetCacheUrl(elmt_conc, link_res)
			}

			// free:
			for _, elmt2 := range link_res {
				// select {
				// case <-terminate:
				// 	once.Do(func() {
				// 		stop = true
				// 	})
				// 	break free
				// default:
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
						fmt.Println(elmt_conc, " ======================= PARENT ===========")
						once.Do(func() {
							// stop = true
							close(terminate)
						})

					}
					mutex.Lock()
					visited_bfs[elmt2] = true
					mutex.Unlock()
					temp_url_list = append(temp_url_list, elmt2)
				}

			}
			// }
		}(elmt)
		fmt.Println(elmt)
		select {
		case <-terminate:
			break free
		default:
		}
	}

	wg.Wait()
	select {
	case <-terminate:
		return
	default:
	}

	*depth = *depth + 1
	BFSSingle(semaphore, temp_url_list, end, resp, depth, articleHit)
}
