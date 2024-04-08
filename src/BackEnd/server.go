package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
+------------------------------------------+
|               MAIN FUNCTION              |
+------------------------------------------+
*/

func main() {
	router := gin.Default()

	// Define API endpoints
	router.GET("/path", getPath)

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
	c.JSON(http.StatusOK, "this is homepage")
}

// Handler for POST /items/add
// func addItem(c *gin.Context) {
// 	var newItem Item
// 	if err := c.BindJSON(&newItem); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	items = append(items, newItem)
// 	c.JSON(http.StatusCreated, gin.H{"message": "Item added successfully"})
// }
