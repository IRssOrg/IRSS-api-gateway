package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/dispatcher"
	"irss-gateway/models"
	"log"
	"net/http"
	"strconv"
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
var UserTopics = make(map[int64][]string)
var Timers = make(map[int64]*cron.Cron)
var LastUpdateTimeMap = make(map[int64]timeList)

//var waitGroup sync.WaitGroup

func SubscriptionTimer(id int64) error {
	config, err := GetUserConfig(id)
	subList, err := GetUserSubscription(id)
	topics, err := GetTopicsList(id)
	if err != nil {
		log.Println("[SubscriptionTimer] get user config fail", err)
		return err
	}
	UserConfigMap[id] = config
	UserSubListMap[id] = subList
	UserTopics[id] = topics
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
	//if err := pushArticleNow(id); err != nil {
	//	log.Println("[SubscriptionTimer] push article now fail", err)
	//}
	_, err = cronObj.AddFunc(config.ArticleTime, func() {
		timeRef, ok := LastUpdateTimeMap[id]
		if !ok {
			timeRef.Zhihu = time.Now().Unix()
			timeRef.Wechat = time.Now().Unix()
			timeRef.Bilibili = time.Now().Unix()
			//timeRef.Zhihu = 100000000
			//timeRef.Wechat = 100000000
			//timeRef.Bilibili = 10000000
		}
		conn, ok := wsPool[int(id)]
		var pushEvent []ArticleResp
		for _, author := range UserSubListMap[id].Zhihu {
			resp, err := GetFromAuthor(author.Id, timeRef.Zhihu, "zhihu")
			if err != nil {
				log.Println("[SubscriptionTimer] get zhihu author fail", err)
				continue
			}
			for _, v := range resp {
				article, err := SearchPassage(strconv.Itoa(int(v.Id)), "zhihu")
				if err != nil {
					log.Println("[SubscriptionTimer] get zhihu article fail", err)
					continue
				}
				article.Time = v.Time
				hash, err := dispatcher.UploadPassage(article.Content)
				topicString, err := GetTopicString(id)
				relatives, err := dispatcher.ConfirmTopicWithRelative(hash, topicString)
				for _, v := range relatives {
					rel, err := strconv.ParseFloat(v.Relative, 64)
					if err != nil {
						log.Println("[SubscriptionTimer] parse relative fail", err)
						continue
					}
					if rel > 0.5 {
						LastId, _ := StoreArticle(id, article, ok)
						article.Id = string(LastId)
						article.Topic = v.Topic
						pushEvent = append(pushEvent, article)
					}
				}

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
	var articleList []ArticleResp
	rows, err := pool.Query("select id, content, time, media_type, topic, platform from " + strconv.Itoa(int(id)) + "_article where checked=0")
	if err != nil {
		log.Println("[pushArticleNow] query fail", err)
		return err
	}

	for rows.Next() {
		var article ArticleResp
		err := rows.Scan(&article.Id, &article.Content, &article.Time, &article.MediaType, &article.Topic, &article.Platform)
		if err != nil {
			log.Println("[pushArticleNow] scan fail", err)
			return err
		}
		articleList = append(articleList, article)
	}
	_, err = pool.Exec("update " + strconv.Itoa(int(id)) + "_article set checked=1 where checked=0")
	conn, ok := wsPool[int(id)]
	if !ok {
		log.Println("[pushArticleNow] conn not exist")
		return nil
	}
	if err := conn.WriteJSON(gin.H{
		"event":    "article notification",
		"articles": articleList,
	}); err != nil {
		log.Println("[pushArticleNow] push article fail", err)
		return err
	}
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
	//log.Println(url)
	resp, err := http.Get(url)
	var respList passageResp
	if err != nil {
		return respList.Ret, err
	}
	err = json.NewDecoder(resp.Body).Decode(&respList)
	//log.Println(respList.Ret)
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
	result, err := pool.Exec("INSERT INTO "+strconv.Itoa(int(id))+"_article (platform, media_type, content, topic, checked, time) VALUES (?, ?, ?, ?, ?, ?)", article.Platform, article.MediaType, article.Content, article.Topic, isChecked, article.Time)
	if err != nil {
		log.Println("[StoreArticle] store article fail", err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
