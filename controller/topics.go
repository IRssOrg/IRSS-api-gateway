package controller

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
)

func SetTopics(c *gin.Context) {
	var req models.SetTopicsReq
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
	topics := req.Topics
	err := RecordTopics(id, topics, where)
	if err != nil {
		log.Println("[setArticleTopics] prepare stmt fail", err)
		return
	}

	c.JSON(200, models.DefaultResp{
		StatusCode: 0,
		StatusMsg:  "设置成功",
	})
	return
}

func GetTopics(c *gin.Context) {
	where := c.Param("type")
	id, ok := c.Get("userId")
	if !ok {
		log.Println("[GetTopics] get userId fail")
		return
	}
	var topicsData []byte
	var stmt *sql.Stmt
	var err error
	if where == "article" {
		stmt, err = pool.Prepare("select article_topic from public.users where id=?")
	} else {
		stmt, err = pool.Prepare("select qq_topic from public.users where id=?")
	}
	if err != nil {
		log.Println("[GetTopics] prepare stmt fail", err)
		return
	}
	err = stmt.QueryRow(id).Scan(&topicsData)
	if topics, err := models.JsonArray2Slice(topicsData); err != nil {
		log.Println("[GetTopics] json to slice fail", err)
		return
	} else {
		c.JSON(200, models.GetTopicsResp{
			StatusCode: 0,
			StatusMsg:  "获取成功",
			Topics:     topics,
		})
		return
	}
}

func RecordTopics(id any, topics []string, where string) error {
	jsonArray, err := json.Marshal(topics)
	if err != nil {
		log.Println("[setArticleTopic] json marshal fail", err)
		return err
	}
	print(string(jsonArray))
	if where == "article" {
		_, err = pool.Exec("update public.users set article_topic=? where id=?", string(jsonArray), id)
	} else {
		_, err = pool.Exec("update public.users set qq_topic=? where id=?", string(jsonArray), id)
	}
	if err != nil {
		log.Println("[setArticleTopic] exec stmt fail", err)
		return err
	}
	return nil
}
