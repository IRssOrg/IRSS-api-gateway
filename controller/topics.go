package controller

import (
	"database/sql"
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
	if where == "article" {
		if err := SetArticleTopics(id, topics); err != nil {
			c.JSON(400, models.DefaultResp{
				StatusCode: 1,
				StatusMsg:  "设置失败",
			})
			return
		}

	} else {
		if err := SetQQTopics(id, topics); err != nil {
			c.JSON(400, models.DefaultResp{
				StatusCode: 1,
				StatusMsg:  "设置失败",
			})
			return
		}
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

func SetArticleTopics(id any, topics []string) error {
	_, err := pool.Exec("update public.users set article_topic='[]' where id=?", id)
	stmt, err := pool.Prepare("update public.users set article_topic=json_array_append(article_topic, '$', ?) where id=?")
	if err != nil {
		log.Println("[setArticleTopic] prepare stmt fail", err)
		return err
	}
	for _, topic := range topics {
		_, err := stmt.Exec(topic, id)
		if err != nil {
			log.Println("[setArticleTopic] exec stmt fail", err)
			return err
		}
	}
	return nil
}

func SetQQTopics(id any, topics []string) error {
	_, err := pool.Exec("update public.users set qq_topic='[]' where id=?", id)
	stmt, err := pool.Prepare("update public.users set qq_topic=json_array_append(qq_topic, '$', ?) where id=?")
	if err != nil {
		log.Println("[setQQTopic] prepare stmt fail", err)
		return err
	}
	for _, topic := range topics {
		_, err := stmt.Exec(topic, id)
		if err != nil {
			log.Println("[setQQTopic] exec stmt fail", err)
			return err
		}
	}
	return nil
}
