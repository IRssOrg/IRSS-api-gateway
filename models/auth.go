package models

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
	Token      string `json:"token"`
	Id         int64  `json:"id"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
	Id         int64  `json:"id"`
}
