/*用户名和密码校验
登录界面输入用户名和密码后，做校验。
a. 校验okay: 返回cookie:
	{
		"devid":设备id //用于分布式之间存储用
		"sessionid":认证key //校验通过后，随机生成的信息，后期完全采用该值来做计算。
		“lifetime”: 30 //cookie过期时间，单位分钟
	}
b. 校验不通过，返回校验失败

sessionid 生成规则:md5sum(linux时间戳 + md5(用户名+密码))
每次校验 seesionid中的linux时间戳+ lifetime ,如果过期了；要求重新设置登录；返回超期登录。

*/

package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type AuthUser struct {
	Username string
	Passwd   string
	Cookies  *http.Cookie
}

const (
	authokay      = 0 //鉴权通过，返回该用户的文件夹列表
	authfailed    = 1 //鉴权失败，请求重新登录
	cookietimeout = 2 //cookie超时，让重新登录
)

var getpaswd func(username string) []byte

func RegistGetPasswdFunc(fun func(username string) []byte) {
	getpaswd = fun
}

var sessionmap map[string]*AuthUser

func init() {
	getpaswd = nil
	sessionmap = make(map[string]*AuthUser, 32)
}

func NewAuthUser(username string) *AuthUser {

	if username == "" {
		return &AuthUser{
			Username: "",
			Passwd:   "",
			Cookies:  nil,
		}
	} else {
		if a, ok := sessionmap[username]; ok {
			return a
		}
	}
	return nil
}

func (a *AuthUser) Auth() int {
	if a.Cookies != nil {
		return a.authCookie()
	}
	return a.authuser()

}

func (a *AuthUser) authuser() int {
	//对用户名和密码做解密操作
	if a.Username == "" || a.Passwd == "" || getpaswd == nil {
		return authfailed
	}

	if string(getpaswd(a.Username)) != a.Passwd {
		return authfailed
	}

	//设置cookie
	a.setcookie()

	sessionmap[a.Username] = a
	return authokay
}

func (a *AuthUser) setcookie() {
	nowtime := time.Now()
	sessionstring := fmt.Sprintf("%d%s%s", nowtime.Unix(), a.Username, a.Passwd)
	sessionbyte := md5.New().Sum([]byte(sessionstring))
	sessionid := base64.StdEncoding.EncodeToString(sessionbyte)
	a.Cookies = &http.Cookie{
		Name:    fmt.Sprintf("%d", nowtime.Unix()),
		Value:   sessionid,
		Expires: time.Now().Add(time.Hour * 1),
	}

}

func (a *AuthUser) authCookie() int {
	if a.Cookies == nil {
		return authfailed
	}

	if a.Cookies.Name != a.Username {
		return authfailed
	}

	sessionstring := fmt.Sprintf("%d%s%s", a.Cookies.Value, a.Username, a.Passwd)

	if sessionstring != a.Cookies.Value {
		return a.authuser()
	}

	return authfailed
}
