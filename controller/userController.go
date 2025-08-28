package controller

import (
	"mobapp/initializer"
	"mobapp/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var body struct {
		Username string
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := model.User{Username: body.Username, Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": "succeed to create account",
	})

}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	var user model.User
	initializer.DB.First(&user, "email = ?", body.Email)

	if user.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 3600*24, "", "", true, true)
	c.JSON(http.StatusOK, gin.H{})
}

func HomePage(c *gin.Context) {
	var toko []model.Store
	var produk []model.Product

	if err := initializer.DB.Find(&toko).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to find store",
		})
		return
	}
	if len(toko) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"store":   toko,
			"message": "no store found in this area",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"store": toko,
	})

	if err := initializer.DB.Find(&produk).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to find product",
		})
		return
	}
	if len(toko) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"product": produk,
			"message": "no product",
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"product": produk,
	})
}
