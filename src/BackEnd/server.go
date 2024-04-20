package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ResponseJSON struct {
	Path   []string `json:"path"`
	Status bool     `json:"status"`
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
	var resp ResponseJSON
	resp.Path = []string{"Nigga-chan", "Nigga-balls", "Ching-chong"}
	resp.Status = true
	c.JSON(http.StatusOK, resp)
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
