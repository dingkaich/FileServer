package httpsvr

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type loginObj struct {
	Name string
	HTTPObj
}

//var Loginobj = NewLoginObj()
var Loginobj = &loginObj{}

func (h *loginObj) ObjHanlde(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.Body)
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

func init() {
	Loginobj.Name = "index"
	Loginobj.Res = make(map[string]string, 32)
	Loginobj.RegistObj("/", "./htmlfile/login.html")
	Loginobj.RegistObj("/login", "./htmlfile/login.html")
	Loginobj.RegistObj("/favicon.ico", "./htmlfile/favicon.ico")
	Loginobj.RegistObj("/login/favicon.ico", "./htmlfile/favicon.ico")
	Loginobj.RegistObj("/logindata.js", "./htmlfile/logindata.js")
	Loginobj.RegistObj("/login/logindata.js", "./htmlfile/logindata.js")
}

func (h *loginObj) Get(res http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	if url != "/" {
		url = strings.TrimSuffix(url, "/")
	}

	if v, ok := h.Res[url]; ok {
		// v = "./htmlfile" + v
		f, err := os.Open(v)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = io.Copy(res, f)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	log.Println("unkown res:", url)
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "{\nResult:\"unkown url[%s],please use xxx.xx.com/login\"\n}", url)
	return
}

func (h *loginObj) Head(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Get Head Okay!!"
		}`))
	return
}

//login页面当前只是实现了简单的通过表单传输的密码
func (h *loginObj) Post(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.FormValue("inputEmail")
	passwd := req.FormValue("inputPassword")
	cookie := &http.Cookie{}
	if len(req.Cookies()) != 0 {
		cookie = req.Cookies()[0]
	}

	log.Println("cookies:", cookie)
	err := AuthUsrPwd(username, passwd, cookie)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("{\"Result\":\"%v\"}", err)))
		return
	}
	http.SetCookie(res, cookie)
	log.Println("set cookie ", cookie)
	//http.Redirect(res, req, "/views/"+username, http.StatusFound)
	http.Redirect(res, req, "/login", http.StatusFound)
	return
}
func (h *loginObj) Delete(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Noting to Delete!!"
		}`))
	return
}

func (h *loginObj) Put(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Please use PUT method!!!"
		}`))
	return
}
