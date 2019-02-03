package main

import (
	mysql "FileServer/sqlite"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type config_parm struct {
	Ip       string `json:"ip"`
	HttpPort string `json:"httpport"`
	SftpPort string `json:"sftpport"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

var default_cfg = &config_parm{
	Ip:       "localhost",
	HttpPort: "6060",
	SftpPort: "20022",
	Username: "dingkai",
	Passwd:   "12345",
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
	log.Println("set by default configuration")
	return default_cfg

}

var localconf *config_parm

func StartFilesever(localconf *config_parm) {
	go StartHTTP(localconf) //http
	go StartSFTP(localconf) //sftp

}

func main() {
	//1. 获取基本配置
	localconf = loadconfig()

	//2.初始化数据
	mysql.DBinit()

	//3. 启动文件服务器
	StartFilesever(localconf)

	select {}
}
