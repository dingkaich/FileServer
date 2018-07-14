/*
本文件主要定义一些抽象概念。
主要用于定义静态文件的http访问模型

*/

package httpsvr

import "net/http"

// HTTPObj  闲着无聊
//主要适用于静态文件的对象
type HTTPObj struct {
	Res  map[string]string //资源对象:主要是用于URL访问的静态文件名
	Name string
}

//HTTPHandle ...
type HTTPHandle interface {
	Get(res http.ResponseWriter, req *http.Request)
	Post(res http.ResponseWriter, req *http.Request)
	Put(res http.ResponseWriter, req *http.Request)
	Head(res http.ResponseWriter, req *http.Request)
	Delete(res http.ResponseWriter, req *http.Request)
}

func NewLoginObj() *HTTPObj {
	return &HTTPObj{
		Res: make(map[string]string, 32),
	}
}

func (h *HTTPObj) RegistObj(urlname, fullfilename string) {
	h.Res[urlname] = fullfilename
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
