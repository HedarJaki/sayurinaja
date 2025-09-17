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
}

var OTPcode struct {
	Otp    string
	ExpOTP time.Time
}

var account model.User

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func GenerateOTPforgotpw() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HasMX(email string) bool {
	part := strings.Split(email, "@")
	if len(part) != 2 {
		return false
	}

	domain := part[1]
	mx, err := net.LookupMX(domain)
	return err == nil && len(mx) > 0
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
	otp := GenerateOTP()
	OTPcode.Otp = otp
	OTPcode.ExpOTP = time.Now().Add(5 * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"kode": otp,
	})
}

func VerifyOTP(c *gin.Context) {
	var body struct {
		Otp string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if body.Otp != OTPcode.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	if time.Now().After(OTPcode.ExpOTP) {
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
	if initializer.DB.Where("email = ?", body.Email).First(&user).Error != nil {
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
	c.JSON(http.StatusOK, gin.H{"message": "login successfull"})
}

func Forgotpw(c *gin.Context) {
	var body struct {
		Email string `json:"Email" form:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	if initializer.DB.First(&account, "email = ?", body.Email).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	OTP := GenerateOTPforgotpw()
	OTPcode.Otp = OTP
	OTPcode.ExpOTP = time.Now().Add(5 * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"kode": OTP,
	})
}

func VerifyForgotPWOTP(c *gin.Context) {
	var body struct {
		Code string `json:"Code"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	if body.Code != OTPcode.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	if time.Now().After(OTPcode.ExpOTP) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code already expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verify success"})
}

func NewPassword(c *gin.Context) {
	var body struct {
		NewPW     string `json:"Password"`
		ConfirmPW string `json:"Confirmed_password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to load body",
		})
		return
	}

	if body.NewPW != body.ConfirmPW {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "confirm password must be the same as your new password",
		})
		return
	}

	Hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPW), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	account.Password = string(Hash)
	if initializer.DB.Save(account).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to change password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your password already changed ",
	})

}

func Home(c *gin.Context) {

	var store []model.Store
	var product []model.Product

	if err := initializer.DB.Find(&store).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to find store",
		})
		return
	}

	if err := initializer.DB.Find(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to find product",
		})
		return
	}

	resp := gin.H{
		"toko":   store,
		"produk": product,
	}

	if len(store) == 0 {
		resp["store_message"] = "no store found in this area"
	}

	if len(product) == 0 {
		resp["product_message"] = "no product found"
	}

	c.JSON(http.StatusOK, resp)
}
