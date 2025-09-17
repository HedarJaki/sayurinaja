package middleware

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Auth")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if claim, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claim["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			var user model.User
			initializer.DB.First(&user, claim["sub"])

			if user.UserID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("userID", user.UserID)

			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func AuthStoreMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		Id, err := c.Get("userID")
		if !err {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var store model.Store
		if err := initializer.DB.Where("userID = ?", Id).First(&store).Error; err == nil {
			c.Set("storeID", store.StoreID)
		}
		c.Next()
	}
}

func CartCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetInt("userID")
		var cart model.Cart
		if initializer.DB.Where("userID = ? AND is_active = ?", user, true).First(&cart).Error != nil {
			cart = model.Cart{UserID: user, Is_active: true}
			if initializer.DB.Create(&cart).Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to create cart",
				})
				return
			}
		}
		c.Set("cartID", cart.CartID)
		c.Next()
	}
}
