package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
	"net/http"
	"sync"
	"time"
)
import "github.com/robfig/cron/v3"

type timeList struct {
	Zhihu    int64
	Wechat   int64
	Bilibili int64
}

type passageResp struct {
	Ret []passages `json:"ret"`
}

type passages struct {
	Title     string `json:"title"`
	Id        int64  `json:"id"`
	Time      string `json:"time"`
	TimeStamp int64  `json:"timestamp"`
}

type ArticleResp struct {
	Platform  string `json:"platform"`
	MediaType string `json:"media_type"`
	Content   string `json:"content"`
	Time      string `json:"time"`
	Topic     string `json:"topic"`
	Id        string `json:"id"`
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

var UserConfigMap = make(map[int64]models.UserConfig) //map剂的使用make函数初始化
var UserSubListMap = make(map[int64]UserSubList)
var Timers = make(map[int64]*cron.Cron)
var LastUpdateTimeMap = make(map[int64]timeList)
var waitGroup sync.WaitGroup

func SubscriptionTimer(id int64) error {
	config, err := GetUserConfig(id)
	subList, err := GetUserSubscription(id)
	if err != nil {
		log.Println("[SubscriptionTimer] get user config fail", err)
		return err
	}
	UserConfigMap[id] = config
	UserSubListMap[id] = subList
	if err != nil {
		log.Println("[SubscriptionTimer] get user config fail", err)
		return err
	}
	cronObj, exist := Timers[id]
	if !exist {
		cronObj = cron.New()
		Timers[id] = cronObj
	} else {
		cronObj.Stop()
		cronObj = cron.New()
		Timers[id] = cronObj
	}
	if err := pushArticleNow(id); err != nil {
		log.Println("[SubscriptionTimer] push article now fail", err)
	}
	_, err = cronObj.AddFunc(config.ArticleTime, func() {
		timeRef, ok := LastUpdateTimeMap[id]
		if !ok {
			timeRef.Zhihu = time.Now().Unix()
			timeRef.Wechat = time.Now().Unix()
			timeRef.Bilibili = time.Now().Unix()
		}
		log.Println(timeRef)
		conn, ok := wsPool[int(id)]
		log.Println(UserSubListMap[id])
		var pushEvent []ArticleResp
		for _, author := range UserSubListMap[id].Zhihu {
			log.Println("[SubscriptionTimer] get zhihu author", author)
			resp, err := GetFromAuthor(author.Id, timeRef.Zhihu, "zhihu")
			if err != nil {
				log.Println("[SubscriptionTimer] get zhihu author fail", err)
				continue
			}
			log.Println("[SubscriptionTimer] get zhihu author success")
			for _, v := range resp {
				article, err := SearchPassage(string(v.Id), "zhihu")
				if err != nil {
					log.Println("[SubscriptionTimer] get zhihu article fail", err)
					continue
				}
				article.Time = v.Time
				LastId, err := StoreArticle(id, article, ok)
				article.Id = string(LastId)
				pushEvent = append(pushEvent, article)

			}
		}
		if ok {
			if err := conn.WriteJSON(gin.H{
				"event":    "article notification",
				"articles": pushEvent,
			}); err != nil {
				log.Println("[SubscriptionTimer] push zhihu article fail", err)
			}
		}

	})
	cronObj.Start()
	log.Println("[SubscriptionTimer] cron start")
	//time.Sleep(time.Hour * 72)
	return nil
}

func pushArticleNow(id int64) error {
	return nil
}

func GetFromAuthor(id string, timeRef int64, platform string) ([]passages, error) {
	url := config.Spider.Zhihu + "/api/passages/" + id + "/0"
	switch platform {
	case "bilibili":
		url = config.Spider.Bilibili + "/api/passages/" + id + "/0"
	case "zhihu":
		url = config.Spider.Zhihu + "/api/passages/" + id + "/0"
	case "wechat":
		url = config.Spider.Wechat + "/api/passages/" + id + "/0"
	}
	resp, err := http.Get(url)
	var respList passageResp
	if err != nil {
		return respList.Ret, err
	}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	defer resp.Body.Close()
	if err != nil {
		return respList.Ret, err
	}
	var ret []passages
	for _, v := range respList.Ret {
		if v.TimeStamp > timeRef {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func SearchPassage(id string, platform string) (ArticleResp, error) {
	url := config.Spider.Zhihu + "/api/passage/" + id
	switch platform {
	case "bilibili":
		url = config.Spider.Bilibili + "/api/passage/" + id
	case "zhihu":
		url = config.Spider.Zhihu + "/api/passage/" + id
	case "wechat":
		url = config.Spider.Wechat + "/api/passage/" + id
	}
	resp, err := http.Get(url)
	var article Article
	if err != nil {
		return ArticleResp{}, err
	}
	err = json.NewDecoder(resp.Body).Decode(&article)
	defer resp.Body.Close()
	if err != nil {
		return ArticleResp{}, err
	}
	var ret ArticleResp
	ret.Content = article.Content
	ret.Platform = platform
	ret.Topic = "miku"
	ret.MediaType = "article"
	ret.Id = id
	return ret, nil
}

func StoreArticle(id int64, article ArticleResp, checked bool) (int64, error) {
	isChecked := 0
	if checked {
		isChecked = 1
	}
	result, err := pool.Exec("INSERT INTO article (id, platform, media_type, content, topic, checked) VALUES (?, ?, ?, ?, ?, ?)", id, article.Platform, article.MediaType, article.Content, article.Topic, isChecked)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
