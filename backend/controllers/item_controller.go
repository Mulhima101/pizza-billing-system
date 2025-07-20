package controllers

import (
    "net/http"
    "pizza-billing-backend/models"
    "pizza-billing-backend/services"
    "strconv"

    "github.com/gin-gonic/gin"
)

type ItemController struct {
    itemService *services.ItemService
}

func NewItemController(itemService *services.ItemService) *ItemController {
    return &ItemController{itemService: itemService}
}

func (ic *ItemController) GetItems(c *gin.Context) {
    items, err := ic.itemService.GetAllItems()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, items)
}

func (ic *ItemController) GetItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    item, err := ic.itemService.GetItemByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
        return
    }
    c.JSON(http.StatusOK, item)
}

func (ic *ItemController) CreateItem(c *gin.Context) {
    var item models.Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdItem, err := ic.itemService.CreateItem(&item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, createdItem)
}

func (ic *ItemController) UpdateItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    var item models.Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    item.ID = id

    updatedItem, err := ic.itemService.UpdateItem(&item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updatedItem)
}

func (ic *ItemController) DeleteItem(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    err = ic.itemService.DeleteItem(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func (ic *ItemController) GetCategories(c *gin.Context) {
    categories, err := ic.itemService.GetAllCategories()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, categories)
}