package httpsvr

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type FileIndexObj struct {
	Name string
	HTTPObj
}

func (h *FileIndexObj) ObjHanlde(res http.ResponseWriter, req *http.Request) {
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
var IndexObject = &FileIndexObj{}

func init() {
	IndexObject.Name = "index"
	IndexObject.Res = make(map[string]string, 32)
	IndexObject.RegistObj("/", "./htmlfile/index.html")
	IndexObject.RegistObj("/index", "./htmlfile/index.html")
	IndexObject.RegistObj("/favicon.ico", "./htmlfile/favicon.ico")
}

//查询内容
func (h *FileIndexObj) Get(res http.ResponseWriter, req *http.Request) {
	//校验用户是否正确
	var (
		err      error
		url      string
		username string
		tmpl     *template.Template
		data     = map[string]interface{}{} //html模板转换公式
	)
	//data := map[string]interface{}{"Username": username}
	err, username = CookiesAuth(res, req)

	if err != nil || username == "" {
		log.Println("cookies auth failed err=", err, "username=", username)
		goto loginfailed
	}

	url = req.URL.String()
	if url != "/" {
		data = map[string]interface{}{"Username": username}
		url = strings.TrimSuffix(url, "/")
	}

	//静态文件获取，同时做template的转换
	if v, ok := h.Res[url]; ok {
		tmpl, err = template.ParseFiles(v) //将一个文件读作模板

		tmpl.Execute(res, data)
		log.Println("get res name:", tmpl.Name())
		return //找到资源就直接退出好了
	}

	log.Println("unkown res:", url)
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "{\nResult:\"unkown url[%s],please use xxx.xx.com/login\"\n}", url)

loginfailed:
	result := []byte("{\n\"Reslt\":\"not login;please re-login\"\n}")
	res.Write(result)
	http.RedirectHandler("/login", http.StatusOK)

	return
}

func (h *FileIndexObj) Head(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Get Head Okay!!"
		}`))
	return
}

//login页面当前只是实现了简单的通过表单传输的密码
func (h *FileIndexObj) Post(res http.ResponseWriter, req *http.Request) {
	return
}
func (h *FileIndexObj) Delete(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Noting to Delete!!"
		}`))
	return
}

func (h *FileIndexObj) Put(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Please use PUT method!!!"
		}`))
	return
}
