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

package httpsvr

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type AuthUser struct {
	Username string //用户名
	Passwd   string //hex sha1的值
	Cookies  *http.Cookie
}

var (
	Errauthuserfailed   = errors.New("wrong user")
	Errauthpwdfailed    = errors.New("wrong passwd")
	ErrauthCookiefailed = errors.New("cookie timeout")
	Errauthunkownfailed = errors.New("wrong unkonw reson")
)

var sessionmap map[string]*AuthUser //username-authinfo
var sessionRWLock sync.RWMutex

func checksessionmap() {

	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			log.Println("checkseesion map")
			sessionRWLock.Lock()
			for k, v := range sessionmap {
				if v.Cookies.Expires.After(time.Now()) {
					log.Println("delete timeout seesion username:", k)
					delete(sessionmap, k)
				}
			}
			sessionRWLock.Unlock()
		}
	}

}

func init() {
	sessionRWLock.Lock()
	sessionmap = make(map[string]*AuthUser, 32)
	sessionRWLock.Unlock()
}

//外部传递的获取用户名和密码的封装函数
var getuserinfo func(username string) []byte

func RegistGetUserInfo(fun func(username string) []byte) {
	getuserinfo = fun
}

func newauthuser(username, passwd string) *AuthUser {
	return &AuthUser{
		Username: username,
		Passwd:   passwd,
	}
}

func AuthUsrPwd(username, passwd string, cookie *http.Cookie) error {

	//find user,存在cookie，直接校验cookie
	var err error
	sessionRWLock.RLock()
	if value, ok := sessionmap[username]; ok {

		if err = value.authCookie(cookie); err == nil {
			sessionRWLock.RUnlock()
			log.Printf("user[%s] cookie auth okay", username)
			return nil
		}
	}
	sessionRWLock.RUnlock()

	//cookie校验失败，直接校验用户名
	storpasswd := getuserinfo(username)
	if storpasswd == nil {
		log.Println("unknow username:", username)
		return Errauthuserfailed
	}
	aa := sha1.Sum(storpasswd)

	bytepasswd, err := hex.DecodeString(passwd)
	if err != nil {
		return Errauthpwdfailed
	}

	if bytes.Compare(aa[0:], bytepasswd) == 0 {
		log.Println("login success :", username)
	} else {
		fmt.Println("login failed :", username)
		return Errauthpwdfailed
	}
	//用户名校验Okay,写入cookie

	sessionRWLock.Lock()
	sessionmap[username] = newauthuser(username, passwd) //按照华为的规范这里密码不能存入内存。不管了。呵呵
	*cookie = *sessionmap[username].setcookie()          //把设置的cookie返回上去
	//= sessionmap[username].Cookies
	sessionRWLock.Unlock()
	return nil

}

//平时登录时携带的cookie校验
func AuthUsrCookie(cookie *http.Cookie) error {
	sessionid := cookie.Value
	sessionRWLock.RLock()
	defer sessionRWLock.RUnlock()
	for _, v := range sessionmap {
		//sessionid 有效，未过期，直接返回okay
		if v.Cookies.Value == sessionid && v.Cookies.Expires.After(time.Now()) {
			return nil
		}
	}

	return ErrauthCookiefailed

}

func (a *AuthUser) setcookie() *http.Cookie {
	nowtime := time.Now()
	sessionstring := fmt.Sprintf("%d%s%s", nowtime.Unix(), a.Username, a.Passwd)
	sessionbyte := md5.New().Sum([]byte(sessionstring))
	sessionid := base64.StdEncoding.EncodeToString(sessionbyte)
	a.Cookies = &http.Cookie{
		Name:    fmt.Sprintf("%d", nowtime.Unix()),
		Value:   sessionid,
		Expires: time.Now().Add(time.Hour * 1),
		MaxAge:  60 * 60, //过期1小时
	}
	return a.Cookies
}

func (a *AuthUser) authCookie(cookie *http.Cookie) error {
	if a.Cookies == nil {
		// cookie是空，那么就要delete session map
		log.Println("cookie is null")
		return ErrauthCookiefailed
	}

	if a.Cookies.Name != cookie.Name {
		log.Println("cookie name error ", a.Cookies, "|||", cookie)
		return ErrauthCookiefailed
	}

	if a.Cookies.Value != cookie.Value {
		log.Println("cookie value error")
		return ErrauthCookiefailed
	}

	if a.Cookies.Expires.Before(time.Now()) {
		log.Println("cookie time out")
		return ErrauthCookiefailed
	}

	return nil
}
