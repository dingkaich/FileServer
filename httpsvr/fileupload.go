package httpsvr

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//注意一下语法只有golang 1.9.2存在
type FileUploadObj struct {
	Name string
	HTTPObj
}

func (h *FileUploadObj) ObjHanlde(res http.ResponseWriter, req *http.Request) {
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
var UploadObject = &FileUploadObj{}

func init() {
	UploadObject.Name = "Fileupload"
	UploadObject.Res = make(map[string]string, 32)
	UploadObject.RegistObj("/", "./htmlfile/uploadfile.html")
}

//查询内容
func (h *FileUploadObj) Get(res http.ResponseWriter, req *http.Request) {
	//校验用户是否正确
	var (
		err      error
		url      string
		username string
		userpath string
		file     os.FileInfo
		tmpl     *template.Template
		data     = map[string]interface{}{} //html模板转换公式
	)

	url = req.URL.String()
	err, username = CookiesAuth(res, req)
	if err != nil || username == "" {
		log.Println("cookies auth failed err=", err, "username=", username)
		goto loginfailed
	}

	//只允许访问自己的路径
	userpath = "/uploadfile/" + username
	if !strings.HasPrefix(url, userpath) {
		fmt.Printf("wrong url[%s] uploadfile[%s] ", url, userpath)
		result := []byte("{\"Reslt\":\"please upload in yourself directory\"}")
		res.Write(result)
		return
	}

	file, err = os.Stat("./upload" + username)
	if err != nil {
		os.Mkdir("./upload"+username, 0666)
	} else {
		if file.IsDir() {
			os.Remove("./upload" + username)
		}
	}
	data = map[string]interface{}{"Username": username}
	tmpl, err = template.ParseFiles("./htmlfile/uploadfile.html")
	if err != nil {
		res.Write([]byte("get upload html error"))
		http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(res, data)
	return

loginfailed:
	result := []byte("{\n\"Reslt\":\"not login;please re-login\"\n}")
	res.Write(result)

	return
}

func (h *FileUploadObj) Head(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Get Head Okay!!"
		}`))
	return
}

//login页面当前只是实现了简单的通过表单传输的密码
func (fp *FileUploadObj) Post(res http.ResponseWriter, req *http.Request) {

	//校验用户是否正确
	var (
		err      error
		url      string
		username string
		userpath string
	)

	url = req.URL.String()
	err, username = CookiesAuth(res, req)
	if err != nil || username == "" {
		log.Println("cookies auth failed err=", err, "username=", username)
		result := []byte("{\n\"Reslt\":\"not login;please re-login\"\n}")
		res.Write(result)
		return
	}

	//只允许访问自己的路径
	userpath = "/uploadfile/" + username
	if !strings.HasPrefix(url, userpath) {
		fmt.Printf("wrong url[%s] uploadfile[%s] ", url, userpath)
		result := []byte("{\"Reslt\":\"please upload in yourself directory\"}")
		res.Write(result)
		return
	}

	fileExist, err := os.Stat("./upload" + username)
	if err != nil {
		os.Mkdir("./upload"+username, 0666)
	} else {
		if fileExist.IsDir() {
			os.Remove("./upload" + username)
		}
	}

	// goto work
	req.ParseMultipartForm(1 << 10) //设置下内存缓存大小，默认为1G

	f, h, err := req.FormFile("uploadfile") //该值有html中定义
	// defer f.Close()
	if err != nil || f == nil || h == nil {
		// if req.Form
		fmt.Fprintln(res, "uploadfile", func(i *multipart.FileHeader) string {
			if i == nil {
				return "N/A"
			} else {
				return i.Filename
			}
		}(h), " failed")
		return
	}
	file, err := os.OpenFile("./upload/"+username+"/"+h.Filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		fmt.Fprintln(res, "create ", h.Filename, " in server failed")
		http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(file, f)
	if err != nil {
		fmt.Fprintln(res, "write ", h.Filename, " in server failed")
		http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
	filedir, _ := filepath.Abs("./myfileserver/upload/" + h.Filename)
	filesize, _ := file.Seek(0, 1)
	file_res := FileOpStat{
		Filename:     h.Filename,
		Filepath:     filedir,
		Uploadstatus: "upload success",
		HttpStatus:   http.StatusOK,
		Description:  fmt.Sprintf("%s%d bytes", h.Filename+"上传完成，文件大小:", filesize),
	}
	// resdata, _ := json.Marshal(file_res)
	resdata, _ := json.MarshalIndent(file_res, "", "\t")

	res.Write(resdata)

	return
}
func (h *FileUploadObj) Delete(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Noting to Delete!!"
		}`))
	return
}

func (h *FileUploadObj) Put(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Please use PUT method!!!"
		}`))
	return
}

type FileOpStat struct {
	Filename     string
	Filepath     string
	Uploadstatus string
	HttpStatus   int
	Description  string
}
