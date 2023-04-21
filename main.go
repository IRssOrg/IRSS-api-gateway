package main

import (
	"github.com/gin-gonic/gin"
	"irss-gateway/controller"
	"log"
)

func initRouter(r *gin.Engine) {
	beforeAuth := r.Group("/auth")
	beforeAuth.POST("/login", controller.Login)
	beforeAuth.POST("/register", controller.Register)

	apiRouter := r.Group("/subscription", controller.Auth)
	apiRouter.POST("/:type/topics", controller.SetTopics)
	apiRouter.GET("/:type/topics", controller.GetTopics)
	apiRouter.POST("/account", controller.SetAccount)
	apiRouter.POST("/config", controller.SetConfig)
	apiRouter.GET("/ws", controller.WsHandler, controller.SubscriptionTimer)

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
