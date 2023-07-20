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

	mail := newMailSender(conf.Mail, conf.Mail, conf.Passwd, conf.MailServer)

	tech = NewHandler("tech", db, mail)

	for i := 0; i < len(conf.WhiteList); i++ {
		tech.whitelist[conf.WhiteList[i]] = struct{}{}
	}

	http.HandleFunc("/add", tech.add)
	http.HandleFunc("/sand", tech.send)
}

func main() {
	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
	}
}
