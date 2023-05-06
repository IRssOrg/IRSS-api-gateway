package main

import (
	"github.com/gin-gonic/gin"
	"irss-gateway/controller"
	"log"
)

func initRouter(r *gin.Engine) {
	r.GET("/auth", controller.Auth, controller.CheckToken)
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
	apiRouter.GET("/ws", controller.WsHandler)
	apiRouter.POST("answer", controller.GetAnswer)

	apiRouter.GET("/subscription/author/:platform/subscribed", controller.GetSubscription)
	apiRouter.POST("/subscription/author/:platform", controller.SearchAuthor)
	apiRouter.POST("/subscription/author/:platform/subscribed", controller.AddSubscription)
	apiRouter.DELETE("/subscription/author/:platform/subscribed", controller.DeleteSubscription)

	userRouter := r.Group("/user", controller.Auth)
	userRouter.GET("/note", controller.GetNote)
	userRouter.POST("/note", controller.SetNote)
	userRouter.DELETE("/note", controller.DeleteNote)
	userRouter.GET("/favorite", controller.GetFavorite)
	userRouter.POST("/favorite", controller.AddFavorite)
	userRouter.DELETE("/favorite", controller.DeleteFavorite)

	r.GET("/account", controller.Auth, controller.GetAccount)
	r.POST("/account", controller.Auth, controller.AddAccount)

	r.POST("/summary/TopicListener", controller.TopicListener)
	r.POST("/summary/Summary", controller.MessageListener)
}

func main() {
	r := gin.Default()
	if err := controller.Init(); err != nil {
		log.Fatal(err)
	}
	initRouter(r)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	} // this method will block the calling goroutine unless an err happens
}
