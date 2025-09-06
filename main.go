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

func main() {
	router := gin.Default()
	router.POST("/signup", controller.SignUp)
	router.POST("/Login", controller.Login)
	router.GET("/home", middleware.AuthMiddleware(), controller.HomePage)
	router.PUT("/sign-store", middleware.AuthMiddleware(), middleware.AuthStoreMiddleware(), controller.SignStore) //masih perlu diupdate
	router.POST("/my_store/addproduct", middleware.AuthStoreMiddleware(), controller.AddProduct)
	router.GET("/store/:id", controller.CheckStore)
	router.GET("/search-product", controller.ShowProductByCategory)
	router.POST("add-product-to-cart/:id", middleware.AuthMiddleware(), middleware.CartCheckMiddleware(), controller.AddToCart)
	router.PUT("/cart", middleware.AuthMiddleware(), middleware.CartCheckMiddleware(), controller.ShowCart)
	router.GET("/store")
	router.Run() // listen and serve on 0.0.0.0:8080
}
