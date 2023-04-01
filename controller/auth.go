package controller

import (
	"connection-gateway/lib"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	c.Set("userId", claims["userId"])
	c.Set("username", claims["username"])
	c.Next()
}

func Login(c *gin.Context) {
	var req lib.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, lib.LoginResp{
			StatusCode: 400,
			StatusMsg:  "Bad Request",
		})
		log.Println(err)
		return
	}
	var passwordref string
	var id int64
	err := pool.QueryRow("select password, id from public.user where username = $1", req.Username).Scan(&passwordref, &id)
	if err != nil {
		log.Println("username invalid in auth ", err)
		return
	}
	if req.Password != passwordref {
		c.JSON(200, lib.LoginResp{
			StatusCode: 1,
			StatusMsg:  "password or username is wrong",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"userId":   id,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, lib.LoginResp{
		StatusCode: 2,
		StatusMsg:  "auth successfully",
		Token:      signedToken,
		Id:         id,
	})

	return
}

func Register(c *gin.Context) {
	var req lib.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, lib.RegisterResp{
			StatusCode: 400,
			StatusMsg:  "Bad Request",
		})
		log.Println(err)
		return
	}
	var id int64
	err := pool.QueryRow("select id from public.user where username = $1", req.Username).Scan(&id)
	if err == nil {
		c.JSON(200, lib.RegisterResp{
			StatusCode: 0,
			StatusMsg:  "username already exists",
		})
		return
	}
	err = pool.QueryRow("insert into public.user(username, password) values($1, $2) returning id", req.Username, req.Password).Scan(&id)
	if err != nil {
		c.JSON(200, lib.RegisterResp{
			StatusCode: 1,
			StatusMsg:  "register failed",
		})
		log.Println("fail to insert into database in register ", err)
		return
	}
	c.JSON(200, lib.RegisterResp{
		StatusCode: 2,
		StatusMsg:  "register successfully",
		Id:         id,
	})
	return
}
