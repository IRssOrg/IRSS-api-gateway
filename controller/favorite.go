package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"irss-gateway/models"
	"log"
)

type Favorite struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type FavoriteList struct {
	Favorites []Favorite `json:"favorites"`
}

type DeleteFavoriteReq struct {
	Type string   `json:"type"`
	Ids  []string `json:"ids"`
}

func AddFavorite(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[AddFavorite] get userId fail")
		return
	}
	id := idCode.(int)
	var req FavoriteList
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[AddFavorite] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	favoriteList, err := GetFavoriteList(int64(id))
	favoriteList = append(favoriteList, req.Favorites...)
	favoriteByte, err := json.Marshal(favoriteList)
	if err != nil {
		log.Println("[AddFavorite] marshal fail", err)
		return
	}
	_, err = pool.Exec("update public.users set favorite_article=? where id=?", string(favoriteByte), id)
	c.JSON(200, gin.H{
		"favorites":   favoriteList,
		"status_code": 0,
	})
	return
}

func GetFavoriteList(id int64) ([]Favorite, error) {
	var favoriteByte []byte
	var favoriteList []Favorite
	err := pool.QueryRow("select favorite_article from public.users where id=?", id).Scan(&favoriteByte)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(favoriteByte, &favoriteList)
	if err != nil {
		return nil, err
	}
	return favoriteList, nil
}

func GetFavorite(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[GetFavorite] get userId fail")
		return
	}
	id := idCode.(int)
	favoriteList, err := GetFavoriteList(int64(id))
	if err != nil {
		log.Println("[GetFavorite] get favorite list fail", err)
		return
	}
	c.JSON(200, gin.H{
		"favorites": favoriteList,
	})
}

func DeleteFavorite(c *gin.Context) {
	idCode, ok := c.Get("userId")
	if !ok {
		log.Println("[DeleteFavorite] get userId fail")
		return
	}
	id := idCode.(int)
	var req DeleteFavoriteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[DeleteFavorite] read json fail", err)
		c.JSON(400, models.DefaultResp{
			StatusCode: 1,
			StatusMsg:  "请求格式错误",
		})
		return
	}
	favoriteList, err := GetFavoriteList(int64(id))
	if err != nil {
		log.Println("[DeleteFavorite] get favorite list fail", err)
		return
	}
	for i, v := range favoriteList {
		for _, id := range req.Ids {
			if v.Id == id {
				favoriteList = append(favoriteList[:i], favoriteList[i+1:]...)
			}
		}
	}
	favoriteByte, err := json.Marshal(favoriteList)
	if err != nil {
		log.Println("[DeleteFavorite] marshal fail", err)
		return
	}
	_, err = pool.Exec("update public.users set favorite_article=? where id=?", string(favoriteByte), id)
	c.JSON(200, gin.H{
		"favorites":   favoriteList,
		"status_code": 0,
	})
}
