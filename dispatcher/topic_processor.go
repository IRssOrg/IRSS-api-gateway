package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"irss-gateway/models"
	"net/http"
)

var config models.Config
var QQUrl = config.Processor.QQSummary
var PassageUrl = config.Processor.TopicProcessor
var token = config.Token

// UploadPassage /*
//
//	上传文本
//
// */
func UploadPassage(content string) (string, error) {
	body := models.PushContentReq{
		Token:   token,
		Content: content,
	}
	byte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(PassageUrl+"/passage", "application/json", bytes.NewReader(byte))
	if err != nil {
		return "", err
	}
	res := models.PushContentRes{}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	if res.ErrNo != 0 {
		return "", fmt.Errorf(res.Message)
	} else {
		return res.Hash, nil
	}
}

func AskQuestion(hash string, question string) (string, error) {
	body := models.TopicActionReq{
		Action: "ask",
		Token:  token,
		Param:  question,
	}
	byte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(PassageUrl+"/passage/"+hash, "application/json", bytes.NewReader(byte))
	if err != nil {
		return "", err
	}
	res := models.AskQuestionRes{}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	if res.ErrNo != 0 {
		return "", fmt.Errorf(res.Message)
	} else {
		return res.Content, nil
	}
}

func Summary(hash string) (string, error) {
	body := models.TopicActionReq{
		Action: "action",
		Token:  token,
	}
	byte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(PassageUrl+"/passage/"+hash, "application/json", bytes.NewReader(byte))
	if err != nil {
		return "", err
	}
	res := models.SummaryRes{}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	if res.ErrNo != 0 {
		return "", fmt.Errorf(res.Message)
	} else {
		return res.Content, nil
	}
}

func GetPassageTopics(hash string) ([]models.TopicWithRelative, error) {
	body := models.TopicActionReq{
		Action: "topic",
		Token:  token,
	}
	byte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(PassageUrl+"/passage/"+hash, "application/json", bytes.NewReader(byte))
	if err != nil {
		return nil, err
	}
	res := models.GetTopicRes{}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	if res.ErrNo != 0 {
		return nil, fmt.Errorf(res.Message)
	} else {
		return res.Topics, nil
	}
}

func ConfirmTopicWithRelative(hash string, topic string) ([]models.TopicWithRelative, error) {
	body := models.TopicActionReq{
		Action: "getTopicRelative",
		Param:  topic,
		Token:  token,
	}
	byte, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(PassageUrl+"/passage/"+hash, "application/json", bytes.NewReader(byte))
	if err != nil {
		return nil, err
	}
	res := models.GetTopicRelativeRes{}
	err = json.NewDecoder(resp.Body).Decode(&resp)
	if res.ErrNo != 0 {
		return nil, fmt.Errorf(res.Message)
	} else {
		return res.Topics, nil
	}
}
