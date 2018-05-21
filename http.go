package main

import (
	mysql "FileServer/sqlite"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Deafult(res http.ResponseWriter, req *http.Request) {
	log.Println("Deafult", req.URL.String())
	// log.Println("head", req.Header)
	if req.URL.Path == "" || req.URL.Path == "/" {
		Index(res, req)
	} else {
		ViewFile(res, req)
	}

	return
}

//为了减少内存copy，使用指针，需要特别注意哦
func AuthAccount(res http.ResponseWriter, req *http.Request) error {

	auth := req.Header.Get("Authorization")
	if len(auth) == 0 {
		res.Header().Set("WWW-Authenticate", `Basic realm="No User Login"`)
		res.WriteHeader(http.StatusUnauthorized)
		log.Println("no auth info")
		return errors.New("no auth info ")
	}
	auth = strings.TrimPrefix(auth, "Basic")
	auth = strings.TrimSpace(auth)
	auth_data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		log.Println("decode failed", auth)
		return errors.New("DecodeString failed")
	}
	name := bytes.Split(auth_data, []byte(":"))[0]
	passwd := bytes.Split(auth_data, []byte(":"))[1]

	dbpasswd := mysql.QueryUserInfo(string(name))

	if len(dbpasswd) == 0 || bytes.Compare(passwd, dbpasswd) != 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return errors.New("auth failed")
	}
	return nil
}

func Index(res http.ResponseWriter, req *http.Request) {
	Login(res, req)
	return
	if AuthAccount(res, req) != nil {
		log.Println("AuthAccount failed")
		return
	}
	// paswd := mysql.QueryUserInfo(auth)

	tmp, err := template.ParseFiles("./myfileserver/htmlfile/index.html")
	if err != nil {
		res.Write([]byte("get index html error"))
		http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
	tmp.Execute(res, nil)

}

func ViewFile(res http.ResponseWriter, req *http.Request) {
	if AuthAccount(res, req) != nil {
		log.Println("AuthAccount failed")
		return
	}
	// defer req.Body.Close()
	log.Println("view")
	if strings.HasPrefix(req.URL.String(), "/viewfile") {
		// 这块用于展示
		// http.StripPrefix("/viewfile", http.FileServer(http.Dir("./myfileserver/upload/"))).ServeHTTP(res, req)
		http.StripPrefix("/viewfile", http.FileServer(http.Dir("./myfileserver/upload/"))).ServeHTTP(res, req)

	} else {
		res.Header().Set("Content-Type", "application/octet-stream") //设置文件下载类型
		http.FileServer(http.Dir("./myfileserver/upload/")).ServeHTTP(res, req)
		log.Println(req.Header.Get("Content-Type"), "|", res.Header().Get("Content-Type"))
	}
	return
}

type FileOpStat struct {
	Filename     string
	Filepath     string
	Uploadstatus string
	HttpStatus   int
	Description  string
}

//靠齐restful接口
//文件的查增删
func UploadFile(res http.ResponseWriter, req *http.Request) {
	if AuthAccount(res, req) != nil {
		log.Println("AuthAccount failed")
		return
	}
	defer req.Body.Close()
	switch req.Method {
	//GET
	case http.MethodGet:
		log.Println("get")
		t, err := template.ParseFiles("./myfileserver/htmlfile/file.html")

		if err != nil {
			res.Write([]byte("get upload html error"))
			http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
			return
		}
		t.Execute(res, "上传文件")
	//POST
	/*restful接口改造
	{
		filename:""
		uploadstatus:""
		httpstatus:""
		filepath:""
		Description:""
	}
	*/
	case http.MethodPost:
		log.Println("post")
		req.ParseMultipartForm(1 << 10) //设置下内存缓存大小，默认为1G
		log.Println(req)

		f, h, err := req.FormFile("uploadfile") //该值有html中定义
		// defer f.Close()
		if err != nil || f == nil || h == nil {
			// if req.Form
			fmt.Fprintln(res, "uploadfile", func(i *multipart.FileHeader) string {
				if i == nil {
					return "N/A"
				} else {
					return i.Filename
				}
			}(h), " failed")
			return
		}

		//查询
		queryres, err := mysql.Queryfile(h.Filename)
		if queryres != nil {
			for _, v := range queryres {
				file_res := FileOpStat{
					Filename:     h.Filename,
					Filepath:     v.Fileserverpath,
					Uploadstatus: "upload failed",
					HttpStatus:   http.StatusOK,
					Description:  fmt.Sprintf("%s%d bytes", h.Filename+"已经存在,", "上传时间:", v.Date),
				}
				// resdata, _ := json.Marshal(file_res)
				resdata, _ := json.MarshalIndent(file_res, "", "\t")
				res.Write(resdata)
			}

			return
		}

		file, err := os.OpenFile("./myfileserver/upload/"+h.Filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		defer file.Close()
		if err != nil {
			fmt.Fprintln(res, "create ", h.Filename, " in server failed")
			http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(file, f)
		if err != nil {
			fmt.Fprintln(res, "write ", h.Filename, " in server failed")
			http.Error(res, "StatusInternalServerError", http.StatusInternalServerError)
			return
		}
		filedir, _ := filepath.Abs("./myfileserver/upload/" + h.Filename)
		filesize, _ := file.Seek(0, 1)
		file_res := FileOpStat{
			Filename:     h.Filename,
			Filepath:     filedir,
			Uploadstatus: "upload success",
			HttpStatus:   http.StatusOK,
			Description:  fmt.Sprintf("%s%d bytes", h.Filename+"上传完成，文件大小:", filesize),
		}
		// resdata, _ := json.Marshal(file_res)
		resdata, _ := json.MarshalIndent(file_res, "", "\t")

		res.Write(resdata)

		//本地有了，那么来存数据库吧
		data := new(mysql.Fileinfo)
		data.Filename = h.Filename
		data.Fileserverpath = filedir
		md5val := md5.New()
		io.Copy(md5val, file)
		data.Filemd5 = fmt.Sprintf("%x", md5val.Sum(nil))
		log.Println("md5:", data.Filemd5)
		data.ServerIp = localconf.Ip
		data.ClientIp = req.RemoteAddr
		err = mysql.Addfile(data)
		if err != nil {
			log.Println(err)
		}

	//DELETE
	case http.MethodDelete:
		//获取文件名
		filename := req.URL.Path
		filename = "./myfileserver/upload" + strings.TrimPrefix(filename, "/uploadfile")

		log.Println("detele", filename)
		_, err := os.Stat(filename)
		if err != nil || os.IsNotExist(err) {

			file_res := FileOpStat{
				Filename:     filename,
				Uploadstatus: "delete success",
				HttpStatus:   http.StatusOK,
				Description:  "file not exist in sever",
			}
			// resdata, _ := json.Marshal(file_res)
			resdata, _ := json.MarshalIndent(file_res, "", "\t")
			res.Write(resdata)
			return

		}

		err = os.Remove(filename)
		if err != nil {
			file_res := FileOpStat{
				Filename:     filename,
				Uploadstatus: "delete failed",
				HttpStatus:   http.StatusOK,
				Description:  "sever delete failed",
			}
			// resdata, _ := json.Marshal(file_res)
			resdata, _ := json.MarshalIndent(file_res, "", "\t")
			res.Write(resdata)

			//此处应该起一个goroutine标记该文件是否还可用

			return
		}

		filedir, _ := filepath.Abs(filename)

		file_res := FileOpStat{
			Filename:     filedir,
			Uploadstatus: "delete success",
			HttpStatus:   http.StatusOK,
		}
		// resdata, _ := json.Marshal(file_res)
		resdata, _ := json.MarshalIndent(file_res, "", "\t")

		res.Write(resdata)
		mysql.Deletefile(filename)
	default:
		fmt.Fprintln(res, "only support get and post")

	}

	return
}

func myhttpmain(localconf *config_parm) {

	if localconf == nil {
		log.Println("localconf  is nil")
		return
	}
	fmt.Println("currentdir=", os.Args[0])
	Mux := http.NewServeMux()
	//根目录是主页
	Mux.HandleFunc("/", Deafult)
	Mux.HandleFunc("/viewfile", ViewFile)
	Mux.HandleFunc("/viewfile/", ViewFile)
	Mux.HandleFunc("/uploadfile", UploadFile)
	Mux.HandleFunc("/uploadfile/", UploadFile)

	server := http.Server{
		Addr:         localconf.Ip + ":" + localconf.Port,
		Handler:      Mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
