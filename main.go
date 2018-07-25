package main

import (
	mysql "FileServer/sqlite"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config_parm struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

func loadconfig() *config_parm {
	conf_bean := &config_parm{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		goto deafultval
	}

	err = json.Unmarshal(data, conf_bean)
	if err != nil {
		fmt.Println("unmarshal failed")
		goto deafultval
	}

	return conf_bean

deafultval:
	return &config_parm{
		Ip:       "localhost",
		Port:     "6060",
		Username: "dingkai",
		Passwd:   "12345",
	}

}

var localconf *config_parm

func main() {
	mysql.DBinit()
	localconf = loadconfig()
	StartFilesever(localconf)

}
