package main

import (
	"FileServer/httpsvr"
	mysql "FileServer/sqlite"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func StartHTTP(localconf *config_parm) {
	// log.SetFlags(flag)
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	//注册获取用户名和密码的地方
	httpsvr.RegistGetUserInfo(mysql.QueryUserInfo)

	if localconf == nil {
		log.Println("localconf  is nil")
		return
	}
	fmt.Println("currentdir=", os.Args[0])
	Mux := http.NewServeMux()
	//根目录是主页
	Mux.HandleFunc("/", httpsvr.Loginobj.ObjHanlde)
	Mux.HandleFunc("/login/", httpsvr.Loginobj.ObjHanlde)
	Mux.HandleFunc("/index", httpsvr.IndexObject.ObjHanlde)
	Mux.HandleFunc("/index/", httpsvr.IndexObject.ObjHanlde)
	// Mux.HandleFunc("/index/", httpsvr.Loginobj.ObjHanlde)
	Mux.HandleFunc("/viewfile", httpsvr.ViewObject.ObjHanlde)
	Mux.HandleFunc("/viewfile/", httpsvr.ViewObject.ObjHanlde)
	Mux.HandleFunc("/uploadfile", httpsvr.UploadObject.ObjHanlde)
	Mux.HandleFunc("/uploadfile/", httpsvr.UploadObject.ObjHanlde)

	server := http.Server{
		Addr:         localconf.Ip + ":" + localconf.HttpPort,
		Handler:      Mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 20,
	}
	fmt.Println("start to listen", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
