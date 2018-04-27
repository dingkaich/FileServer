package sqlite

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
sqllite 貌似没有行锁等，只有一个库锁，而且还要求应用自己保证锁okay。
这里套接一层，
，也不知道go-sqlite3有没有帮忙实现,打算维持一个fd来开库，并且使用读写锁

先这么滴吧。写一个玩玩吧

*/

/*
create table fileinfo (
    filename varchar(1024) not null CONSTRAINT filename_pk PRIMARY KEY,
    filemd5  varchar(512),
    fileserverpath varchar(1024),
    serverIp    varchar(25),
    clientIp varchar(25)
);
*/

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

var Dbsqlite *sql.DB = nil

var DSN = "./myfileserver/sqlite/filesever.db"

// var DSN = "/root/dingkai/filesever.db"
var SQLfile = "./myfileserver/sqlite/sqlite3.sql"

var rwlock *sync.RWMutex

func DBinit() error {
	//连接db
	var err error
	Dbsqlite, err = sql.Open("sqlite3", DSN)
	if err != nil {
		log.Println(err)
		return err
	}

	/*
		//建表
		sql, err := ioutil.ReadFile(SQLfile)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = Dbsqlite.Exec(string(sql))
		if err != nil {
			log.Println(err, sql)
			return err
		}
	*/
	rwlock = new(sync.RWMutex)
	log.Println("rwlock:", rwlock)
	return nil
}

func DBfinish() {
	Dbsqlite.Close()
	Dbsqlite = nil
}

func QueryUserInfo(username string) []byte {
	rwlock.RLock()
	defer rwlock.RUnlock()
	var passwd []byte
	err := Dbsqlite.QueryRow("select passwd from userinfo where username=?", username).Scan(&passwd)
	if err != nil {
		log.Println(err)
		return nil
	}
	return passwd
}

func Addfile(data *Fileinfo) error {
	rwlock.Lock()
	defer rwlock.Unlock()
	_, err := Dbsqlite.Exec("insert into fileinfo  values(?,?,?,?,?,?)", data.Filename, data.Filemd5, data.Fileserverpath, data.ServerIp, data.ClientIp, time.Now().String())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Queryfile(filename string) ([]Fileinfo, error) {
	rwlock.RLock()
	defer rwlock.RUnlock()
	rows, err := Dbsqlite.Query("select * from fileinfo where filename=?", filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var result []Fileinfo

	for rows.Next() {
		var rowresult Fileinfo
		err = rows.Scan(&rowresult.Filename, &rowresult.Filemd5, &rowresult.Fileserverpath, &rowresult.ServerIp, &rowresult.ClientIp, &rowresult.Date)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, rowresult)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func Deletefile(filename string) error {
	rwlock.Lock()
	defer rwlock.Unlock()
	_, err := Dbsqlite.Exec("delete from fileinfo where filename=?", filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func Updatefile(data *Fileinfo) error {
	rwlock.Lock()
	defer rwlock.Unlock()
	_, err := Dbsqlite.Exec("update  fileinfo  set filemd5=?,fileserverpath=?,serverIp=?, clientIp=?,date=? where filename=?",
		data.Filemd5, data.Fileserverpath, data.ServerIp, data.ClientIp, data.Filename, time.Now().String())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
