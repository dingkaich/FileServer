package httpsvr

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
)

var getuserinfo func(username string) []byte = nil

func RegistGetUserInfo(fun func(username string) []byte) {
	getuserinfo = fun
}

func NewLoginObj() *HTTPObj {
	return &HTTPObj{
		Res: make(map[string]string, 32),
	}
}

func (h *HTTPObj) RegistObj(urlname, fullfilename string) {
	h.Res[urlname] = fullfilename
}

func (h *HTTPObj) Get(res http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
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
	}
}

func (h *HTTPObj) Head(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Get Head Okay!!"
		}`))
	return
}

//login页面当前只是实现了简单的通过表单传输的密码
func (h *HTTPObj) Post(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	username := req.FormValue("inputEmail")
	storpasswd := getuserinfo(username)
	aa := sha1.Sum(storpasswd)
	passwd := req.FormValue("inputPassword")
	if bytes.Compare(aa[0:], []byte(passwd)) == 0 {
		fmt.Println("login success :", username)

	} else {
		fmt.Println("login failed", storpasswd, "&&", aa)
		res.Write([]byte(`
			{
				"Result":"Login failed!!"
			}
			`))
		return
	}
	http.Redirect(res, req, "/views/"+"username", http.StatusFound)
	return
}
func (h *HTTPObj) Delete(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Noting to Delete!!"
		}`))
	return
}

func (h *HTTPObj) Put(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte(`
		{
			"Result":"Please use PUT method!!!"
		}`))
	return
}

func (h *HTTPObj) ObjHanlde(res http.ResponseWriter, req *http.Request) {
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

var Loginobj = NewLoginObj()

func init() {
	Loginobj.Name = "index"
	Loginobj.RegistObj("/", "./htmlfile/login.html")
	Loginobj.RegistObj("/favicon.ico", "./htmlfile/favicon.ico")
	Loginobj.RegistObj("/logindata.js", "./htmlfile/logindata.js")
}
