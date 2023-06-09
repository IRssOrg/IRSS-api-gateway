package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"sync"
	"time"
)
import "github.com/gorilla/websocket"

var (
	wsPool   = make(map[int64]*websocket.Conn)
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	mutex = &sync.Mutex{}
)

type UserSubList struct {
	Zhihu    []Author `json:"zhihu"`
	Wechat   []Author `json:"wechat"`
	Bilibili []Author `json:"bilibili"`
}

func WsHandler(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[WsHandler] get userId fail")
		return
	}
	id := idCode.(int64)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
		})
		return
	}
	wsPool[id] = conn

	if err := pushArticleNow(id); err != nil {
		log.Println("[WsHandler] pushArticleNow fail", err)
	}
	if err := pushMessageNow(id, 1); err != nil {
		log.Println("[WsHandler] pushMessageNow fail", err)
	}
	if err := pushMessageNow(id, 2); err != nil {
		log.Println("[WsHandler] pushMessageNow fail", err)
	}
	go func() {
		_ = SubscriptionTimer(id)
	}()
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
	time.Sleep(120 * time.Second)
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

func GetUserSubscription(id int64) (UserSubList, error) {
	log.Println("[GetUserSubscription] running id:", id)
	var subList UserSubList
	var zhihuByte, wechatByte, bilibiliByte []byte
	err := pool.QueryRow("select zhihu_sub,wechat_sub,bilibili_sub from public.users where id=?", id).Scan(&zhihuByte, &wechatByte, &bilibiliByte)
	if err != nil {
		log.Println("[GetUserSubscription] query fail", err)
		return subList, err
	}
	log.Println("zhihu:", string(zhihuByte))
	err = json.Unmarshal(zhihuByte, &subList.Zhihu)
	err = json.Unmarshal(wechatByte, &subList.Wechat)
	err = json.Unmarshal(bilibiliByte, &subList.Bilibili)
	log.Println("[GetUserSubscription] subList:", subList)
	if err != nil {
		log.Println("[GetUserSubscription] query fail", err)
		return subList, err
	}
	return subList, nil
}
