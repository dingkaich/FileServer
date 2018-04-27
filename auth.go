/*用户名和密码校验
登录界面输入用户名和密码后，做校验。
a. 校验okay: 返回cookie:
	{
		"devid":设备id //用于分布式之间存储用
		"sessionid":认证key //校验通过后，随机生成的信息，后期完全采用该值来做计算。
		“lifetime”: 30 //cookie过期时间，单位分钟
	}
b. 校验不通过，返回校验失败

sessionid 生成规则: linux时间戳 + md5(用户名+密码)
每次校验 seesionid中的linux时间戳+ lifetime ,如果过期了；要求重新设置登录；返回超期登录。

*/

package main

import "net/http"

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

func NewAuthUser() *AuthUser {
	return &AuthUser{
		Username: "",
		Passwd:   "",
		Cookies:  nil,
	}
}

func (a *AuthUser) Auth() int {
	if a.Cookies != nil {
		return a.authCookie()
	}
	return a.authuser()

}

func (a *AuthUser) authuser() int {
	return authfailed
}

func (a *AuthUser) authCookie() int {
	return authfailed
}
