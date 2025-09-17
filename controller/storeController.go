package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"strconv"

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

func ShowstorebyID(id int) (model.Store, error) {
	var store model.Store
	if err := initializer.DB.Where("storeID = ?", id).First(&store).Error; err != nil {
		return store, err
	}
	return store, nil
}

func CheckStore(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get store id",
		})
		return
	}

	storeData, err := ShowstorebyID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get store data",
		})
		return
	}

	product, err := ShowProductbyID(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get product",
		})
		return
	}

	response := gin.H{
		"store":    storeData,
		"products": product,
	}
	if len(product) == 0 {
		response["message"] = "there is no product in this store"
	}

	c.JSON(http.StatusOK, response)
}

func UserStoreHomePage(c *gin.Context) {
	storeID := c.GetInt("storeID")

	var store model.Store
	result := initializer.DB.Preload("Product").Preload("Order").First(&store, "storeID = ?", storeID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get your customers order",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "store not found",
		})
		return
	}

	response := gin.H{
		"name":        store.StoreName,
		"rating":      store.Rating,
		"description": store.StoreDesription,
		"address":     store.StoreAddress,
		"product":     store.Products,
		"order":       store.Orders,
	}

	if len(store.Products) == 0 {
		response["message"] = "there is no product in your store"
	}

	if len(store.Orders) == 0 {
		response["message"] = "there is no order in your store"
	}

	c.JSON(http.StatusOK, response)
}

func StoreReview(c *gin.Context) {
	var body struct {
		star int
		desc string
	}
	user := c.GetInt("userID")
	storeid, err := strconv.Atoi(c.Param("id"))
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

	review := model.StoreReview{
		StoreID:           storeid,
		UserID:            user,
		Star:              body.star,
		Store_review_desc: body.desc,
	}

	if initializer.DB.Create(&review).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create review",
		})
		return
	}

	var store model.Store
	if initializer.DB.Preload("Review").Where("storeId = ?", store).First(&store).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed load store's review",
		})
		return
	}

	if len(store.Reviews) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "no review",
		})
		return
	}

	var totalstar int
	for _, review := range store.Reviews {
		totalstar += review.Star
	}

	store.Rating = float64(totalstar) / float64(len(store.Reviews))
	if initializer.DB.Save(&store).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update the rating",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "you review created successfully"})
}
