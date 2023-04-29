package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
)

func AddAccount(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[SetAccount] get userId fail")
		return
	}
	log.Println(idCode)
	id := idCode.(int)
	var req models.AddAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[SetAccount] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}

	var accountByte []byte
	err := pool.QueryRow("select qq from public.users where id=?", id).Scan(&accountByte)
	if err != nil {
		log.Println("[SetAccount] query fail", err)
		return
	}
	var accounts models.Accounts
	err = json.Unmarshal(accountByte, &accounts)
	log.Println(string(accountByte))
	if err != nil {
		log.Println("[SetAccount] unmarshal fail", err)
		return
	}
	accounts.Accounts = append(accounts.Accounts, req)
	accountByte, err = json.Marshal(accounts)
	_, err = pool.Exec("update public.users set qq=? where id=?", string(accountByte), id)
	if err != nil {
		log.Println("[SetAccount] update fail", err)
		return
	}
	c.JSON(200, models.DefaultResp{
		StatusCode: 0,
		StatusMsg:  "添加成功",
	})
	return
}

func GetAccount(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetAccount] get userId fail")
		return
	}
	id := idCode.(int)
	var accountByte []byte
	err := pool.QueryRow("select qq from public.users where id=?", id).Scan(&accountByte)
	log.Println(string(accountByte))
	if err != nil {
		log.Println("[GetAccount] query fail", err)
		return
	}
	var accounts models.Accounts
	err = json.Unmarshal(accountByte, &accounts)
	if err != nil {
		log.Println("[GetAccount] unmarshal fail", err)
		return
	}
	c.JSON(200, accounts)
	return
}
