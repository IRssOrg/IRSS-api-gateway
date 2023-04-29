package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"strconv"
)

func AddTopics(c *gin.Context) {
	var req models.AddTopicsReq
	where := c.Param("type")
	log.Println(where)
	id, ok := c.Get("userId")
	if !ok {
		log.Println("[setArticleTopics] get userId fail")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[setArticleTopics] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	topic := req.Topic
	log.Println("topic:", topic)
	var topicbyte []byte
	var topicIdByte []byte
	err := pool.QueryRow("select article_topic from public.users where id=?", id).Scan(&topicbyte)
	err = pool.QueryRow("select selected_topic from public.users where id=?", id).Scan(&topicIdByte)
	if err != nil {
		log.Println("[setArticleTopics] query fail", err)
		return
	}
	var topics []string
	var topicIds []int
	err = json.Unmarshal(topicbyte, &topics)
	err = json.Unmarshal(topicIdByte, &topicIds)
	if err != nil {
		log.Println("[setArticleTopics] unmarshal fail", err)
		return
	}
	topics = append(topics, topic)
	topicId := len(topics)
	log.Println("topicId:", topicId)
	topicIds = append(topicIds, topicId)
	jsonArray, err := json.Marshal(topics)
	jsonArrayId, err := json.Marshal(topicIds)
	if err != nil {
		log.Println("[setArticleTopic] json marshal fail", err)
		return
	}
	_, err = pool.Exec("update public.users set article_topic=? where id=?", string(jsonArray), id)
	_, err = pool.Exec("update public.users set selected_topic=? where id=?", string(jsonArrayId), id)
	if err != nil {
		log.Println("[setArticleTopic] exec stmt fail", err)
		return
	}
	var resp models.TopicsResp
	for i, v := range topics {
		t := models.Topic{
			Id:    strconv.Itoa(i),
			Topic: v,
		}
		resp.Topics = append(resp.Topics, t)
	}
	resp.StatusCode = 0
	c.JSON(200, resp)
	return
}

func GetTopics(c *gin.Context) {
	where := c.Param("type")
	log.Println(where)
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetTopics] get userId fail")
		return
	}
	id := idCode.(int)
	var topicByte []byte
	err := pool.QueryRow("select article_topic from public.users where id=?", id).Scan(&topicByte)
	if err != nil {
		log.Println("[GetTopics] query fail", err)
		return
	}
	var topics []string
	err = json.Unmarshal(topicByte, &topics)
	if err != nil {
		log.Println("[GetTopics] unmarshal fail", err)
		return
	}
	var resp models.TopicsResp
	for i, v := range topics {
		t := models.Topic{
			Id:    strconv.Itoa(i),
			Topic: v,
		}
		resp.Topics = append(resp.Topics, t)
	}
	resp.StatusCode = 0
	c.JSON(200, resp)
	return

}

func GetSelectedTopics(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetTopics] get userId fail")
		return
	}
	id := idCode.(int)
	var topicIdByte []byte
	err := pool.QueryRow("select selected_topic from public.users where id=?", id).Scan(&topicIdByte)
	if err != nil {
		log.Println("[GetTopics] query fail", err)
		return
	}
	var topicIds []int
	err = json.Unmarshal(topicIdByte, &topicIds)
	if err != nil {
		log.Println("[GetTopics] unmarshal fail", err)
		return
	}
	var IdString []string
	for _, v := range topicIds {
		IdString = append(IdString, strconv.Itoa(v))
	}
	c.JSON(200, models.SelectedTopics{
		Ids: IdString,
	})
}

func SetTopics(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetTopics] get userId fail")
		return
	}
	id := idCode.(int)
	var req models.SelectedTopics
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[setArticleTopics] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	var topicInts []int
	for _, v := range req.Ids {
		topicId, _ := strconv.Atoi(v)
		topicInts = append(topicInts, topicId)
	}
	topicIds, err := json.Marshal(topicInts)
	if err != nil {
		log.Println("[setArticleTopics] marshal fail", err)
		return
	}
	_, err = pool.Exec("update public.users set selected_topic=? where id=?", string(topicIds), id)
	if err != nil {
		log.Println("[setArticleTopic] exec stmt fail", err)
		return
	}
	c.JSON(200, req)

}
