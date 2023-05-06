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
	var topicIds []string
	err = json.Unmarshal(topicbyte, &topics)
	err = json.Unmarshal(topicIdByte, &topicIds)
	if err != nil {
		log.Println("[setArticleTopics] unmarshal fail", err)
		return
	}
	topics = append(topics, topic)
	topicIds = append(topicIds, topic)
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
	go func() {
		_ = SubscriptionTimer(id.(int64))
	}()
}

func AddTopicList(topicsRef []string, id int64) error {
	var topicByte []byte
	var topics []string
	err := pool.QueryRow("select article_topic from public.users where id=?", id).Scan(&topicByte)
	if err != nil {
		log.Println("[setArticleTopics] query fail", err)
		return err
	}
	err = json.Unmarshal(topicByte, &topics)
	if err != nil {
		log.Println("[setArticleTopics] unmarshal fail", err)
		return err
	}
	topics = append(topics, topicsRef...)
	jsonArray, err := json.Marshal(topics)
	if err != nil {
		log.Println("[setArticleTopic] json marshal fail", err)
		return err
	}
	_, err = pool.Exec("update public.users set article_topic=? where id=?", string(jsonArray), id)
	if err != nil {
		log.Println("[setArticleTopic] exec stmt fail", err)
		return err
	}
	conn, ok := wsPool[id]
	if !ok {
		log.Println("[setArticleTopic] get conn fail")
		return nil
	}
	if err := conn.WriteJSON(gin.H{
		"event":  "article topic update",
		"topics": topics,
	}); err != nil {
		log.Println("[setArticleTopic] write json fail", err)
	}
	return nil
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

func GetTopicsList(id int64) ([]string, error) {
	var topicByte []byte
	err := pool.QueryRow("select article_topic from public.users where id=?", id).Scan(&topicByte)
	if err != nil {
		log.Println("[GetTopics] query fail", err)
		return nil, err
	}
	var topics []string
	err = json.Unmarshal(topicByte, &topics)
	if err != nil {
		log.Println("[GetTopics] unmarshal fail", err)
		return nil, err
	}
	return topics, nil
}

func GetTopicString(id int64) (string, error) {
	topics, err := GetTopicsList(id)
	if err != nil {
		return "", err
	}
	var topicString string
	for _, v := range topics {
		topicString += v + ","
	}
	return topicString, nil
}

func GetSelectedTopics(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetTopics] get userId fail")
		return
	}
	id := idCode.(int)
	var topicByte []byte
	err := pool.QueryRow("select selected_topic from public.users where id=?", id).Scan(&topicByte)
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
	c.JSON(200, models.SelectedTopics{
		Topics: topics,
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

	topicIds, err := json.Marshal(req.Topics)
	if err != nil {
		log.Println("[setArticleTopics] marshal fail", err)
		return
	}
	_, err = pool.Exec("update public.users set selected_topic=? where id=?", string(topicIds), id)
	if err != nil {
		log.Println("[setArticleTopic] exec stmt fail", err)
		return
	}
	go func() {
		_ = SubscriptionTimer(int64(id))
	}()
	c.JSON(200, req)

}
