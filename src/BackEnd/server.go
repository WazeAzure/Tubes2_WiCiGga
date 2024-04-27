package main

import (
	"backend/algorithm"
	"backend/caching"
	"backend/scraper"
	"backend/util"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type pathJSON struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	PageStart string `json:"page_start"`
	PageEnd   string `json:"page_end"`
	Method    string `json:"method"`
	Type      string `json:"type"`
}

/*
+------------------------------------------+
|               MAIN FUNCTION              |
+------------------------------------------+
*/

func main() {
	// initialize cache - PASTIKAN SELALU DIPANGGIL
	caching.CheckCacheFolder()
	caching.InitCache()

	/* SERVER PART */
	router := gin.Default()

	router.Use(cors.Default())

	// Define API endpoints
	router.POST("/api", getPath)

	// Start the server
	router.Run(":4000")

	/* TESTING PART */
	// page1 := scraper.SendApi("Jokowi")
	// page2 := scraper.SendApi("Central Java")

	// 	// get initial value
	// 	fmt.Println(PrettyPrint(page1))
	// 	fmt.Println(PrettyPrint(page2))

	// 	// start scraping
	// 	// max_depth := 3
	// hasil := algorithm.BFShandler(page1.Url, page2.Url, "single")
	// hasil2 := bfsHandler(page1.Url, page2.Url)
	// 	// fmt.Println(hasil.Message)
	// 	// fmt.Println(hasil.Status)
	// fmt.Println(hasil.Time)

	// 	// fmt.Println(hasil.Nodes)
	// 	// fmt.Println(hasil.Edges)

	// x := algorithm.IDShandler(page1.Url, page2.Url, "single")
	// fmt.Println(x.Time)
}

/*
+------------------------------------------+
|               ROUTE HANDLER              |
+------------------------------------------+
*/

// Handler for GET /api
func getPath(c *gin.Context) {

	// clean global var
	algorithm.IDSReset()
	algorithm.BFSReset()

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

	fmt.Println(requestData.Method)
	if requestData.Method == "BFS" {
		resp = *algorithm.BFShandler(page1.Url, page2.Url, requestData.Type)
	} else if requestData.Method == "IDS" {
		resp = *algorithm.IDShandler(page1.Url, page2.Url, requestData.Type)
	}

	c.JSON(http.StatusOK, resp)
}
