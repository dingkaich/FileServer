package fileredis

import (
	"fmt"

	//"errors"

	"time"

	"github.com/go-redis/redis"
)

type Fileinfo struct {
	Filename       string
	Filemd5        string
	Fileserverpath string
	ServerIp       string
	ClientIp       string
	Date           string
}

type Userinfo struct {
	Username string
	Passwd   string
}

var client *redis.Client

func DBinit() {
	client = redis.NewClient(&redis.Options{
		//Addr:     "192.168.0.104:6379",
		Addr:     "106.14.179.186:6379",
		Password: "dingkai", // no password set
		DB:       0,         // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
		return
	}

	err = AddUserInfo("dingkai", "123456")
	if err != nil {
		fmt.Println("add user info failed")
	}
	return
}

// type Userinfo struct {
// 	Username string
// 	Passwd   string
// }
// set username passwd

func QueryUserInfo(username string) []byte {

	val, err := client.Get(username).Result()
	if err != nil {
		return nil
	}
	return []byte(val)
}

func AddUserInfo(username, passwd string) error {
	err := client.SetNX(username, passwd, time.Duration(0)).Err()
	if err != nil {
		return err
	}
	return nil
}

// type Fileinfo struct {
// 	Filename       string
// 	Filemd5        string
// 	Fileserverpath string
// 	ServerIp       string
// 	ClientIp       string
// 	Date           string
// }
// hash
// hmset Filename  Filemd5 xx  Fileserverpath xx  ServerIp xx  ClientIp xx  Date xx

func Addfile(data *Fileinfo) error {

	fields := make(map[string]interface{})
	fields["Filemd5"] = data.Filemd5
	fields["Fileserverpath"] = data.Fileserverpath
	fields["ServerIp"] = data.ServerIp
	fields["ClientIp"] = data.ClientIp
	fields["Date"] = data.Date

	client.HMSet(data.Filename, fields)
	return nil
}

func Queryfile(filename string) ([]Fileinfo, error) {
	val, err := client.HGetAll(filename).Result()
	if err != nil || val == nil || len(val) == 0 {
		return nil, err
	}

	//暂时认为到这里是必然有值的
	var filebean Fileinfo
	filebean.Filename = filename
	filebean.Filemd5 = val["Filemd5"]
	filebean.Fileserverpath = val["Fileserverpath"]
	filebean.ClientIp = val["ClientIp"]
	filebean.ServerIp = val["ServerIp"]
	filebean.Date = val["Date"]

	var result []Fileinfo
	result = append(result, filebean)
	return result, nil

}

func Deletefile(filename string) error {

	_, err := client.Del(filename).Result()
	if err != nil {
		return err
	}
	return nil
}

func Updatefile(data *Fileinfo) error {
	//更新的话，加上一个强校验，必须要先存在这个值
	_, err := client.Exists(data.Filename).Result()
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})
	fields["Filemd5"] = data.Filemd5
	fields["Fileserverpath"] = data.Fileserverpath
	fields["ServerIp"] = data.ServerIp
	fields["ClientIp"] = data.ClientIp
	fields["Date"] = data.Date

	client.HMSet(data.Filename, fields)
	return nil
}
