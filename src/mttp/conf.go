package main

import (
	"database/sql"
	"fmt"

	"github.com/xrfang/go-conf"
)

type Configuration struct {
	ALLOW_WRITE   bool
	DB_USER       string
	DB_PASS       string
	DB_HOST       string
	DB_PORT       string
	DB_NAME       string
	HTTP_PORT     string
	READ_LIMIT    int
	READ_TIMEOUT  int
	WRITE_TIMEOUT int
}

var (
	cf Configuration
	db *sql.DB
)

func loadConfig(fn string) {
	//default values
	cf.HTTP_PORT = "6887"
	cf.READ_TIMEOUT = 60
	cf.WRITE_TIMEOUT = 60
	cf.READ_LIMIT = 1000
	err := conf.ParseFile(fn, &cf)
	assert(err)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&sql_mode=ANSI_QUOTES",
		cf.DB_USER, cf.DB_PASS, cf.DB_HOST, cf.DB_PORT, cf.DB_NAME)
	db, err = sql.Open("mysql", dsn)
	assert(err)
}
