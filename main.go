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
	router.PUT("/sign-store", middleware.AuthStoreMiddleware(), controller.SignStore)
	router.POST("/my_store/addproduct", controller.AddProduct)
	router.GET("/store/:id", controller.CheckStore)
	router.GET("/search-product", controller.ShowProductByCategory)
	router.Run() // listen and serve on 0.0.0.0:8080
}
