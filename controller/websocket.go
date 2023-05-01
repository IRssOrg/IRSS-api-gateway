package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"sync"
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

type UserSubList struct {
	Zhihu    []Author `json:"zhihu"`
	Wechat   []Author `json:"wechat"`
	Bilibili []Author `json:"bilibili"`
}

type LastUpdateTime struct {
}

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
	if err := pushArticleNow(int64(id)); err != nil {
		log.Println("[WsHandler] pushArticleNow fail", err)
	}
	if err := pushMessageNow(int64(id)); err != nil {
		log.Println("[WsHandler] pushMessageNow fail", err)
	}
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
	c.Next()
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
