package controller

import (
	"database/sql"
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"irss-gateway/models"
	"log"
	"os"
	"strconv"
)

var pool *sql.DB
var db *gorm.DB

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
	db, err = gorm.Open(mysql.Open(connStr), &gorm.Config{})
	pool, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	return nil
}
