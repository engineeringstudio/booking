package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	EnableTLS  bool     `json:"enabletls"`
	CertPath   string   `json:"certpath"`
	KeyPath    string   `json:"keypath"`
	Port       string   `json:"port"`
	Mail       string   `json:"mail"`
	MailServer string   `json:"mailserver"`
	Passwd     string   `json:"passwd"`
	MaxLength  int      `json:"maxlength"`
	DataBase   string   `json:"database"`
	WhiteList  []string `json:"whitelist"`
	MailList   []string `json:"maillist"`
}

func ReadConf(path string, conf *Config) error {
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
