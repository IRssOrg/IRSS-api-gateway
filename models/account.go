package models

type AddAccountReq struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type Accounts struct {
	Accounts []AddAccountReq `json:"accounts"`
}
