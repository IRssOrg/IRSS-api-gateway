package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Note struct {
	Content string `json:"content"`
	Time    string `json:"time"`
}

type NoteResp struct {
	ArticleId string `json:"article_id"`
	Note      Note   `json:"note"`
}

type DeleteNoteReq struct {
	ArticleIds []string `json:"article_ids"`
}

func GetNote(c *gin.Context) {
	idCode, ok := c.Get("id")
	if !ok {
		return
	}
	id := idCode.(int64)
	noteList, err := GetNoteList(id)
	if err != nil {
		log.Println("[GetNote] get note list fail", err)
		c.JSON(500, gin.H{
			"status": "fail",
		})
		return
	}
	c.JSON(200, gin.H{
		"notes": noteList,
	})
}

func GetNoteList(id int64) ([]NoteResp, error) {
	rows, err := pool.Query("select content, time, article_id from public.note where user_id = ?", id)
	if err != nil {
		log.Println("[GetNoteList] query fail", err)
		return nil, err
	}
	var list []NoteResp
	for rows.Next() {
		var note Note
		var articleId string
		err := rows.Scan(&note.Content, &note.Time, &articleId)
		if err != nil {
			log.Println("[GetNoteList] scan fail", err)
			continue
		}
		list = append(list, NoteResp{
			ArticleId: articleId,
			Note:      note,
		})

	}
	return list, nil
}

func SetNote(c *gin.Context) {
	idCode, ok := c.Get("id")
	if !ok {
		return
	}
	id := idCode.(int64)
	var req NoteResp
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"msg":    "invalid request",
		})
		return
	}
	_, err := pool.Exec("insert into public.note (article_id, user_id, content, time) values (?, ?, ?, ?)", req.ArticleId, id, req.Note.Content, req.Note.Time)
	if err != nil {
		log.Println("[SetNote] insert note fail", err)
		return
	}
	noteList, err := GetNoteList(id)
	if err != nil {
		log.Println("[SetNote] get note list fail", err)
		return
	}
	c.JSON(200, gin.H{
		"status_code": "success",
		"notes":       noteList,
	})
}

func DeleteNote(c *gin.Context) {
	idCode, ok := c.Get("id")
	if !ok {
		return
	}
	id := idCode.(int64)
	var req DeleteNoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"msg":    "invalid request",
		})
		return
	}
	for _, articleId := range req.ArticleIds {
		_, err := pool.Exec("delete from public.note where article_id = ? and user_id = ?", articleId, id)
		if err != nil {
			log.Println("[DeleteNote] delete note fail", err)
			continue
		}
	}
	noteList, err := GetNoteList(id)
	if err != nil {
		log.Println("[DeleteNote] get note list fail", err)
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"notes":       noteList,
	})

}
