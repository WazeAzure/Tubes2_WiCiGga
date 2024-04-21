package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type pathJSON struct {
	Path_start string `json:"path_start"`
	Path_end   string `json:"path_end"`
	Method     string `json:"method"`
}

/*
+------------------------------------------+
|               MAIN FUNCTION              |
+------------------------------------------+
*/

var path = pathJSON{
	Path_start: "dummy_start",
	Path_end:   "dummy_end",
	Method:     "method",
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	// Define API endpoints
	router.GET("/path", getPath)
	router.POST("/path", addPath)

	// Start the server
	router.Run(":4000")
}

/*
+------------------------------------------+
|               ROUTE HANDLER              |
+------------------------------------------+
*/

// Handler for GET /items
func getPath(c *gin.Context) {
	c.JSON(http.StatusOK, path)
}

// Handler for POST /items/add
func addPath(c *gin.Context) {
	var newPath pathJSON
	if err := c.BindJSON(&newPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	path = newPath
	c.JSON(http.StatusCreated, gin.H{"message": "Item added successfully"})
}
