package controller

import (
	"irss-gateway/models"
	"log"
	"time"
)
import "github.com/robfig/cron/v3"

type timeList struct {
	Zhihu    int64
	Wechat   int64
	Bilibili int64
}

var UserConfigMap map[int64]models.UserConfig
var UserSubListMap map[int64]UserSubList
var Timers map[int64]*cron.Cron
var LastUpdateTimeMap map[int64]timeList

func SubscriptionTimer(id int64) error {
	config, err := GetUserConfig(id)
	subList, err := GetUserSubscription(id)
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
	}
	_, err = cronObj.AddFunc(config.ArticleTime, func() {
		timeRef, ok := LastUpdateTimeMap[id]
		if !ok {
			timeRef.Zhihu = time.Now().Unix()
			timeRef.Wechat = time.Now().Unix()
			timeRef.Bilibili = time.Now().Unix()
		}
		for _, author := range UserSubListMap[id].Zhihu {
			ZhihuAuthor(author.Id, timeRef.Zhihu)
		}
	})
	return nil
}

func ZhihuAuthor(id string, timeRef int64) {

}
