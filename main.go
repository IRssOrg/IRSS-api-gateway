package main

import (
	"connection-gateway/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func initRouter(r *gin.Engine) {
	beforeAuth := r.Group("/auth")
	beforeAuth.POST("/login", controller.Login)
	beforeAuth.POST("/register", controller.Register)

	apiRouter := r.Group("/connection", controller.Auth)
	apiRouter.POST("/test")
}

func main() {
	r := gin.Default()
	initRouter(r)
	if err := r.Run("0.0.0.0:11451"); err != nil {
		log.Fatal(err)
	}
}
