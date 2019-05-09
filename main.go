package main

import (
	"flag"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
)

type shortyContext struct {
	Config *config
	DB     *gorm.DB
}

func main() {
	confPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	conf, err := loadConfig(*confPath)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("sqlite3", conf.DatabasePath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Link{})

	ctx := shortyContext{conf, db}
	a := api{&ctx}

	router := httprouter.New()
	a.route(router)

	err = http.ListenAndServe(conf.ListenAddress, router)
	if err != nil {
		panic(err)
	}
}
