package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	cartID := c.GetInt("cartID")
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get product id",
		})
		return
	}

	var body struct {
		quantity int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	var price float64
	initializer.DB.Table("products").Select("price_each").Where("productID = ?", productID).First(&price)

	cartItem := model.CartItem{CartID: cartID, ProductID: productID, Quantity: body.quantity, Price: float64(body.quantity) * price}
	if initializer.DB.Create(&cartItem).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product added",
	})
}

func UpdateCart(c *gin.Context) {
	cart := c.GetInt("cartID")

	var body struct {
		productID int
		quantity  int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	var cartItem model.CartItem
	if body.quantity == 0 {
		if initializer.DB.Where("productID = ? AND cartID = ?", body.productID, cart).Delete(&cartItem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to remove product from cart",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "cart updated successfully",
		})
	} else if cartItem.Quantity == body.quantity {
		cartItem.Price = float64(body.quantity) * (cartItem.Price / float64(cartItem.Quantity))
		cartItem.Quantity = body.quantity
		if initializer.DB.Where("productID = ? AND cartID = ?", body.productID, cart).Save(&cartItem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update the cart",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "cart updated successfully",
		})
	}
}

func ShowCart(c *gin.Context) {
	cartID := c.GetInt("cartID")

	type body struct {
		name     string
		quantity int
	}
	var items []body
	if initializer.DB.Table("cartItems").Select("products.product_name, cartItems.quantity").Joins("products ON cartItem.productID = products.productID").Where("cartItems.cartID = ?", cartID).Find(&items).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to find item",
		})
		return
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"eror": "cart is empty"})
		return
	}

	for _, item := range items {
		c.JSON(http.StatusOK, gin.H{
			"product":  item.name,
			"quantity": item.quantity,
		})
	}
}

func AddOrder(c *gin.Context) {
	user := c.GetInt("userID")
	var cart model.Cart
	if initializer.DB.Preload("Item").Where("userID = ? AND is_active = ?", user, true).First(&cart).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no cart is active"})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
	}

	var product model.Product
	if initializer.DB.Where("productID = ?", cart.CartItems[0].ProductID).First(&product).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant find the store you ordered"})
		return
	}

	tx := initializer.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total float64
	for _, item := range cart.CartItems {
		total += item.Price
	}

	order := model.Order{
		UserID:      user,
		StoreID:     product.StoreID,
		Total_price: total,
		Status:      "dalam proses",
	}

	if tx.Create(&order).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	for _, item := range cart.CartItems {
		orderitem := model.OrderDetail{
			ProductID: item.ProductID,
			OrderID:   order.OrderID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if tx.Create(&orderitem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order item"})
			return
		}
	}

	cart.Is_active = false
	if tx.Save(&cart).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faied to close cart"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "order created successfully"})
}
