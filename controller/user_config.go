package controller

import (
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
)

func SetConfig(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[SetConfig] get userId fail")
		return
	}
	id := idCode.(int)
	var req models.UserConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[SetConfig] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}

	_, err := pool.Exec("update public.users set article_time=?,message_time=?,video_time=? where id=?",
		req.ArticleTime, req.MessageTime, req.VideoTime, id)
	if err != nil {
		log.Println("[SetConfig] exec fail", err)
		return
	}
	c.JSON(200, models.DefaultResp{
		StatusCode: 0,
		StatusMsg:  "设置成功",
	})
	return

}
