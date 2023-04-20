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
