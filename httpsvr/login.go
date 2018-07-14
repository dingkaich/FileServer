package httpsvr

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var getuserinfo func(username string) []byte

func RegistGetUserInfo(fun func(username string) []byte) {
	getuserinfo = fun
}

var Loginobj = NewLoginObj()

func init() {
	Loginobj.Name = "index"
	Loginobj.RegistObj("/", "./htmlfile/login.html")
	Loginobj.RegistObj("/login", "./htmlfile/login.html")
	Loginobj.RegistObj("/favicon.ico", "./htmlfile/favicon.ico")
	Loginobj.RegistObj("/login/favicon.ico", "./htmlfile/favicon.ico")
	Loginobj.RegistObj("/logindata.js", "./htmlfile/logindata.js")
	Loginobj.RegistObj("/login/logindata.js", "./htmlfile/logindata.js")
}

func (h *HTTPObj) Get(res http.ResponseWriter, req *http.Request) {
	url := strings.TrimSuffix(req.URL.String(), "/")
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

	log.Println("unkown res:", url)
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(res, "{\nResult:\"unkown url[%s],please use xxx.xx.com/login\"\n}", url)
	return
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

	if storpasswd == nil {
		res.Write([]byte(`
{
	"Result":"Login failed!! unkown username"
}`))
		return
	}

	log.Printf("get db usrname[%s],passwd[%s]", username, string(storpasswd))
	aa := sha1.Sum(storpasswd)
	passwd := req.FormValue("inputPassword")
	bytepasswd, err := hex.DecodeString(passwd)
	if err != nil {
		res.Write([]byte(`
{
	"Result":"Login failed!! wrong passwd"
}`))
		return
	}

	if bytes.Compare(aa[0:], bytepasswd) == 0 {
		fmt.Println("login success :", username)
	} else {
		fmt.Println("login failed :", username)
		res.Write([]byte(`
{
	"Result":"Login failed!!"
}
			`))
		return
	}

	http.Redirect(res, req, "/views/"+username, http.StatusFound)
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
