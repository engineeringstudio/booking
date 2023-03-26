package main

import (
	"encoding/json"
	"os"
)

type config struct {
	EnableTLS bool     `json:"enabletls"`
	CertPath  string   `json:"certpath"`
	KeyPath   string   `json:"keypath"`
	Port      string   `json:"port"`
	MaxLength int      `json:"maxlength"`
	DataBase  string   `json:"database"`
	WhiteList []string `json:"whitelist"`
	MailList  []string `json:"maillist"`
}

func readConf(path string, conf *config) error {
	_, err := os.Stat(path)
	if err != nil || os.IsExist(err) {
		return err
	}
	tmp, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, conf)
}
