package httpsvr

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

//注意一下语法只有golang 1.9.2存在
type FileViewObj struct {
	Name string
	HTTPObj
}

func (h *FileViewObj) ObjHanlde(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.Post(res, req)
	case http.MethodGet:
		h.Get(res, req)
	case http.MethodHead:
		h.Head(res, req)
	case http.MethodDelete:
		h.Delete(res, req)
	case http.MethodPut:
		h.Put(res, req)
	}
}

//var ViewObject = io.Closer()
var ViewObject = &FileViewObj{}

func init() {
	ViewObject.Name = "Fileview"
	ViewObject.Res = make(map[string]string, 32)
	ViewObject.RegistObj("/", "./upload")
}

//查询内容
func (h *FileViewObj) Get(res http.ResponseWriter, req *http.Request) {
	//校验用户是否正确
	var (
		err      error
		url      string
		username string
		userpath string
		file     os.FileInfo
	)

	url = req.URL.String()
	err, username = CookiesAuth(res, req)
	if err != nil || username == "" {
		log.Println("cookies auth failed err=", err, "username=", username)
		goto loginfailed
	}

	//只允许访问自己的路径
	userpath = "/viewfile/" + username
	if !strings.HasPrefix(url, userpath) {
		fmt.Printf("wrong url[%s] viewpath[%s] ", url, userpath)
		result := []byte("{\"Reslt\":\"please view yourself directory\"}")
		res.Write(result)
		return
	}

	file, err = os.Stat("./upload" + username)
	if err != nil {
		os.Mkdir("./upload/"+username, os.ModePerm)
	} else {
		if !file.IsDir() {
			os.Remove("./upload/" + username)
		}
	}

	//目录的话直接浏览
	if strings.HasSuffix(url, "/") {
		log.Printf("view path=%s", userpath)
		http.StripPrefix(userpath, http.FileServer(http.Dir("./upload/"+username))).ServeHTTP(res, req)
	} else {
		log.Printf("download file=%s", url)
		//文件的话直接下载
		res.Header().Set("Content-Type", "application/octet-stream") //设置文件下载类型
		// http.FileServer(http.Dir("./upload/"+username)).ServeHTTP(res, req)
		// http.FileServer(http.Dir("./")).ServeHTTP(res, req)
		http.StripPrefix(userpath, http.FileServer(http.Dir("./upload/"+username))).ServeHTTP(res, req)
		log.Println(req.Header.Get("Content-Type"), "|", res.Header().Get("Content-Type"))
	}
	return

loginfailed:
	result := []byte("{\n\"Reslt\":\"not login;please re-login\"\n}")
	res.Write(result)

	return
}

func (h *FileViewObj) Head(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Get Head Okay!!"
		}`))
	return
}

//login页面当前只是实现了简单的通过表单传输的密码
func (h *FileViewObj) Post(res http.ResponseWriter, req *http.Request) {
	return
}
func (h *FileViewObj) Delete(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Noting to Delete!!"
		}`))
	return
}

func (h *FileViewObj) Put(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Please use PUT method!!!"
		}`))
	return
}
