package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"net/http"
)

type Author struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

type SearchAuthorResp struct {
	Ret []Author `json:"ret"`
}

type AddAuthorReq struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}
type DeleteAuthorReq struct {
	Id string `json:"id"`
}

func SearchAuthor(c *gin.Context) {
	var req models.SearchAuthor
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[SearchAuthor] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	platform := c.Param("platform")
	url := config.Spider.Zhihu
	switch platform {
	case "bilibili":
		url = config.Spider.Bilibili
	case "zhihu":
		url = config.Spider.Zhihu
	case "wechat":
		url = config.Spider.Wechat
	}
	resp, err := http.Get(url + "/api/search/author/" + req.Question)
	if err != nil {
		log.Println("[SearchAuthor] get fail", err)
		c.JSON(500, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "服务器错误",
		})
		return
	}
	reader := resp.Body
	defer reader.Close()
	var body SearchAuthorResp
	err = json.NewDecoder(reader).Decode(&body)
	if err != nil {
		log.Println("[SearchAuthor] decode fail", err)
		return
	}
	c.JSON(200, gin.H{
		"authors": body.Ret,
	})
	return
}

func AddSubscription(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[AddSubscription] get userId fail")
		return
	}
	id := idCode.(int)
	var req AddAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[AddSubscription] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	platform := c.Param("platform")
	col := "zhihu_sub"
	switch platform {
	case "bilibili":
		col = "bilibili_sub"
	case "zhihu":
		col = "zhihu_sub"
	case "wechat":
		col = "wechat_sub"
	}
	var subList []byte
	err := pool.QueryRow("select "+col+" from public.users where id=?", id).Scan(&subList)
	if err != nil {
		log.Println("[AddSubscription] query fail", err)
		return
	}
	var subscriptions []AddAuthorReq
	err = json.Unmarshal(subList, &subscriptions)
	if err != nil {
		log.Println("[AddSubscription] unmarshal fail", err)
		return
	}
	subscriptions = append(subscriptions, req)
	subList, err = json.Marshal(subscriptions)
	_, err = pool.Exec("update public.users set "+col+"=? where id=?", string(subList), id)
	if err != nil {
		log.Println("[AddSubscription] update fail", err)
		return
	}
	c.JSON(200, gin.H{
		"subscription": subscriptions,
	})
	return
}

func DeleteSubscription(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[DeleteSubscription] get userId fail")
		return
	}
	id := idCode.(int)
	var req DeleteAuthorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[DeleteSubscription] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	platform := c.Param("platform")
	col := "zhihu_sub"
	switch platform {
	case "bilibili":
		col = "bilibili_sub"
	case "zhihu":
		col = "zhihu_sub"
	case "wechat":
		col = "wechat_sub"
	}
	var subList []byte
	err := pool.QueryRow("select "+col+" from public.users where id=?", id).Scan(&subList)
	if err != nil {
		log.Println("[DeleteSubscription] query fail", err)
		return
	}
	var subscriptions []AddAuthorReq
	err = json.Unmarshal(subList, &subscriptions)
	if err != nil {
		log.Println("[DeleteSubscription] unmarshal fail", err)
		return
	}
	for i, v := range subscriptions {
		if v.Id == req.Id {
			subscriptions = append(subscriptions[:i], subscriptions[i+1:]...) //...将切片打散
			break
		}
	}
	subList, err = json.Marshal(subscriptions)
	_, err = pool.Exec("update public.users set "+col+"=? where id=?", string(subList), id)
	if err != nil {
		log.Println("[DeleteSubscription] update fail", err)
		return
	}
	c.JSON(200, gin.H{
		"subscription": subscriptions,
	})
}
