package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	wg    sync.WaitGroup
)

var visited_node = make(map[string]map[string]bool)

func DLS(start string, end string, maxdepth int, saved_path []string, ans *[][]string) bool {
	defer wg.Done()
	if start == end {
		*ans = append(*ans, saved_path)
		for _, elmt := range saved_path {
			fmt.Println(elmt)
		}
		fmt.Println(start, "=================================================================================================================================================\n\n\n\n\n\n\n\n\n\n")
		return true
	}

	if maxdepth <= 0 {
		return false
	}

	// cek apakah ada elemen map yang key nya bernilai start di visited_node
	var url_scrap []string
	var url_list = make(map[string]bool)

	mutex.Lock()
	url_visited, exist := visited_node[start]
	mutex.Unlock()
	if !exist {
		url_scrap = scrapWeb(start)
		mutex.Lock()
		visited_node[start] = make(map[string]bool)
		mutex.Unlock()
	} else {
		url_list = url_visited
	}

	for _, elmt := range url_scrap {
		mutex.Lock()
		_, err := visited_node[elmt]
		mutex.Unlock()

		if !err {
			// kalau ga ada di map visited_node
			mutex.Lock()
			visited_node[start][elmt] = true
			url_list[elmt] = true
			mutex.Unlock()
		}
	}

	for key, _ := range url_list {
		saved_path2 := append(saved_path, key)
		fmt.Println(key)
		wg.Add(1)
		if DLS(key, end, maxdepth-1, saved_path2, ans) {
			return true
		}
	}
	return false
}

func IDS(start string, end string) [][]string {
	var resp ResponseAPI
	multipath := [][]string{}
	defer timeTrack(time.Now(), "IDS", &resp.Time)

	isFound := false
	// ctx, cancel := context.WithCancel(context.Background())
	saved_path := []string{}
	// var semaphore = make(chan struct{}, 10)

	var i int = 0
	stop := make(chan bool)
	semaphore := make(chan struct{}, 10)
	for !isFound {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(depth int) {
			defer func() { <-semaphore }()
			select {
			case <-stop:
				break
			default:
				if DLS(start, end, depth, saved_path, &multipath) {
					stop <- true
				}
			}
		}(i)
		i++
	}
	wg.Wait()

	resp.Status = isFound
	multipath = append(multipath, saved_path)
	if isFound {
		resp.Message = "Path Found"
		for _, elmt := range saved_path {
			fmt.Println(elmt)
		}
		return multipath
	} else {
		dummy := [][]string{}
		return dummy
		// resp.Message = "Path not found"
	}
}
