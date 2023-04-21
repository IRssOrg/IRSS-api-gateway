package controller

import (
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
)

func SetAccount(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		log.Println("[SetAccount] get userId fail")
		return
	}
	var req models.SetAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[SetAccount] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	stmt, err := pool.Prepare("update public.users set zhihu=?,bilibili=?,wechat=?,qq=? where id=?")
	if err != nil {
		log.Println("[SetAccount] prepare stmt fail", err)
		return
	}
	_, err = stmt.Exec(req.Zhihu, req.Bilibili, req.Wechat, req.Qq, id)
	if err != nil {
		log.Println("[SetAccount] exec fail", err)
		return
	}
	c.JSON(200, models.DefaultResp{
		StatusCode: 0,
		StatusMsg:  "设置成功",
	})
	return
}
