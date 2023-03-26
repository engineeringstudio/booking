package main

import (
	"database/sql"
	"flag"
	"net/http"
)

var conf config
var confpath string

var tech *handler

func init() {
	flag.StringVar(&confpath, "c", "./config.json", "Set the config path")

	flag.Parse()

	err := readConf(confpath, &conf)
	if err != nil {
		panic("OpenConfigError")
	}

	db, err := sql.Open("sqlite3", conf.DataBase)
	if err != nil {
		panic("OpenDatabaseError")
	}

	tech = NewHandler("tech", db)

	for i := 0; i < len(conf.WhiteList); i++ {
		tech.whitelist[conf.WhiteList[i]] = struct{}{}
	}

	http.HandleFunc("/"+tech.name, tech.add)
}

func main() {
	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
	}
}
