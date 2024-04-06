package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var items []Item

func main() {
	router := gin.Default()

	// Define API endpoints
	router.GET("/items", getItems)
	router.POST("/items/add", addItem)

	// Start the server
	router.Run(":4000")
}

// Handler for GET /items
func getItems(c *gin.Context) {
	c.JSON(http.StatusOK, items)
}

// Handler for POST /items/add
func addItem(c *gin.Context) {
	var newItem Item
	if err := c.BindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items = append(items, newItem)
	c.JSON(http.StatusCreated, gin.H{"message": "Item added successfully"})
}
