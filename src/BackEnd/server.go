package main

import (
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
}

/*
+------------------------------------------+
|               MAIN FUNCTION              |
+------------------------------------------+
*/

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	// Define API endpoints
	router.POST("/api", getPath)

	// Start the server
	router.Run(":4000")
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

	page1 := sendApi(requestData.Start)
	page2 := sendApi(requestData.End)

	if !page1.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page with start value title not found"})
		return
	}

	if !page2.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page with end value title not found"})
		return
	}

	// clean input
	var resp ResponseAPI

	fmt.Print(requestData.Method)
	if requestData.Method == "BFS" {
		resp = *bfsHandler(page1.Url, page2.Url)
	} else if requestData.Method == "IDS" {
		fmt.Print("OIT INI DARI IDS")
		resp = *IDShandler(page1.Url, page2.Url)
	}

	c.JSON(http.StatusOK, resp)
}
