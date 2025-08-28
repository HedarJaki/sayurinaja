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
	router.PUT("/my-store", middleware.AuthStoreMiddleware(), controller.SignStore)
	router.Run() // listen and serve on 0.0.0.0:8080
}
