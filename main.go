package main

import (
	"mobapp/controller"
	"mobapp/initializer"
	"mobapp/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVar()
	initializer.ConnecttoDB()
	initializer.SyncDatabase()
}

//root:MHAIDARZAKI1@tcp(127.0.0.1:3306)/sayurinaja?charset=utf8mb4&parseTime=True&loc=Local

func main() {
	router := gin.Default()
	router.POST("/signup", controller.SignUp)
	router.POST("/verify-otp", controller.VerifyOTP)
	router.POST("/forgot-pw-insert-email", controller.Forgotpw)
	router.POST("/verify-otp-forgot-pw", controller.VerifyForgotPWOTP)
	router.PUT("/change-pw", controller.NewPassword)
	router.POST("/Login", controller.Login)
	router.GET("/home", controller.Home)
	router.PUT("/sign-store", middleware.AuthMiddleware(), middleware.AuthStoreMiddleware(), controller.SignStore) //masih perlu diupdate
	router.PUT("/review-store/:id", middleware.AuthMiddleware(), controller.StoreReview)
	router.POST("/addproduct", middleware.AuthMiddleware(), middleware.AuthStoreMiddleware(), controller.AddProduct)
	router.POST("review-product/:id", middleware.AuthMiddleware(), controller.ProductReview)
	router.GET("/store/:id", controller.CheckStore)
	router.GET("/search-product", controller.ShowProductByCategory)
	router.POST("add-product-to-cart/:id", middleware.AuthMiddleware(), middleware.CartCheckMiddleware(), controller.AddToCart)
	router.PUT("/cart", middleware.AuthMiddleware(), middleware.CartCheckMiddleware(), controller.ShowCart)
	router.POST("/add-order", middleware.AuthMiddleware(), controller.AddOrder)
	router.Run() // listen and serve on 0.0.0.0:8080
}
