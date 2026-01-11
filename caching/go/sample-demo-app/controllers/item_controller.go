package controllers

import (
	"net/http"
	"strconv"

	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/models"
	"github.com/gin-gonic/gin"
)

type ItemController struct {
	DataCtrl *DataController
}

func NewItemController() *ItemController {
	return &ItemController{
		DataCtrl: NewDataController(),
	}
}

func (ic *ItemController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := ic.DataCtrl.Get(id)
	if err != nil {
		// Simple error handling, could distinguish Not Found vs Internal Error
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (ic *ItemController) GetRandom(c *gin.Context) {
	items, err := ic.DataCtrl.GetMultiple(10) // Matches Java's TOTAL_RANDOM_ITEMS
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch random items"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (ic *ItemController) Create(c *gin.Context) {
	var newItem models.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := ic.DataCtrl.Create(&newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (ic *ItemController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = ic.DataCtrl.Delete(id)
	if err != nil {
		// Java sample doesn't return error on delete failure usually, creates valid JSON response
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
