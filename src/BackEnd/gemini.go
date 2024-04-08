package main

import (
    "fmt"
    "sync"

    "github.com/gocolly/colly/v2"
)

type WikiPage struct {
    Title string
    Links []string
}

func scrapeWikiPage(c *colly.Collector, url string, wg *sync.WaitGroup) (WikiPage, error) {
    var page WikiPage

    c.OnHTML("title", func(e *colly.HTMLElement) {
        page.Title = e.Text
    })

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        link := e.Attr("href")
        if link[0] != '#' { // Ignore internal links (e.g., jump to sections)
            page.Links = append(page.Links, link)
        }
    })

    err := c.Visit(url)
    wg.Done()
    return page, err
}

func findShortestPath(startURL, endURL string) ([]string, error) {
    visitedForward := make(map[string]bool)
    visitedBackward := make(map[string]bool)
    queueForward := []string{startURL}
    queueBackward := []string{endURL}
    parentsForward := make(map[string]string)
    parentsBackward := make(map[string]string)
    var wg sync.WaitGroup

    for len(queueForward) > 0 || len(queueBackward) > 0 {
        // Forward BFS
        for i := 0; i < len(queueForward); i++ {
            url := queueForward[i]
            queueForward = queueForward[i+1:]

            if visitedForward[url] {
                continue
            }

            visitedForward[url] = true

            if visitedBackward[url] {
                path := reconstructPath(parentsForward, url, parentsBackward)
                return path, nil
            }

            wg.Add(1)
            go func(url string) {
                page, err := scrapeWikiPage(c, url, &wg)
                if err != nil {
                    fmt.Println(err)
                    return
                }

                for _, link := range page.Links {
                    if !visitedForward[link] {
                        queueForward = append(queueForward, link)
                        parentsForward[link] = url
                    }
                }
            }(url)
        }

        // Backward BFS
        for i := 0; i < len(queueBackward); i++ {
            url := queueBackward[i]
            queueBackward = queueBackward[i+1:]

            if visitedBackward[url] {
                continue
            }

            visitedBackward[url] = true

            if visitedForward[url] {
                path := reconstructPath(parentsBackward, url, parentsForward)
                return path, nil
            }

            wg.Add(1)
            go func(url string) {
                page, err := scrapeWikiPage(c, url, &wg)
                if err != nil {
                    fmt.Println(err)
                    return
                }

                for _, link := range page.Links {
                    if !visitedBackward[link] {
                        queueBackward = append(queueBackward, link)
                        parentsBackward[link] = url
                    }
                }
            }(url)
        }

        wg.Wait()
    }

    return nil, fmt.Errorf("No path found from %s to %s", startURL, endURL)
}

func reconstructPath(parents map[string]string, current string, otherParents map[string]string) []string {
    path := []string{current}
    for current != ""; current = parents[current] {
        path = append(path, current)
    }

    reversedPath := []string{}
    for i := len(path) - 1; i >= 0; i-- {
        reversedPath = append(reversedPath, path[i])
    }

    for i := 0; i < len(reversedPath); i++ {
        if otherParents[reversedPath[i]] != "" {
            return append(reversedPath[:i+1], reversePath(otherParents, otherParents[reversedPath[i]])...)
        }
    }

    return nil
}

func reversePath(parents map[string]string, current string) []string {
	path := []string{}
	for current != "" {
		path = append(path, current)
		current = parents[current]
	}
 
	reversedPath := []string{}
	for i := len(path) - 1; i >= 0; i-- {
		reversedPath = append(reversedPath, path[i])
	}
 
	return reversedPath
 }