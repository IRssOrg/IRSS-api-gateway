package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"irss-gateway/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Auth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	log.Println("tokenString", tokenString)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	id := int(claims["userId"].(float64))
	log.Println(id)
	c.Set("userId", id)
	c.Set("username", claims["username"])
	c.Next()
}

func Login(c *gin.Context) {
	var req models.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, models.LoginResp{
			StatusCode: 400,
			StatusMsg:  "Bad Request",
		})
		log.Println(err)
		return
	}
	var passwordref string
	var id int64
	err := pool.QueryRow("select password, id from public.users where username = ?", req.Username).Scan(&passwordref, &id)
	if err != nil {
		log.Println("username invalid in auth ", err)
		return
	}
	if req.Password != passwordref {
		c.JSON(200, models.LoginResp{
			StatusCode: 1,
			StatusMsg:  "password or username is wrong",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"userId":   id,
		"exp":      time.Now().Add(time.Hour * 720).Unix(), //todo 为了测试方便，暂时设置为一个token有效期为30天
	})
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, models.LoginResp{
		StatusCode: 2,
		StatusMsg:  "auth successfully",
		Token:      signedToken,
		Id:         id,
	})
	go func() {
		_ = SubscriptionTimer(id)
	}()
	return
}

func Register(c *gin.Context) {
	var req models.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, models.RegisterResp{
			StatusCode: 400,
			StatusMsg:  "Bad Request",
		})
		log.Println(err)
		return
	}
	var id int64
	err := pool.QueryRow("select id from public.users where username = ?", req.Username).Scan(&id)
	if err == nil {
		c.JSON(200, models.RegisterResp{
			StatusCode: 0,
			StatusMsg:  "username already exists",
		})
		return
	}

	stmt, err := pool.Prepare("insert into public.users(username, password, article_topic, qq_topic, qq, selected_topic, wechat_sub, zhihu_sub, bilibili_sub) values (?, ?, ?, ?, ?, ?)")
	result, err := stmt.Exec(req.Username, req.Password, "[]", "[]", "{\"accounts\":[]}", "[]", "[]", "[]", "[]")
	id, err = result.LastInsertId()
	if err != nil {
		c.JSON(200, models.RegisterResp{
			StatusCode: 1,
			StatusMsg:  "register failed",
		})
		log.Println("[Register]fail to insert into database in register ", err)
		return
	}
	_, err = pool.Exec("CREATE TABLE " + strconv.Itoa(int(id)) + "_article" + " (id bigint NOT NULL AUTO_INCREMENT,  content varchar(6000) NULL,  time varchar(255) NULL, media_type varchar(255) NULL,  topic varchar(255) NULL,  PRIMARY KEY (id), checked int NULL, platform varchar(255) NULL);")
	if err != nil {
		log.Println("[Register]fail to create table in register ", err)
		return
	}
	go func() {
		_ = SubscriptionTimer(id)
	}()

	c.JSON(200, models.RegisterResp{
		StatusCode: 2,
		StatusMsg:  "register successfully",
		Id:         id,
	})

}
