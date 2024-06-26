package algorithm

import (
	"backend/caching"
	"backend/scraper"
	"backend/util"
	"fmt"
	"sync"
	"time"
)

// global variable to track IDS visited nodes
var visited_node = make(map[string]map[string]bool)

func IDS(start string, end string, resp *util.ResponseAPI) [][]string {
	defer util.TimeTrack(time.Now(), "IDS", &resp.Time)

	multipath := [][]string{}

	isFound := false
	saved_path := []string{}

	articleHit := 0
	var i int = 0
	var wait sync.WaitGroup
	for !isFound {
		// wg.Add(1)
		fmt.Println("=============================================================")
		fmt.Println("|                         CURRENT DEPTH ", i, "              |")
		fmt.Println("=============================================================")
		if DLS(start, end, i, saved_path, &multipath, &wait, &articleHit) {
			isFound = true
		}
		wait.Wait()
		time.Sleep(1 * time.Second)
		i++
	}

	resp.Degree = i - 1
	resp.Hit = articleHit
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

func DLS(start string, end string, maxdepth int, saved_path []string, ans *[][]string, wg *sync.WaitGroup, articleHit *int) bool {
	*articleHit++
	if start == end {
		fmt.Println("\n\n\nFOUND\n\n\n")
		*ans = append(*ans, saved_path)
		// for _, elmt := range saved_path {
		// 	fmt.Println(elmt)
		// }
		// fmt.Println(start, "=================================================================================================================================================\n\n\n\n\n\n\n\n\n\n")
		return true
	}

	if maxdepth <= 0 {
		return false
	}

	// cek apakah ada elemen map yang key nya bernilai start di visited_node
	var url_scrap []string
	var url_list = make(map[string]bool)

	mutex.Lock()
	// url_visited, exist := visited_node[start]
	visited_node[start] = make(map[string]bool)
	mutex.Unlock()
	if caching.CheckCacheFile(start) {
		url_scrap = caching.GetCacheUrl(start)
	} else {
		mutex.Lock()
		url_scrap = scraper.ScrapWeb(start)
		caching.SetCacheUrl(start, url_scrap)
		mutex.Unlock()
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
		} else {
			// kalau ada
			mutex.Lock()
			url_list[elmt] = true
			mutex.Unlock()
		}
	}

	stop := false
	semaphore := make(chan struct{}, 3)
	for key := range url_list {
		semaphore <- struct{}{}
		saved_path2 := append(saved_path, key)
		fmt.Println(key, maxdepth)
		(*wg).Add(1)
		go func(key_conc string) {
			defer func() { <-semaphore }()
			defer (*wg).Done()
			x := DLS(key_conc, end, maxdepth-1, saved_path2, ans, wg, articleHit)
			if x {
				stop = true
			}
		}(key)
	}

	// wg.Wait()

	return stop
}

func IDSSingle(start string, end string, resp *util.ResponseAPI) [][]string {
	defer util.TimeTrack(time.Now(), "IDSSingle", &resp.Time)

	multipath := [][]string{}
	saved_path := []string{}

	isFound := false

	articleHit := 0
	// ctx, cancel := context.WithCancel(context.Background())

	// var semaphore = make(chan struct{}, 10)

	var i int = 0
	for !isFound {
		if DLSSingle(start, end, i, saved_path, &multipath, &articleHit) {
			isFound = true
		}
		i++
	}
	resp.Degree = i - 1
	resp.Hit = articleHit

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

func DLSSingle(start string, end string, maxdepth int, saved_path []string, ans *[][]string, articleHit *int) bool {
	*articleHit++
	if start == end {
		*ans = append(*ans, saved_path)
		// for _, elmt := range saved_path {
		// 	fmt.Println(elmt)
		// }
		// fmt.Println(start, "=================================================================================================================================================\n\n\n\n\n\n\n\n\n\n")
		return true
	}

	if maxdepth <= 0 {
		return false
	}

	// cek apakah ada elemen map yang key nya bernilai start di visited_node
	var url_scrap []string
	var url_list = make(map[string]bool)

	mutex.Lock()
	// url_visited, exist := visited_node[start]
	visited_node[start] = make(map[string]bool)
	mutex.Unlock()
	if caching.CheckCacheFile(start) {
		url_scrap = caching.GetCacheUrl(start)
	} else {
		url_scrap = scraper.ScrapWeb(start)
		mutex.Lock()
		caching.SetCacheUrl(start, url_scrap)
		mutex.Unlock()
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
		} else {
			// kalau ada
			mutex.Lock()
			url_list[elmt] = true
			mutex.Unlock()
		}
	}

	for key := range url_list {

		saved_path2 := append(saved_path, key)
		fmt.Println(key, maxdepth)
		if DLSSingle(key, end, maxdepth-1, saved_path2, ans, articleHit) {
			return true
		}
	}

	// wg.Wait()

	return false
}

func IDSReset() {
	for k := range visited_node {
		delete(visited_node, k)
	}
}

func IDShandler(start string, end string, ans_type string) *util.ResponseAPI {
	var resp util.ResponseAPI
	fmt.Println(ans_type)

	var multipath [][]string
	if ans_type == "multi" {
		multipath = IDS(start, end, &resp)

		// for i := range multipath {
		// 	fmt.Println(multipath[i])
		// }
	} else if ans_type == "single" {
		multipath = IDSSingle(start, end, &resp)
	}
	var temp = make(map[string]map[string]bool)
	resp.Nodes, resp.Edges = util.ConvertToVisualizerHandler(start, end, temp, 0, multipath, "IDS")

	return &resp
}
