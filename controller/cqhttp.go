package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type TopicGet struct {
	Topic   string `json:"Topic"`
	Group   string `json:"Group"`
	User    string `json:"User"`
	RawText string `json:"RawText"`
	Summary string `json:"Summary"`
}

type Message struct {
	Content         string `json:"content"`
	Time            string `json:"time"`
	Topic           string `json:"topic"`
	OriginalContent string `json:"original_content"`
	Id              string `json:"id"`
	NickName        string `json:"nickname"`
}

type Summary struct {
	Group   string `json:"Group"`
	RawText string `json:"RawText"`
	Summary string `json:"Summary"`
}

func TopicListener(c *gin.Context) {
	var req TopicGet
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "请求格式错误",
		})
		return
	}
	var id int64
	err := pool.QueryRow("select user_id from account where account=?", req.User).Scan(&id)
	if err != nil {
		log.Println("[TopicListener] query fail", err)
		return
	}
	original := Message{
		Content:         req.Summary,
		Time:            strconv.FormatInt(time.Now().Unix(), 10),
		OriginalContent: req.RawText,
		NickName:        req.User,
		Topic:           req.Topic,
	}

	conn, ok := wsPool[id]
	isChecked := 1
	if !ok {
		isChecked = 0
		_, err := StoreMessage(id, original, isChecked, 1)
		if err != nil {
			log.Println("[TopicListener] store message fail", err)
			return
		}
		return
	}
	lastId, err := StoreMessage(id, original, isChecked, 1)
	original.Id = strconv.Itoa(int(lastId))
	var list []Message
	list = append(list, original)
	msg := gin.H{
		"event":    "message notification",
		"messages": list,
	}
	if req.Topic == config.VIPUserNickName {
		msg["event"] = "mentioned notification"
	}
	err = conn.WriteJSON(msg)
}

func StoreMessage(id int64, msg any, checked int, typeCode int) (int64, error) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return -1, err
	}
	result, err := pool.Exec("insert into message (user_id, message, is_checked, type) values (?, ?, ?. ?)", id, string(msgJson), checked, typeCode)
	if err != nil {
		return -1, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return lastId, nil
}

func pushMessageNow(id int64, typeCode int) error {
	IsNew := true
	var msgRef []byte
	rows, err := pool.Query("select message from message where user_id=? and type=1", id)
	if err := pool.QueryRow("select message from message where user_id=? and type=1 and checked=0", id).Scan(&msgRef); err != nil {
		IsNew = false
	}
	if typeCode == 2 {
		rows, err = pool.Query("select message from message where user_id=? and type=2", id)
		if err := pool.QueryRow("select message from message where user_id=? and type=2 and checked=0", id).Scan(&msgRef); err != nil {
			IsNew = false
		}
	}
	if err != nil {
		log.Println("[pushMessageNow] query fail", err)
		return err
	}
	var list []Message
	var sumList []Summary
	for rows.Next() {
		var msg Message
		var sum Summary
		var msgJson string
		err := rows.Scan(&msgJson)
		if err != nil {
			log.Println("[pushMessageNow] scan fail", err)
			return err
		}
		if typeCode == 1 {
			err = json.Unmarshal([]byte(msgJson), &msg)
			if err != nil {
				log.Println("[pushMessageNow] unmarshal fail", err)
				return err
			}
			list = append(list, msg)
			continue
		}
		err = json.Unmarshal([]byte(msgJson), &sum)
		if err != nil {
			log.Println("[pushMessageNow] unmarshal fail", err)
			return err
		}
		sumList = append(sumList, sum)

	}
	conn, ok := wsPool[id]
	if !ok {
		return nil
	}
	_, err = pool.Exec("update message set is_checked=1 where user_id=? and is_checked=0 and type=1", id)
	if typeCode == 2 {
		_, err = pool.Exec("update message set is_checked=1 where user_id=? and is_checked=0 and type=2", id)
	}
	if err != nil {
		log.Println("[pushMessageNow] update fail", err)
		return err
	}
	if typeCode == 1 {
		if err := conn.WriteJSON(gin.H{
			"event":    "message notification",
			"messages": list,
		}); err != nil {
			log.Println("[pushMessageNow] write json fail", err)
			return err
		}
		return nil
	}
	if err := conn.WriteJSON(gin.H{
		"event":    "message summary",
		"is_new":   IsNew,
		"messages": sumList,
	}); err != nil {
		log.Println("[pushMessageNow] write json fail", err)
		return err
	}
	return nil
}

func MessageListener(c *gin.Context) {
	id := config.VIPUser
	var req Summary
	_ = c.ShouldBindJSON(&req)
	conn, ok := wsPool[id]
	if ok {
		_, err := StoreMessage(int64(id), req, 1, 2)
		if err != nil {
			log.Println("[MessageListener] store message fail", err)
			return
		}
		if err := conn.WriteJSON(gin.H{
			"event":   "message summary",
			"summary": req.Summary,
		}); err != nil {
			log.Println("[MessageListener] write json fail", err)
		}
		return
	}
	_, err := StoreMessage(int64(id), req, 0, 2)
	if err != nil {
		log.Println("[MessageListener] store message fail", err)
		return
	}
}
