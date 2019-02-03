package main

import "FileServer/sftpsvr"

func StartSFTP(localconf *config_parm) {
	sftpsvr.SftpServer(localconf.Ip + ":" + localconf.SftpPort)
}
