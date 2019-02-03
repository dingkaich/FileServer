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

package sftpsvr

import (
	"errors"
	"log"
)

type AuthUser struct {
	Username string //用户名
	Passwd   string //hex sha1的值
}

var (
	Errauthuserfailed = errors.New("wrong user")
	Errauthpwdfailed  = errors.New("wrong passwd")
)

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

//sftp传来的密码都是明文
func AuthUsrPwd(username, passwd string) error {

	//数据库里存的是明文
	storpasswd := getuserinfo(username)
	if storpasswd == nil {
		log.Println("unknow username:", username)
		return Errauthuserfailed
	}

	if passwd == string(storpasswd) {
		return nil
	} else {
		return Errauthpwdfailed
	}

}
