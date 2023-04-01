package main

import (
	"connection-gateway/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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
	if err := controller.Init(); err != nil {
		log.Fatal(err)
	}
	initRouter(r)
	if err := r.Run("0.0.0.0:11451"); err != nil {
		log.Fatal(err)
	}
}
