package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "strings"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal("panic occured:", err)
		panic(err)
	}
}


//------------------发布推文列表-----------------------------

type Tweetinfo struct {
	User_id int
	Memo    string
	Image_url []string
	Time    string
}

func AddTweetHanler(w http.ResponseWriter, r *http.Request) {
	//获取请求报文的内容长度
	lens := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, lens)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Tweetinfo
	err := json.Unmarshal(body, &n)
	checkErr(err)

	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	insert, err := db.Prepare("insert into memo set user_id =?,memo =?,time=?")
	checkErr(err)
	ress, err := insert.Exec(n.User_id,n.Memo,n.Time)
	checkErr(err)

	id, err := ress.LastInsertId()

	fmt.Println("insert succ:", id)

	for i := 0 ;i< len(n.Image_url);i++ {
	stmt, err := db.Prepare("insert into images set image_url=?,publish_id=? ")
	checkErr(err)
	res, err := stmt.Exec(n.Image_url[i], id)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("受影响行数:", num)
	}
}

//------------------查询推文列表-----------------------------

type Search struct {
	ID   int
	Tweetinfo string
	User_name string
	Headimage string
	User_code string
	Imagelist [] string
}

type Searchlic struct {
	Searchs []Search
}
type UserID struct{
	User_id int
}
func SearchTweetHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	var t Searchlic
	var k  Search
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n UserID
	err := json.Unmarshal(body, &n)
	checkErr(err)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	
	rows, err := db.Query("SELECT memo.id,memo.memo,user_info.user_name,user_info.headimage,user_info.user_code FROM memo,user_info where memo.user_id=user_info.user_id;")
	checkErr(err)

	for rows.Next() {
		var id 	 int
		var memo string
		var user_name string
		var headimage string
		var user_code string
		err = rows.Scan(&id,&memo,&user_name,&headimage,&user_code)
		checkErr(err)

		k.Imagelist = nil
		products, abc := db.Query("select image_url from images where publish_id = ?;",id)
		checkErr(abc)
		for products.Next() {
			var image_url 	 string
			abc = products.Scan(&image_url)
			checkErr(abc)
			k.Imagelist = append(k.Imagelist,image_url)
		}
		t.Searchs = append(t.Searchs, Search{Tweetinfo: memo, User_name: user_name,Headimage:headimage,User_code:user_code,Imagelist:k.Imagelist})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}


//------------------接口配置-----------------------------

func main() {
	http.HandleFunc("/AddTweet", AddTweetHanler)
	http.HandleFunc("/SearchTweet", SearchTweetHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
