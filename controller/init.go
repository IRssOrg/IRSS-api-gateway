package controller

import (
	"database/sql"
	"encoding/json"
	"io"
	"irss-gateway/models"
	"log"
	"os"
	"strconv"
)

var pool *sql.DB

func Init() error {
	var err error
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config models.Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		log.Fatal(err)
	}
	connStr := config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + config.Database.DatabaseName
	log.Println(connStr)
	pool, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	return nil
}
