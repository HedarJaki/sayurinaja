package controller

import (
	"fmt"
	"math/rand"
	"mobapp/initializer"
	"mobapp/model"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var user struct {
	Username string
	Email    string
	Password string
	Otp      string
	ExpOTP   time.Time
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func HasMX(email string) bool {
	part := strings.Split(email, "@")
	if len(part) != 2 {
		return false
	}

	domain := part[1]
	mx, err := net.LookupMX(domain)
	return err != nil && len(mx) > 0
}

func SignUp(c *gin.Context) {
	if c.Bind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	if !HasMX(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email doesnt exist",
		})
		return
	}
}

func GetOTP(c *gin.Context) {
	otp := GenerateOTP()
	user.Otp = otp
	user.ExpOTP = time.Now().Add(5 * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"kode": otp,
	})
}

func VerifyOTP(c *gin.Context) {
	var otp string
	if c.Bind(&otp) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if user.Otp != otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	if time.Now().After(user.ExpOTP) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code already expired"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := model.User{Username: user.Username, Email: user.Email, Password: string(hash)}
	if initializer.DB.Create(&user).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": "succeed to create account",
	})
	c.Redirect(http.StatusFound, "/Login")
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
	c.Redirect(http.StatusFound, "/home")
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
