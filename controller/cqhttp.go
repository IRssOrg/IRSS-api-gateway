package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type TopicGet struct {
	Topic   string `json:"topic"`
	Group   string `json:"group"`
	User    string `json:"user"`
	RawText string `json:"raw_text"`
	Summary string `json:"summary"`
}

type Message struct {
	Content         string `json:"content"`
	Time            string `json:"time"`
	Topic           string `json:"topic"`
	OriginalContent string `json:"original_content"`
	Id              string `json:"id"`
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
		Time:            time.Now().Format("2006-01-02 15:04:05"),
		Topic:           req.Topic,
		OriginalContent: req.RawText,
	}
	conn, ok := wsPool[int(id)]
	isChecked := 1
	if !ok {
		isChecked = 0
		_, err := StoreMessage(id, original, isChecked)
		if err != nil {
			log.Println("[TopicListener] store message fail", err)
			return
		}
		return
	}
	lastId, err := StoreMessage(id, original, isChecked)
	original.Id = string(lastId)
	var list []Message
	list = append(list, original)
	msg := gin.H{
		"event":    "message notification",
		"messages": list,
	}
	err = conn.WriteJSON(msg)
}

func StoreMessage(id int64, msg Message, checked int) (int64, error) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return -1, err
	}
	result, err := pool.Exec("insert into message (user_id, message, is_checked) values (?, ?, ?)", id, string(msgJson), checked)
	if err != nil {
		return -1, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return lastId, nil
}

func pushMessageNow(id int64) error {
	rows, err := pool.Query("select message from message where user_id=? and is_checked=0", id)
	if err != nil {
		log.Println("[pushMessageNow] query fail", err)
		return err
	}
	var list []Message
	for rows.Next() {
		var msg Message
		var msgJson string
		err := rows.Scan(&msgJson)
		if err != nil {
			log.Println("[pushMessageNow] scan fail", err)
			return err
		}
		err = json.Unmarshal([]byte(msgJson), &msg)
		if err != nil {
			log.Println("[pushMessageNow] unmarshal fail", err)
			return err
		}
		list = append(list, msg)
	}
	conn, ok := wsPool[int(id)]
	if !ok {
		return nil
	}
	_, err = pool.Exec("update message set is_checked=1 where user_id=? and is_checked=0", id)
	if err != nil {
		log.Println("[pushMessageNow] update fail", err)
		return err
	}
	if err := conn.WriteJSON(gin.H{
		"event":    "message notification",
		"messages": list,
	}); err != nil {
		log.Println("[pushMessageNow] write json fail", err)
		return err
	}
	return nil
}
