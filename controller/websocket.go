package controller

import (
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"sync"
	"time"
)
import "github.com/gorilla/websocket"

var (
	wsPool   = make(map[int]*websocket.Conn)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mutex = &sync.Mutex{}
)

func WsHandler(c *gin.Context) {
	idCode, ok := c.Get("id")
	if !ok {
		return
	}
	id := int(idCode.(float64))
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
		})
		return
	}
	mutex.Lock()
	wsPool[id] = conn
	mutex.Unlock()
	c.Next()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws read fail", err)
			break
		}
	}
	mutex.Lock()
	delete(wsPool, id)
	mutex.Unlock()
}

func SubscriptionTimer(c *gin.Context) {
	log.Println("[SubscriptionTimer] running")
	idCode, ok := c.Get("id")
	if !ok {
		return
	}
	id := idCode.(int)
	config, err := GetUserConfig(int64(id))
	if err != nil {
		log.Println("[SubscriptionTimer] get user config fail", err)
		return
	}
	go func(id int) { // article timer
		for {
			time.Sleep(time.Second * time.Duration(config.ArticleTime))
		}
	}(id)
	go func(id int) { // message timer
		for {
			time.Sleep(time.Second * time.Duration(config.MessageTime))
		}
	}(id)
	go func(id int) { // video timer
		for {
			time.Sleep(time.Second * time.Duration(config.VideoTime))
		}
	}(id)
}

func GetUserConfig(id int64) (models.UserConfig, error) {
	log.Println("[GetUserConfig] running id:", id)
	var config models.UserConfig
	err := pool.QueryRow("select article_time,message_time,video_time from public.users where id=?", id).Scan(&config.ArticleTime, &config.MessageTime, &config.VideoTime)
	if err != nil {
		log.Println("[GetUserConfig] query fail", err)
		return config, err
	}
	return config, nil
}
