package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/dispatcher"
	"log"
	"net/http"
)

type Question struct {
	Question string `json:"question"`
}

type Answer struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type AnswerResp struct {
	Ret []Answer `json:"ret"`
}

type HandledAnswer struct {
	Summary    string          `json:"summary"`
	References []AnswerContent `json:"references"`
}

type AnswerContent struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}

type AnswerItem struct {
	Ret []Content `json:"ret"`
}

type Content struct {
	Content string `json:"content"`
}

func GetAnswer(c *gin.Context) {
	var req Question
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[GetAnswer] read json fail", err)
		c.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	url := config.Spider.Zhihu + "/api/question"
	reqByte, _ := json.Marshal(gin.H{"keyword": req.Question})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqByte))
	if err != nil {
		log.Println("[GetAnswer] get fail", err)
		c.JSON(500, gin.H{
			"message": "服务器错误",
		})
		return
	}
	var answerList AnswerResp
	_ = json.NewDecoder(resp.Body).Decode(&answerList)
	defer resp.Body.Close()
	var answers []HandledAnswer
	for _, v := range answerList.Ret[:3] {
		var answerItems AnswerItem
		resp, err := http.Get(config.Spider.Zhihu + "/api/answer/" + v.Id)
		if err != nil {
			log.Println("[GetAnswer] get fail", err)
			continue
		}
		var preSummary string
		_ = json.NewDecoder(resp.Body).Decode(&answerItems)
		var references []AnswerContent
		for _, vv := range answerItems.Ret {
			references = append(references, AnswerContent{
				Content: vv.Content,
				Title:   v.Title,
			})
			preSummary += vv.Content
		}
		hash, err := dispatcher.UploadPassage(preSummary)
		if err != nil {
			log.Println("[GetAnswer] upload fail", err)
			continue
		}
		summary, err := dispatcher.Summary(hash)
		answer := HandledAnswer{
			Summary:    summary,
			References: references,
		}
		answers = append(answers, answer)
	}
	c.JSON(200, gin.H{
		"answers": answers,
	})
}
