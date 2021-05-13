package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal("panic occured:", err)
		panic(err)
	}
}

type Person struct {
	Username string
	PassWord string
	Email    string
}

func searchHanler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析
	fmt.Println(r.Form) //输出到服务器端的打印信息,map类型   type Values map[string][]string
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	for _, v := range r.Form {
		fmt.Println("val:", strings.Join(v, ""))
		//查询数据
		for i := 0; i < len(v); i++ {
			search_query := "select user_name,user_pass,user_mail from user_info where user_id=" + v[i] + ";"
			fmt.Println(search_query)
			rows, err := db.Query(search_query)
			checkErr(err)
			for rows.Next() {
				var user_name string
				var user_pass string
				var user_mail string
				err = rows.Scan(&user_name, &user_pass, &user_mail)
				checkErr(err)
				//将这一步改成JSON字符串传输至前面页面
				person := Person{
					Username: user_name,
					PassWord: user_pass,
					Email:    user_mail,
				}
				b, err := json.Marshal(person)
				checkErr(err)
				fmt.Fprintf(w, string(b))
			}
		}
	}
}

func loginHanler(w http.ResponseWriter, r *http.Request) {

	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Gonote
	err := json.Unmarshal(body, &n)
	checkErr(err)
	fmt.Println(n.Email, n.PassWord)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("select user_mail,user_pass from user_info where user_mail=? and user_pass=?", n.Email, n.PassWord)
	checkErr(err)

	var data bool
	data = rows.Next()
	//数据库返回参数为空时进入
	if !data {
		fmt.Println("入参:", n.Email, n.PassWord, "数据库:")
		fmt.Fprintf(w, `{"code":0}`)
	}
	//数据库返回参数不为空时进入
	for data {
		var user_mail string
		var user_pass string
		err = rows.Scan(&user_mail, &user_pass)
		checkErr(err)
		//将这一步改成JSON字符串传输至前面页面
		person := Person{
			Email:    user_mail,
			PassWord: user_pass,
		}
		fmt.Println("入参:", n.Email, n.PassWord, "数据库", person.Email, person.PassWord)
		fmt.Fprintf(w, `{"code":1}`)
		data = rows.Next()
	}
}

type Gonote struct {
	Username string
	PassWord string
	Email    string
	Phone    string
}

type Message struct {
	Message string
}

func AddHanler(w http.ResponseWriter, r *http.Request) {
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Gonote
	err := json.Unmarshal(body, &n)
	checkErr(err)
	fmt.Println(n.Username, n.PassWord, n.Email, n.Phone)

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
	if !data {
		fmt.Println("入参:", n.Email, "数据库:")
		fmt.Fprintf(w, `{"code":1}`)

		update, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
		checkErr(err)
		defer db.Close()
		//插入数据
		stmt, err := update.Prepare("INSERT user_info SET user_name=?,user_pass=?,user_mail=?,user_phone=?")
		checkErr(err)
		res, err := stmt.Exec(n.Username, n.PassWord, n.Email, n.Phone)
		checkErr(err)
		num, err := res.RowsAffected()
		checkErr(err)
		fmt.Println("受影响行数:", num)

	}
	//数据库返回参数不为空时进入
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

type NoteModify struct {
	UserId   int
	Username string
}

func ModifyHanler(w http.ResponseWriter, r *http.Request) {
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var nm NoteModify
	err := json.Unmarshal(body, &nm)
	checkErr(err)
	fmt.Println(nm.UserId, nm.Username)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()

	//更新数据
	stmt, err := db.Prepare("UPDATE user_info SET user_name=? where user_id=?")
	checkErr(err)
	res, err := stmt.Exec(nm.Username, nm.UserId)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("Affect rows :", affect)
	m := Message{
		Message: "Modify Success!",
	}
	b, err := json.Marshal(m)
	checkErr(err)
	fmt.Fprintf(w, string(b))
	//fmt.Fprintln(w,"modify success!")
}

func DeleteHanler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析
	fmt.Println(r.Form) //输出到服务器端的打印信息,map类型   type Values map[string][]string
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	for _, v := range r.Form {
		//fmt.Println("val:",strings.Join(v,""))
		//删除数据
		for i := 0; i < len(v); i++ {
			delete_query := "delete  from user_info where user_id=?"
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

type Tagore struct {
	Title       string
	Author      string
	Description string
	ImageUrl    string
}

type Tagoreslic struct {
	Tagores []Tagore
}

func TagoreHanler(w http.ResponseWriter, r *http.Request) {
	var t Tagoreslic
	r.ParseForm() //解析
	fmt.Println(r.Form)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("select title,author,description,imageUrl from Tagore;")
	checkErr(err)
	for rows.Next() {
		var title string
		var author string
		var description string
		var imageUrl string
		err = rows.Scan(&title, &author, &description, &imageUrl)
		checkErr(err)
		//将这一步改成JSON字符串传输至前面页面

		t.Tagores = append(t.Tagores, Tagore{Title: title, Author: author, Description: description, ImageUrl: imageUrl})
		// fmt.Println(t.Tagores)
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}

func main() {
	//mux :=http.NewServeMux()
	http.HandleFunc("/search", searchHanler)
	http.HandleFunc("/login", loginHanler)
	http.HandleFunc("/Add", AddHanler)
	http.HandleFunc("/Modify", ModifyHanler)
	http.HandleFunc("/Delete", DeleteHanler)
	http.HandleFunc("/Tagore", TagoreHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
