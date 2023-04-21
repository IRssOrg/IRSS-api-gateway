package models

type SetTopicsReq struct {
	Topics []string `json:"topics"`
}

type DefaultResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
}

type GetTopicsResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
	Topics     []string
}

type SetAccountReq struct {
	Zhihu    string `json:"zhihu"`
	Qq       string `json:"qq"`
	Bilibili string `json:"bilibili"`
	Wechat   string `json:"wechat"`
}

type UserConfig struct {
	ArticleTime int `json:"articleTime"`
	MessageTime int `json:"messageTime"`
	VideoTime   int `json:"videoTime"`
}
