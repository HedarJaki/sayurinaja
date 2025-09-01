package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
	val, exits := c.Get("store")
	if !exits {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "cant find your store",
		})
		return
	}
	store := val.(model.Store)
	var body struct {
		Product_name string
		category     string
		Stock        int
		Price        int
		Description  string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}
	product := model.Product{Product_name: body.Product_name, Category: body.category, StoreID: store.StoreID, Stock: body.Stock, Price: body.Price, Description: body.Description}
	if initializer.DB.Create(&product).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create product",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": "Product created succeessfully",
	})
}

func ShowProductbyID(id int) ([]model.Product, error) {
	var Product []model.Product
	if err := initializer.DB.Where("storeID = ?", id).Find(&Product).Error; err != nil {
		return nil, err
	}
	return Product, nil
}

func ShowProductByCategory(c *gin.Context) {
	category := c.Query("category")
	var product []model.Product
	if initializer.DB.Where("category = ?", category).Find(&product).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to find product",
		})
		return
	}

	if len(product) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"Products": product,
			"message":  "there is no product in this category",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Products": product,
	})
}
