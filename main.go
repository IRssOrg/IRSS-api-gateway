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

	apiRouter := r.Group("/info", controller.Auth)
	apiRouter.PATCH("/:type/topics/selected", controller.AddTopics)
	apiRouter.GET("/:type/topics/selected", controller.GetSelectedTopics)
	apiRouter.POST("/:type/topics/selected", controller.SetTopics)
	//apiRouter.DELETE("/:type/topics/selected", controller.DeleteTopics)
	apiRouter.GET("/:type/topics", controller.GetTopics)
	apiRouter.POST("/config", controller.SetConfig)
	apiRouter.GET("/ws", controller.WsHandler, controller.SubscriptionTimer)

	apiRouter.POST("/subscription/author/:platform", controller.SearchAuthor)
	apiRouter.POST("/subscription/author/:platform/subscribed", controller.AddSubscription)
	apiRouter.DELETE("/subscription/author/:platform/subscribed", controller.DeleteSubscription)

	r.GET("/account", controller.Auth, controller.GetAccount)
	r.POST("/account", controller.Auth, controller.AddAccount)

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
