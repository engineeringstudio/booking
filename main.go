package main

import (
	"database/sql"
	"flag"
	"net/http"

	utils "booking/internel"
)

var conf utils.Config
var confpath string

var tech *utils.Handler

func init() {
	flag.StringVar(&confpath, "c", "./config.json", "Set the config path")

	flag.Parse()

	err := utils.ReadConf(confpath, &conf)
	if err != nil {
		panic("OpenConfigError")
	}

	db, err := sql.Open("sqlite3", conf.DataBase)
	if err != nil {
		panic("OpenDatabaseError")
	}

	mail := utils.NewMailSender(conf.MailList, conf.Mail, conf.Mail, conf.Passwd, conf.MailServer)

	tech = utils.NewHandler("tech", conf.WhiteList, conf.MaxLength, db, mail)

	http.HandleFunc("/add", tech.Add)
	http.HandleFunc("/sand", tech.Send)
}

func main() {
	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
	}
}
