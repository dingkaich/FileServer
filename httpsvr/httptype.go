/*
本文件主要定义一些抽象概念。
*/

package httpsvr

import "net/http"

//HttpObj 闲着无聊
type HTTPObj struct {
	Res  map[string]string //资源对象:主要是用于URL访问的文件名
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
