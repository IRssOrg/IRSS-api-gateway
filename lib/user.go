package lib

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       int64  `json:"id"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type LoginResq struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type RegisterResp struct {
	StatusCode int    `json:"statusCode"`
	StatusMsg  string `json:"statusMsg"`
	User       User   `json:"user"`
}
