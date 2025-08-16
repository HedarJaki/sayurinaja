package main

import (
	"mobapp/controller"
	"mobapp/initializer"

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
	router.Run() // listen and serve on 0.0.0.0:8080
}
