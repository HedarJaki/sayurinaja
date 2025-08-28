package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct() {
	/*var body struct {
		Product_name string
		StoreID      int
		Stock        int
		Price        int
		Description  string
	}*/

}

func ShowProduct(c *gin.Context) {
	id := c.Param("id")
	var Product []model.Product
	if err := initializer.DB.Where("storeID = ?", id).Find(&Product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Products": "there is no product in this store",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Products": Product,
	})
}
