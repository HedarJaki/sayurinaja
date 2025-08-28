package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignStore(c *gin.Context) {
	userID := c.GetInt("userID")

	var body struct {
		StoreName        string
		StoreDescription string
		StoreAddress     string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	var product []model.Product
	if val, isexist := c.Get("store"); isexist {
		exitStore := val.(model.Store)
		initializer.DB.Where("storeId = ?", exitStore.StoreID).Find(&product)
		c.JSON(http.StatusOK, gin.H{
			"store":   exitStore,
			"product": product,
		})
		return
	}

	tx := initializer.DB.Begin()

	store := model.Store{
		UserID:          userID,
		StoreName:       body.StoreName,
		StoreDesription: body.StoreDescription,
		StoreAddress:    body.StoreAddress,
	}

	if tx.Create(&store).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create store",
		})
		return
	}

	var user model.User
	if tx.First(&user, userID).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	user.Is_seller = true
	if tx.Save(&user).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to change status",
		})
		return
	}
	tx.Commit()
	tx.Rollback()
	c.JSON(http.StatusOK, gin.H{
		"message": "store created successfully",
	})
	c.JSON(http.StatusOK, gin.H{
		"store":   store,
		"product": product,
	})
}
