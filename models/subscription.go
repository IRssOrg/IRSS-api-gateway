package models

type AddTopicsReq struct {
	Topic string `json:"topic"`
}

type Topic struct {
	Topic string `json:"topic"`
	Id    string `json:"id"`
}

type SelectedTopics struct {
	Ids []string `json:"ids"`
}

type TopicsResp struct {
	Topics     []Topic `json:"topics"`
	StatusCode int     `json:"statusCode"`
}
type DefaultResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
}

type GetTopicsResp struct {
	StatusCode int      `json:"statusCode"`
	StatusMsg  string   `json:"statusMsg"`
	Topics     []string `json:"topics"`
}

type UserConfig struct {
	ArticleTime int `json:"articleTime"`
	MessageTime int `json:"messageTime"`
	VideoTime   int `json:"videoTime"`
}
