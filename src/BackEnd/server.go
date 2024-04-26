package main

import (
	"backend/algorithm"
	"backend/scraper"
	"backend/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type pathJSON struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	PageStart string `json:"page_start"`
	PageEnd   string `json:"page_end"`
	Method    string `json:"method"`
}

/*
+------------------------------------------+
|               MAIN FUNCTION              |
+------------------------------------------+
*/

// func main() {
// 	router := gin.Default()

// 	router.Use(cors.Default())

// 	// Define API endpoints
// 	router.POST("/api", getPath)

// 	// Start the server
// 	router.Run(":4000")
// }

func main() {
	page1 := scraper.SendApi("Jokowi")
	page2 := scraper.SendApi("Central Java")

	// 	// get initial value
	// 	fmt.Println(PrettyPrint(page1))
	// 	fmt.Println(PrettyPrint(page2))

	// 	// start scraping
	// 	// max_depth := 3
	hasil := algorithm.BFShandler(page1.Url, page2.Url)
	// hasil2 := bfsHandler(page1.Url, page2.Url)
	// 	// fmt.Println(hasil.Message)
	// 	// fmt.Println(hasil.Status)
	fmt.Println(hasil.Time)

	// 	// fmt.Println(hasil.Nodes)
	// 	// fmt.Println(hasil.Edges)

	// x := IDShandler(page1.Url, page2.Url)
	// fmt.Println(x.Time)
}

/*
+------------------------------------------+
|               ROUTE HANDLER              |
+------------------------------------------+
*/

// Handler for GET /api
func getPath(c *gin.Context) {

	var requestData pathJSON

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page1 := scraper.SendApi(requestData.Start)
	page2 := scraper.SendApi(requestData.End)

	if !page1.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page with start value title not found"})
		return
	}

	if !page2.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page with end value title not found"})
		return
	}

	// clean input
	var resp util.ResponseAPI

	fmt.Print(requestData.Method)
	if requestData.Method == "BFS" {
		resp = *algorithm.BFShandler(page1.Url, page2.Url)
	} else if requestData.Method == "IDS" {
		fmt.Print("OIT INI DARI IDS")
		resp = *algorithm.IDShandler(page1.Url, page2.Url)
	}

	c.JSON(http.StatusOK, resp)
}
