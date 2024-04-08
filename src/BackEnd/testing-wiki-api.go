package main

import (
	"fmt"
	"net/http"
)

func formatParams(params map[string]string) string {
	queryString := "?"
	for key, value := range params {
		queryString += key + "=" + value + "&"
	}
	return queryString[:len(queryString)-1]
}

func sendApi(search string) {
	search := "Jokowi"
	apiUrl := "https://en.wikipedia.org/w/api.php"

	params := map[string]string{
		"action":        "parse",
		"format":        "json",
		"formatversion": "2",
		"list":          "search",
		"srsearch":      search,
	}

	fmt.Println(formatParams(params))
	fmt.Println(searchTitle)
	resp, err := http.Get(apiUrl + formatParams(params))
	fmt.Println(resp)
	fmt.Println(err)

}

func main() {
	fmt.Println("hello world")

	sendApi()
}
