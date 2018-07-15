package httpsvr

import (
	"net/http"
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
	ViewObject.RegistObj("/", "sd")
}

//查询内容
func (h *FileViewObj) Get(res http.ResponseWriter, req *http.Request) {
	//校验用户是否正确

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
