package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"strconv"

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
		Price        float64
		Description  string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}
	product := model.Product{Product_name: body.Product_name, Category: body.category, StoreID: store.StoreID, Stock: body.Stock, Price_Each: body.Price, Product_description: body.Description}
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

func ProductReview(c *gin.Context) {
	var body struct {
		star int
		desc string
	}
	user := c.GetInt("userID")
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get store id",
		})
		return
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	review := model.Productreview{
		ProductID:   productId,
		UserID:      user,
		Star:        body.star,
		Description: body.desc,
	}

	if initializer.DB.Create(&review).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create review",
		})
		return
	}

	/*var product model.Product
	if initializer.DB.Preload("Review").Where("productId = ?", productId).First(&product).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed load store's review",
		})
		return
	}

	if len(product.Review) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "no review",
		})
		return
	}

	var totalstar int
	for _, review := range product.Review {
		totalstar += review.Star
	}

	product.Rating = float64(totalstar) / float64(len(product.Review))
	if initializer.DB.Save(&product).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update the rating",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "you review created successfully"})*/
}
