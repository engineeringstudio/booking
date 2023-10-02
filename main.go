package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	"booking/internel"
)

var conf internel.Config
var confpath string

var tech *internel.Handler

func init() {
	flag.StringVar(&confpath, "c", "./config.json", "Set the config path")

	flag.Parse()

	err := internel.ReadConf(confpath, &conf)
	if err != nil {
		fmt.Println(err)
		panic("OpenConfigError")
	}

	db, err := sql.Open("sqlite3", conf.DataBase)
	if err != nil {
		fmt.Println(err)
		panic("OpenDatabaseError")
	}

	mail := internel.NewMailSender(conf.MailList, conf.Mail, conf.Mail, conf.Passwd, conf.MailServer)

	tech, err = internel.NewHandler(&conf, db, mail)
	if err != nil {
		fmt.Println(err)
		panic("OpenHandlerErr")
	}

	go tech.UpdateTask()
	http.HandleFunc("/add", tech.Add)
}

func main() {
	if conf.EnableTLS {
		http.ListenAndServeTLS("0.0.0.0:"+conf.Port,
			conf.CertPath, conf.KeyPath, nil)
	} else {
		http.ListenAndServe("0.0.0.0:"+conf.Port, nil)
	}
}
