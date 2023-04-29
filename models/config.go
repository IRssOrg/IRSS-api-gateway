package models

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	Database struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		User         string `json:"user"`
		Password     string `json:"password"`
		DatabaseName string `json:"database_name"`
	} `json:"database"`
	QQTopics struct {
		Topics []string `json:"topics"`
	}
	Spider struct {
		Zhihu    string `json:"zhihu"`
		Wechat   string `json:"wechat"`
		Bilibili string `json:"bilibili"`
	}
}
