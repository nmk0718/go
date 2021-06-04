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



//------------------用户注册接口-----------------------------
type Gonote struct {
	Username string
	PassWord string
	Email    string
}
type Person struct {
	Email    string
}
func UserRegisterHanler(w http.ResponseWriter, r *http.Request) {
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Gonote
	err := json.Unmarshal(body, &n)
	checkErr(err)
	fmt.Println(n.Username, n.PassWord, n.Email)

	//查询邮箱是否重复
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("select user_mail from user_info where user_mail=?", n.Email)
	checkErr(err)

	var data bool
	data = rows.Next()
	//数据库返回参数为空时进入
	//邮箱未注册返回code为1
	if !data {
		fmt.Println("入参:", n.Email, "数据库:")
		fmt.Fprintf(w, `{"code":1}`)

		update, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
		checkErr(err)
		defer db.Close()
		//插入数据
		stmt, err := update.Prepare("INSERT user_info SET user_name=?,user_pass=?,user_mail=?")
		checkErr(err)
		res, err := stmt.Exec(n.Username, n.PassWord, n.Email)
		checkErr(err)
		num, err := res.RowsAffected()
		checkErr(err)
		fmt.Println("受影响行数:", num)

	}
	//数据库返回参数不为空时进入
	//邮箱已经注册,返回code为0
	for data {
		var user_mail string
		err = rows.Scan(&user_mail)
		checkErr(err)
		//将这一步改成JSON字符串传输至前面页面
		person := Person{
			Email: user_mail,
		}
		fmt.Println("入参:", n.Email, "数据库", person.Email)
		fmt.Fprintf(w, `{"code":0}`)
		data = rows.Next()
	}

}

//------------------用户登录接口-----------------------------
type Login struct {
	Email    string
	PassWord string
}
type UserID struct{
	User_id int
}
func UserLoginHanler(w http.ResponseWriter, r *http.Request) {

	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Login
	err := json.Unmarshal(body, &n)
	checkErr(err)
	fmt.Println(n.Email, n.PassWord)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("select user_id from user_info where user_mail=? and user_pass=?", n.Email, n.PassWord)
	checkErr(err)

	var data bool
	data = rows.Next()
	//数据库返回参数为空时进入
	if !data {
		fmt.Println("入参:", n.Email, n.PassWord, "数据库:")
		fmt.Fprintf(w, `{"user_id":0}`)
	}
	//数据库返回参数不为空时进入
	for data {
		var user_id int
		err = rows.Scan(&user_id)
		checkErr(err)
		//将这一步改成JSON字符串传输至前面页面
		b, err := json.Marshal(UserID{User_id: user_id})
		checkErr(err)
		fmt.Fprintf(w, string(b))
		data = rows.Next()
	}
}


//------------------用户新增Memo接口-----------------------------
type Memolic struct {
	Memos []Memo
}
type Memo struct {
	User_id	  int
	Memo      string
	Time      string
}
func AddMemoHanler(w http.ResponseWriter, r *http.Request) {
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Memo
	err := json.Unmarshal(body, &n)
	checkErr(err)

	update, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	//插入数据
	stmt, err := update.Prepare("INSERT memo SET user_id=?,memo=?,time=?")
	checkErr(err)
	res, err := stmt.Exec(n.User_id,n.Memo, n.Time)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("受影响行数:", num)
	fmt.Fprintf(w, `{"code":1}`)

}

//------------------用户查询Memo接口-----------------------------

type Search struct {
	ID   int
	Memo string
	Time string
}
type Searchlic struct {
	Searchs []Search
}
func SearchMemoHanler(w http.ResponseWriter, r *http.Request) {
	var t Searchlic
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
	rows, err := db.Query("select id,memo,time from memo where user_id=?",n.User_id)
	checkErr(err)
	for rows.Next() {
		var id 	 int
		var memo string
		var time string
		err = rows.Scan(&id,&memo, &time)
		checkErr(err)
		//将这一步改成JSON字符串传输至前面页面

		t.Searchs = append(t.Searchs, Search{ID: id, Memo: memo, Time: time})
	}
	b, err := json.Marshal(t)
	//{"Memos":[{"我","Time":"2021-04-15 18:23:39"},{"Memo":"你","Time":"2021-04-16 09:25:33"}]}
	//	b, err := json.Marshal(t.Memos)
	//[{"Memo":"我","Time":"2021-04-15 18:23:39"},{"Memo":"你","Time":"2021-04-16 09:25:33"}]
	checkErr(err)
	fmt.Fprintf(w, string(b))

}

//------------------用户删除Memo接口-----------------------------

type Message struct {
	Message string
}
func DeleteMemoHanler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析
	fmt.Println(r.Form) //输出到服务器端的打印信息,map类型   type Values map[string][]string
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	for _, v := range r.Form {
		fmt.Println("val:",strings.Join(v,""))
		//删除数据
		for i := 0; i < len(v); i++ {
			delete_query := "delete from memo where id=?"
			stmt, err := db.Prepare(delete_query)
			checkErr(err)

			res, err := stmt.Exec(v[i])
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("affect rows:", affect)
			if affect != 0 {
				m := Message{
					Message: "Delete Success!",
				}
				b, err := json.Marshal(m)
				checkErr(err)
				fmt.Fprintf(w, string(b))
				//fmt.Fprintln(w,"Delete success! ")
			} else {
				m := Message{
					Message: "The id not exist!",
				}
				b, err := json.Marshal(m)
				checkErr(err)
				fmt.Fprintf(w, string(b))
				//fmt.Fprintf(w,"The id not exist!")
			}
		}
	}

}

//------------------接口配置-----------------------------

func main() {
	http.HandleFunc("/UserRegister", UserRegisterHanler)
	http.HandleFunc("/UserLogin", UserLoginHanler)
	http.HandleFunc("/SearchMemo", SearchMemoHanler)
	http.HandleFunc("/AddMemo", AddMemoHanler)
	http.HandleFunc("/DeleteMemo", DeleteMemoHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
