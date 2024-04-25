package main

import (
	"fmt"
	"time"
)

// var (
// 	mutex sync.Mutex
// 	wg    sync.WaitGroup
// )

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

// CONCURRENT

// func IDS(parent_url string, end string) (*ResponseAPI, [][]string) {
// 	var resp ResponseAPI
// 	multipath := [][]string{}
// 	defer timeTrack(time.Now(), "IDS", &resp.Time)

// 	isFound := false
// 	savedPath := []string{}
// 	stop := make(chan struct{})
// 	semaphore := make(chan struct{}, 20)
// 	var once sync.Once

// 	i := 0
// 	for !isFound {
// 		semaphore <- struct{}{}
// 		go func() {
// 			defer func() { <-semaphore }()
// 			if DLS(start, end, i, savedPath, &multipath) {
// 				once.Do(func() {
// 					close(stop)
// 					isFound = true
// 				})
// 			}
// 		}()
// 		i++
// 	}

// 	<-stop

// 	resp.Status = true
// 	resp.Message = "Path Found"
// 	return &resp, multipath
// }

//NOT CONCURRENT

func IDS(start string, end string) [][]string {
	var resp ResponseAPI
	multipath := [][]string{}
	defer timeTrack(time.Now(), "IDS", &resp.Time)

	isFound := false
	saved_path := []string{}

	var i int = 0
	for !isFound {
		if DLS(start, end, i, saved_path, &multipath) {
			isFound = true
		}
		i++
	}
	// DLS(start, end, 4, saved_path, &multipath)

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
	}
}
