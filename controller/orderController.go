package controller

import (
	"fmt"
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var total_cart float64

func AddToCart(c *gin.Context) {
	cartID := c.GetInt("cartID")
	//productID, err := strconv.Atoi(c.Param("id"))

	var body struct {
		Item []struct {
			ProductID int
			Quantity  int
			Note      string
		}
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	for _, item := range body.Item {
		var Price_Each float64
		if initializer.DB.Table("products").Where("productID = ?", item.ProductID).First(&Price_Each).Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %d not found", item.ProductID)})
			return
		}

		cartitem := model.CartItem{
			CartID:    cartID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Note:      item.Note,
		}
		if initializer.DB.Create(&cartitem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to add the item",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product added",
	})
}

/*func UpdateCart(c *gin.Context) {
	cart := c.GetInt("cartID")

	var body struct {
		ProductID int
		Quantity  int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	var cartItem model.CartItem
	if body.Quantity == 0 {
		if initializer.DB.Where("productID = ? AND cartID = ?", body.ProductID, cart).Delete(&cartItem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to remove product from cart",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "cart updated successfully",
		})
	} else if cartItem.Quantity == body.Quantity {
		cartItem.Price = float64(body.Quantity) * (cartItem.Price / float64(cartItem.Quantity))
		cartItem.Quantity = body.Quantity
		if initializer.DB.Where("productID = ? AND cartID = ?", body.ProductID, cart).Save(&cartItem).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update the cart",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "cart updated successfully",
		})
	}
}*/

func ShowCart(c *gin.Context) {
	cartID := c.GetInt("cartID")
	shipment := c.Param("shipment")

	var cart model.Cart
	if initializer.DB.Preload("CartItems").First(&cart, "cartID = ? AND is_active = ?", cartID, true) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to find item",
		})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"eror": "cart is empty"})
		return
	}

	var shipment_price float64
	if shipment == "express" {
		shipment_price = 2500
	} else {
		shipment_price = 5000
	}

	for _, item := range cart.CartItems {
		total_cart += item.Price
	}

	total_cart += shipment_price

	c.JSON(http.StatusOK, gin.H{
		"cart":     cart,
		"shipment": shipment_price,
		"total":    total_cart,
	})
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

	order := model.Order{
		UserID:      user,
		StoreID:     product.StoreID,
		Total_price: total_cart,
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

func UpdateOrdershipment(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get order data",
			})
			return
		}
	}

	var order model.Order
	if initializer.DB.First(&order, "orderID = ?", orderID).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get order data",
		})
		return
	}

	order.Status = "sedang dikirim"
	if initializer.DB.Save(&order).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update order",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "order updated successfully",
	})
}

func Checkorder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get orderid",
		})
		return
	}

	var order model.Order
	if initializer.DB.Preload("OrderDetails").Where("orderID = ?", orderID).First(&order).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get order data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func AddAddress(c *gin.Context) {
	Id := c.GetInt("userID")

	var body struct {
		Fulladdress string
		Latitude    float64
		Longitude   float64
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	Address := model.Address{
		UserID:      Id,
		Fulladdress: body.Fulladdress,
		Latitude:    body.Latitude,
		Longitude:   body.Longitude,
	}

	if initializer.DB.Create(&Address).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add your address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "address added successfully",
	})
}
