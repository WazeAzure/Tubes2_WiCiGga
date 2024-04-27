package scraper

import (
	"backend/caching"
	"backend/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// var (
// 	wg    sync.WaitGroup
// 	mutex sync.Mutex
// )

var semaphore_scraping = make(chan struct{}, 10)

func formatParams(params map[string]string) string {
	queryString := "?"
	for key, value := range params {
		queryString += key + "=" + url.QueryEscape(value) + "&"
	}
	return queryString[:len(queryString)-1]
}

/*
* Fungsi untuk searching dan get recommendation
 */
func SendApi(search string) util.Page {

	apiUrl := "https://en.wikipedia.org/w/api.php"
	webUrl := "https://en.wikipedia.org/wiki/"

	params := map[string]string{
		"action":        "query",
		"format":        "json",
		"formatversion": "2",
		"list":          "search",
		"srsearch":      search,
	}

	var ansPage util.Page

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Create an HTTP request
	req, err := http.NewRequest("GET", apiUrl+formatParams(params), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		ansPage.Status = false
		return ansPage
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "WikiRaceSolverGoQuery/1.1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		ansPage.Status = false
		return ansPage
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		ansPage.Status = false
		return ansPage
	}
	defer resp.Body.Close()

	var result util.AutoGenerated
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if len(result.Query.Search) == 0 {
		ansPage.Status = false
		return ansPage
	}

	ansPage.Pageid = result.Query.Search[0].Pageid
	ansPage.Title = result.Query.Search[0].Title
	ansPage.Url = webUrl + strings.ReplaceAll(ansPage.Title, " ", "_")
	ansPage.Status = true

	return ansPage
}

func ExtractTitle(n *html.Node) string {
	// look until HEAD done. dont go to <body>
	// if n.Type == html.ElementNode && n.Data == "body" {
	// 	return ""
	// }

	// if n.Type == html.ElementNode && n.Data == "title" {
	// If the node is a <title> element, return its text content
	// str := strings.TrimSpace(textContent(n))
	// fmt.Println("TITLE FOUND ================ ", str[6:len(str)-12])

	// if len(str)-12 <= 3 {
	// 	fmt.Println("=========================================")
	// 	fmt.Println("|              ", str, "                  |")
	// 	fmt.Println("=========================================")
	// }
	// return str[:len(str)-12]
	// }

	// Recursively search for <title> element in child nodes
	// for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 	title := ExtractTitle(c)
	// 	if title != "" {
	// 		return title
	// 	}
	// }

	doc := goquery.NewDocumentFromNode(n)
	title := doc.Find("title").Text()
	title = strings.TrimSpace(title)
	return title[:len(title)-12]
}

// func textContent(n *html.Node) string {
// 	var text string
// 	if n.Type == html.TextNode {
// 		text = n.Data
// 	}
// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		text += textContent(c)
// 	}
// 	return text
// }

func handleRedirect(url string) string {

	// defer func() { <-semaphore_scraping }()
	fmt.Println("CALLED ", url, " ==================")
	webUrl := "https://en.wikipedia.org"
	resp, err := http.Get(webUrl + url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body response", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal("Error parsing HTML", err)
	}

	title := ExtractTitle(doc)

	// var temp Page
	// fmt.Println(url, "===========")
	// temp = sendApi(url[6:])

	return "/wiki/" + strings.ReplaceAll(title, " ", "_")
}

func customFilterURL(url string) bool {
	a := strings.HasPrefix(url, "/wiki/")
	if a {
		namespace_list := []string{"User:", "File:", "MediaWiki:", "Template:", "Help:", "Category:", "Special:", "Talk:", "Template_talk:", "Wikipedia:", "Main_Page", "#", "Portal:"}

		for _, elmt := range namespace_list {
			a = a && !strings.Contains(url, elmt)
		}
	}
	return a
}

// Function to extract links from a node
func ExtractLinks(n *html.Node, links *[]string, wg *sync.WaitGroup) {
	// stack := []*html.Node{}

	// var stackNotEmpty int32

	// // Atomic operations to set and check the stackNotEmpty flag
	// setStackNotEmpty := func() { atomic.StoreInt32(&stackNotEmpty, 1) }
	// isStackNotEmpty := func() bool { return atomic.LoadInt32(&stackNotEmpty) == 1 }

	// push := func(node *html.Node) {
	// 	mutex.Lock()
	// 	defer mutex.Unlock()
	// 	stack = append(stack, node)
	// 	setStackNotEmpty()
	// }

	// pop := func() (node *html.Node) {
	// 	mutex.Lock()
	// 	defer mutex.Unlock()
	// 	if len(stack) == 0 {
	// 		return nil
	// 	}
	// 	node = stack[len(stack)-1]
	// 	stack = stack[:len(stack)-1]
	// 	if len(stack) == 0 {
	// 		atomic.StoreInt32(&stackNotEmpty, 0)
	// 	}
	// 	return node
	// }

	// processNode := func(current_node *html.Node) {
	// 	// fmt.Println("CALLED PROCESS NODE")
	// 	var temp_link []string

	// 	if current_node.Type == html.ElementNode && current_node.Data == "a" {
	// 		url := ""
	// 		isTrue := false
	// 		for _, attr := range current_node.Attr {
	// 			if attr.Key == "href" {
	// 				if customFilterURL(attr.Val) {
	// 					url = attr.Val
	// 					isTrue = true
	// 					break
	// 				}
	// 			}
	// 		}
	// 		if isTrue {
	// 			for _, attr := range current_node.Attr {
	// 				if attr.Key == "class" {
	// 					if attr.Val == "mw-redirect" {

	// 						semaphore_scraping <- struct{}{}

	// 						if caching.CheckCacheRedirect(url) {
	// 							url = caching.GetCacheRedirect(url)
	// 						} else {
	// 							url2 := url
	// 							url = handleRedirect(url)
	// 							caching.SetCacheRedirect(url2, url)
	// 						}
	// 						break
	// 					}
	// 				}
	// 			}
	// 		}
	// 		if url != "" {
	// 			temp_link = append(temp_link, "https://en.wikipedia.org"+url)
	// 		}
	// 	}

	// 	// add child node to stack
	// 	for c := current_node.FirstChild; c != nil; c = c.NextSibling {
	// 		// fmt.Println("called for loop ========")
	// 		push(c)
	// 	}

	// 	// update *link
	// 	mutex.Lock()
	// 	defer mutex.Unlock()
	// 	*links = append(*links, temp_link...)
	// 	defer wg.Done()
	// }

	// push(n)

	// for {
	// 	if !isStackNotEmpty() {
	// 		break
	// 	}

	// 	semaphore_scraping <- struct{}{}

	// 	node := pop()

	// 	if node == nil {
	// 		<-semaphore_scraping
	// 		continue
	// 	}

	// 	wg.Add(1)
	// 	go func(node *html.Node) {
	// 		defer func() { <-semaphore_scraping }()
	// 		processNode(node)
	// 	}(node)

	// 	time.Sleep(50 * time.Millisecond)
	// }

	// wg.Wait()

	doc := goquery.NewDocumentFromNode(n)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		href, exist := s.Attr("href")
		if exist {
			if customFilterURL(href) {
				url := href
				class, exist2 := s.Attr("class")
				if exist2 {
					if class == "mw-redirect" {
						if caching.CheckCacheRedirect(url) {
							url = caching.GetCacheRedirect(url)
						} else {
							// semaphore_scraping <- struct{}{}
							url2 := url
							url = handleRedirect(url)
							caching.SetCacheRedirect(url2, url)
						}
					}
				}

				*links = append(*links, "https://en.wikipedia.org"+url)
			}
		}
		// }()
	})

	wg.Wait()
}

func findLinksInContent(n *html.Node) []string {
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "main" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "content" {
					// semaphore := make(chan struct{}, 10)
					var wait sync.WaitGroup
					// If found node with class "content", extract links
					// semaphore <- struct{}{}
					// go func() {
					// defer func() { <-semaphore }()
					// defer wg.Done()
					ExtractLinks(n, &links, &wait)
					// }()
					return
				}
			}
		}

		// Recursively traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	return links
}

func ScrapWeb(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body response", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal("Error parsing HTML", err)
	}

	// find all links
	final_ans := findLinksInContent(doc)

	// namespace_list := []string{"User:", "File:", "MediaWiki:", "Template:", "Help:", "Category:", "Special:", "Talk:", "Template_talk:", "Wikipedia:"}

	// urlHTML := doc.Find("#content").Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
	// 	link, _ := s.Attr("href")

	// 	a := strings.HasPrefix(link, "/wiki/")

	// 	for _, elmt := range namespace_list {
	// 		a = a && !strings.Contains(link, elmt)
	// 	}
	// 	return a
	// })

	// final_ans := urlHTML.Map(func(i int, s *goquery.Selection) string {
	// 	link, _ := s.Attr("href")

	// 	return "https://en.wikipedia.org" + link
	// })

	return final_ans
}
