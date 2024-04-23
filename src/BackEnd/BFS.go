package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var visited_bfs = make(map[string]bool)
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

	current_url = url

	for current_url != end {

		go BFS(current_url, end, &resp)
		// time.Sleep(time.Millisecond * randomTime())

		if len(url_queue) > 0 {
			current_url = url_queue[0]
			url_queue = url_queue[1:]
		}
	}

	return &resp
}

func BFS(current_url string, end string, resp *ResponseAPI) {

	fmt.Println(current_url)

	link_res := scrapWeb(current_url)

	for _, elmt := range link_res {
		// check if
		go func() {
			mutex.Lock()
			_, err := visited_bfs[elmt]
			mutex.Unlock()
			if !err {
				// not exist
				mutex.Lock()
				visited_bfs[elmt] = true
				url_queue = append(url_queue, elmt)
				mutex.Unlock()
			}
		}()
	}
}
