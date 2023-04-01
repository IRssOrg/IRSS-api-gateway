package controller

import (
	"connection-gateway/lib"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	pool    *sql.DB
	connStr = "root:zkw030813@tcp(101.43.168.188:3306)/public"
)

func Init() error {
	var err error
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config lib.Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		log.Fatal(err)
	}
	connStr = config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + config.Database.DatabaseName
	log.Println(connStr)
	pool, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	return nil
}
