package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Node struct {
	Id    int    `json:"id"`
	Level int    `json:"level"`
	Label string `json:"label"`
	Color string `json:"color"`
}

type Edge struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type ResponseAPI struct {
	Nodes   []Node        `json:"nodes"`
	Edges   []Edge        `json:"edges"`
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Time    time.Duration `json:"time"`
}

type Page struct {
	Title  string `json:"title"`
	Pageid int    `json:"pageid"`
	Url    string `json:"url"`
	Status bool   `json:"status"`
}

type AutoGenerated struct {
	Batchcomplete bool `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Searchinfo struct {
			Totalhits         int    `json:"totalhits"`
			Suggestion        string `json:"suggestion"`
			Suggestionsnippet string `json:"suggestionsnippet"`
		} `json:"searchinfo"`
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Size      int       `json:"size"`
			Wordcount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

var colorlist = []string{"#ce76fe", "#ed6f71", "#feae33", "#5358e2", "#fec46b"}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func formatParams(params map[string]string) string {
	queryString := "?"
	for key, value := range params {
		queryString += key + "=" + url.QueryEscape(value) + "&"
	}
	return queryString[:len(queryString)-1]
}

func isIn(s string, arr []string) bool {
	for _, elmt := range arr {
		if elmt == s {
			return true
		}
	}
	return false
}

/*
* Fungsi untuk searching dan get recommendation
 */
func sendApi(search string) Page {

	apiUrl := "https://en.wikipedia.org/w/api.php"
	webUrl := "https://en.wikipedia.org/wiki/"

	params := map[string]string{
		"action":        "query",
		"format":        "json",
		"formatversion": "2",
		"list":          "search",
		"srsearch":      search,
	}

	var ansPage Page

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

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		ansPage.Status = false
		return ansPage
	}

	var result AutoGenerated
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

// func toFile(res []byte) {
// 	os.WriteFile("result", res, 0644)
// }

// func toFileS(res []string) {
// 	f, err := os.Create("result")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	w := bufio.NewWriter(f)

// 	temp_string := ""

// 	for _, element := range res {
// 		temp_string = temp_string + "\n" + element
// 	}

// 	n4, err := w.WriteString(temp_string)
// 	fmt.Println(n4)
// 	w.Flush()
// }

func timeTrack(start time.Time, name string, executedTime *time.Duration) {
	elapsed := time.Since(start)
	*executedTime = elapsed
	log.Printf("%s took %s", name, elapsed)

}

func customFileterURL(url string) bool {
	a := strings.HasPrefix(url, "/wiki/")
	if a {
		namespace_list := []string{"User:", "File:", "MediaWiki:", "Template:", "Help:", "Category:", "Special:", "Talk:", "Template_talk:", "Wikipedia:", "Main_Page", "#"}

		for _, elmt := range namespace_list {
			a = a && !strings.Contains(url, elmt)
		}
	}
	return a
}

// Function to extract links from a node
func extractLinks(n *html.Node, links *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				if customFileterURL(attr.Val) {
					*links = append(*links, "https://en.wikipedia.org"+attr.Val)
				}
			}
		}
	}

	// Recursively extract links from child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, links)
	}
}

func findLinksInContent(n *html.Node) []string {
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "main" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "content" {
					// If found node with class "content", extract links
					extractLinks(n, &links)
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

func scrapWeb(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read body response", err)
	}

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

func convertToVisualizer(id *int, depth int, temp_nodes *[]Node, temp_edges *[]Edge, temp_visited map[string]int, graph map[string]map[string]bool, url string, start string, maxdepth int) bool {
	if url == start {
		// create edge
		return true
	}

	if depth > maxdepth {
		// fmt.Println(depth, " ", maxdepth, " - depth | max depth")
		return false
	}

	fmt.Println("current url: ", url)
	fmt.Println("current depth: ", depth)

	_, err := temp_visited[url]
	if !err {
		// not exist in visited
		temp_visited[url] = *id

		// create Node
		fmt.Println("========= NODE CREATED =========")
		*temp_nodes = append(*temp_nodes, Node{Id: *id, Level: depth, Label: url[30:], Color: colorlist[depth]})
		*id++
	}

	for key, _ := range graph[url] {

		// create edges
		fmt.Println("url: ", key)

		// time.Sleep(1 * time.Second)
		// traverse
		x := convertToVisualizer(id, depth+1, temp_nodes, temp_edges, temp_visited, graph, key, start, maxdepth)
		if x {
			// create edge
			var from_int int = temp_visited[key]
			var to_int int = temp_visited[url]
			*temp_edges = append(*temp_edges, Edge{From: from_int, To: to_int})
		}
	}

	return true
}

func convertToVisualizerHandler(start string, end string, temp map[string]map[string]bool, maxdepth int) ([]Node, []Edge) {

	var temp_nodes []Node
	var temp_edges []Edge

	var temp_visited = make(map[string]int)
	id := 0
	convertToVisualizer(&id, 0, &temp_nodes, &temp_edges, temp_visited, temp, end, start, maxdepth)

	fmt.Println(temp_nodes)
	return temp_nodes, temp_edges
}

// func main() {
// 	page1 := sendApi("Jokowi")
// 	page2 := sendApi("Central Java")

// 	// get initial value
// 	fmt.Println(PrettyPrint(page1))
// 	fmt.Println(PrettyPrint(page2))

// 	// start scraping
// 	// max_depth := 3
// 	hasil := bfsHandler(page1.Url, page2.Url)
// 	fmt.Println(hasil.Message)
// 	fmt.Println(hasil.Status)
// 	fmt.Println(hasil.Time)

// 	fmt.Println(hasil.Nodes)
// 	fmt.Println(hasil.Edges)

// 	// x := IDS(page1.Url, page2.Url)
// 	// fmt.Println(x)
// }
