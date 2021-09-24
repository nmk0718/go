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
type UserRegister struct {
	Username string
	PassWord string
	Email    string
}
func UserRegisterHanler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n UserRegister
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
	
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	//获取请求报文的内容长度
	len := r.ContentLength
	if(len < 1){
		fmt.Println(len)
		fmt.Println("当前为空")
	}
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Login
	err := json.Unmarshal(body, &n)
	checkErr(err)
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
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
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
	if num != 0 {
		fmt.Fprintf(w, `{"code":1}`)
	} else {
		fmt.Fprintf(w, `{"code":0}`)
	}

}

//------------------用户删除Memo接口-----------------------------

type Message struct {
	Message string
}
func DeleteMemoHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	r.ParseForm()       //解析
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	for _, v := range r.Form {
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
				fmt.Fprintf(w, `{"code":1}`)
			} else {
				fmt.Fprintf(w, `{"code":0}`)
			}
		}
	}

}

//------------------用户修改Memo接口-----------------------------

type UpdateMemo struct {
	Id	  int
	Memo      string
	Time      string
}

func UpdateMemoHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n UpdateMemo
	err := json.Unmarshal(body, &n)
	checkErr(err)
	fmt.Println(n.Id, n.Memo,n.Time)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()

	//更新数据
	stmt, err := db.Prepare("UPDATE memo SET memo=?,time=? where id=?")
	checkErr(err)
	res, err := stmt.Exec(n.Memo, n.Time,n.Id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("Affect rows :", affect)
	if affect != 0 {
		fmt.Fprintf(w, `{"code":1}`)
	} else {
		fmt.Fprintf(w, `{"code":0}`)
	}
}


//------------------用户查询Memo接口-----------------------------

type SearchMemo struct {
	ID   int
	Memo string
	Time string
}
type SearchMemolic struct {
	SearchMemos []SearchMemo
}
func SearchMemoHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	var t SearchMemolic
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

		t.SearchMemos = append(t.SearchMemos, SearchMemo{ID: id, Memo: memo, Time: time})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}

//------------------获取banner图片-----------------------------


type Banner struct {
	ID   int
	Banner_Url string
}

type Bannerlic struct {
	Banners []Banner
}

func SearchBannerHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json

	var t Bannerlic
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	
	rows, err := db.Query("SELECT id,banner_url FROM banner")
	checkErr(err)

	for rows.Next() {
		var id 	 int
		var banner_url string
		err = rows.Scan(&id,&banner_url)
		checkErr(err)

		t.Banners = append(t.Banners, Banner{ID: id, Banner_Url: banner_url})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}

//------------------获取金刚区文字图标-----------------------------


type Icon struct {
	ID   int
	Icon_Text string
	Icon_Url string
}

type Iconlic struct {
	Jingang []Icon
}

func SearchIconsHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	
	var t Iconlic
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	
	rows, err := db.Query("SELECT id,icon_text,icon_url FROM Icons")
	checkErr(err)

	for rows.Next() {
		var id 	 int
		var icon_text string
		var icon_url string
		err = rows.Scan(&id,&icon_text,&icon_url)
		checkErr(err)

		t.Jingang = append(t.Jingang, Icon{ID: id, Icon_Text: icon_text,Icon_Url:icon_url})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}

//------------------购物车订单列表-----------------------------

type SearchShopping struct {
	Store_ID   int
	Store_Name string
	BaoYou string
	ProductList []Product
}

type Product struct {
	ID   int
	Product_Name string
	Color string
	Number int
	Price int
	Selected string
}

type SearchShoppinglic struct {
	SearchShoppings []SearchShopping
}
func SearchShoppingCartHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	var t SearchShoppinglic
	var k  SearchShopping
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
	
	rows, err := db.Query("SELECT C.store_id,C.store_name,C.baoyou FROM shoppingcart as A INNER JOIN shopping AS B ON A.shopping_id=B.id INNER JOIN store as C ON B.store_id =C.store_id WHERE A.user_id=? GROUP BY B.store_id;",n.User_id)
	checkErr(err)

	for rows.Next() {
		var store_id 	 int
		var store_name string
		var baoyou string
		err = rows.Scan(&store_id,&store_name,&baoyou)
		checkErr(err)

		k.ProductList = nil
		products, abc := db.Query("SELECT B.id,B.`name`,B.color,A.number,B.price,B.selected FROM shoppingcart as A INNER JOIN shopping AS B ON A.shopping_id=B.id WHERE A.user_id=? AND B.store_id=?;",n.User_id,store_id)
		checkErr(abc)
		for products.Next() {
			var id 	 int
			var product_name string
			var color string
			var number int
			var price int
			var selected string
			abc = products.Scan(&id,&product_name,&color,&number,&price,&selected)
			checkErr(abc)
			k.ProductList = append(k.ProductList, Product{ID: id, Product_Name: product_name,Color:color,Number: number,Price : price,Selected:selected})
		}
		t.SearchShoppings = append(t.SearchShoppings, SearchShopping{Store_ID: store_id, Store_Name: store_name,BaoYou:baoyou,ProductList:k.ProductList})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}


//------------------修改购物车订单数量-----------------------------
type ShoppingCart struct{
	User_id int
	Shopping_id int
	Shopping_number int
}

func UpdateShoppingCartNumberHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n ShoppingCart
	err := json.Unmarshal(body, &n)
	checkErr(err)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()

	//更新数据
	stmt, err := db.Prepare("UPDATE  shoppingcart SET number=? where user_id=? and shopping_id =?")
	checkErr(err)
	res, err := stmt.Exec(n.Shopping_number,n.User_id,n.Shopping_id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("Affect rows :", affect)
	if affect != 0 {
		fmt.Fprintf(w, `{"code":1}`)
	} else {
		fmt.Fprintf(w, `{"code":0}`)
	}
}

//------------------删除购物车订单-----------------------------

func DeleteShoppingCartOrderHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	r.ParseForm()       //解析
	fmt.Println(r.Form) 
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	for _, v := range r.Form {
		//删除数据
		for i := 0; i < len(v); i++ {
			delete_query := "delete from shoppingcart where shopping_id=?"
			stmt, err := db.Prepare(delete_query)
			checkErr(err)

			res, err := stmt.Exec(v[i])
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("affect rows:", affect)
			if affect != 0 {
				fmt.Fprintf(w, `{"code":1}`)
			} else {
				fmt.Fprintf(w, `{"code":0}`)
			}
		}
	}

}


//------------------分页查询商品列表-----------------------------

type Shopping_Info struct {
	TotalPageNum int
	ProductinfoList []Productinfo
}

type Productinfo struct {
	ID   int
	Store_Name string
	Store_Icon string
	Product_Name string
	ImageUrl string
	Color string
	Price int
	Sales_Volume int
}

type Searchshoppinginfolic struct {
	Shoppinginfos []Shopping_Info
}
type PagedQuery struct{
	PageNo int
	PageSize int
	Input_text string
	Status int
}
func PagedQueryShoppingHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	var t Searchshoppinginfolic
	var k  Shopping_Info
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n PagedQuery
	var paixu string
	err := json.Unmarshal(body, &n)
	checkErr(err)
	var page_no = (n.PageNo-1)*n.PageSize
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT COUNT(*) FROM shopping where name like CONCAT('%',?,'%');",n.Input_text)
	checkErr(err)

	for rows.Next() {
		var totalrecord  int
		err = rows.Scan(&totalrecord)
		checkErr(err)

		if(n.Status == 1 ){
			paixu = "id ASC"
		}else if (n.Status == 2){
			paixu = "sales_volume ASC"
		}else if (n.Status ==3){
			paixu = "sales_volume DESC"
		}else if (n.Status ==4){
			paixu = "price ASC"
		}else if (n.Status ==5){
			paixu = "price DESC"
		}

		search_query := "SELECT p.id,s.store_name,s.store_icon,p.`name`,p.imageurl,p.color,p.price,p.sales_volume FROM shopping p,store s where p.store_id = s.store_id and  p.name like CONCAT('%',?,'%') order by "+ paixu +" LIMIT ?,?;"
		productinfos, abc := db.Query(search_query,n.Input_text,page_no,n.PageSize)
		checkErr(abc)
		for productinfos.Next() {
			var id 	 int
			var store_name string
			var store_icon string
			var product_name string
			var imageurl string
			var color string
			var price int
			var sales_volume int 
			abc = productinfos.Scan(&id,&store_name,&store_icon,&product_name,&imageurl,&color,&price,&sales_volume)
			checkErr(abc)
			k.ProductinfoList = append(k.ProductinfoList, Productinfo{ID: id,Store_Name:store_name,Store_Icon:store_icon,Product_Name: product_name,ImageUrl:imageurl,Color:color,Price : price,Sales_Volume:sales_volume})
		}
		var totalnumber = (totalrecord + n.PageSize -1 )/n.PageSize
		t.Shoppinginfos = append(t.Shoppinginfos, Shopping_Info{TotalPageNum: totalnumber,ProductinfoList:k.ProductinfoList})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}



//------------------商品详情-----------------------------

type PictureDetails struct {
	Picture string
	Category string
}
type ProductDetailslic struct {
	ProductDetails []PictureDetails
}
type Shopping_ID struct{
	Shopping_ID int
}
func ProductDetailsHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	var t ProductDetailslic
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var n Shopping_ID
	err := json.Unmarshal(body, &n)
	checkErr(err)

	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()
	rows, err := db.Query("SELECT picture,category from ProductDetails WHERE shopping_id = ?",n.Shopping_ID)
	checkErr(err)

	for rows.Next() {
		var picture  string
		var category  string
		err = rows.Scan(&picture,&category)
		checkErr(err)

		t.ProductDetails = append(t.ProductDetails, PictureDetails{Picture: picture,Category:category})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}


//------------------发布推文列表-----------------------------

type Tweetinfo struct {
	User_id int
	Tweet    string
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
	insert, err := db.Prepare("insert into tweet set user_id =?,tweet =?,time=?")
	checkErr(err)
	ress, err := insert.Exec(n.User_id,n.Tweet,n.Time)
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
	Retweets int
	Likes int
	Retweet_status int
	Like_status int
	Comments [] string
	Time string
	Imagelist [] string
}

type Searchlic struct {
	Searchs []Search
}
type Requestinfo struct{
	User_id int
	PageNo int
	PageSize int
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
	var n Requestinfo
	err := json.Unmarshal(body, &n)
	checkErr(err)
	//定义分页
	var page_no = (n.PageNo-1)*n.PageSize
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT A.id,A.tweet,B.user_name,B.headimage,B.user_code,A.retweets,A.likes,ifnull(C.retweet_status,0) as retweet_status,ifnull(C.like_status,0) as like_status,A.time FROM tweet AS A LEFT JOIN user_info AS B ON A.user_id = B.user_id LEFT JOIN tweet_status AS C ON C.tweet_id = A.id ORDER BY A.time DESC LIMIT ?,?",page_no,n.PageSize)
	checkErr(err)

	for rows.Next() {
		var id 	 int
		var tweet string
		var user_name string
		var headimage string
		var user_code string
		var retweets int
		var likes int
		var retweet_status int
		var like_status int
		var time string
			err = rows.Scan(&id,&tweet,&user_name,&headimage,&user_code,&retweets,&likes,&retweet_status,&like_status,&time)
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

		k.Comments = nil
		products, comment := db.Query("select comments from comments where tweet_id = ?;",id)
		checkErr(comment)
		for products.Next() {
			var comments 	 string
			comment = products.Scan(&comments)
			checkErr(comment)
			k.Comments = append(k.Comments,comments)
		}
		t.Searchs = append(t.Searchs, Search{ID:id,Tweetinfo: tweet, User_name: user_name,Headimage:headimage,Retweets:retweets,Likes:likes,Retweet_status:retweet_status,Like_status:like_status,Comments:k.Comments,Time:time,User_code:user_code,Imagelist:k.Imagelist})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}


//------------------转推和点赞推文-----------------------------
type ClickInfo struct{
	User_id int
	tweet_id int
	Clicktype int  //0为转发 1为点赞
	Click int     //0为取消点击 1为点击
}

func UpdateTweetStatusHanler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE") //允许的请求类型
	w.Header().Set("content-type", "application/json;charset=UTF-8") //返回数据格式是json
	//获取请求报文的内容长度
	len := r.ContentLength
	//新建一个字节切片，长度与请求报文的内容长度相同
	body := make([]byte, len)
	//读取r的请求主体，并将具体内容读入body中
	r.Body.Read(body)
	var c ClickInfo
	err := json.Unmarshal(body, &c)
	checkErr(err)
	//打开数据库操作
	db, err := sql.Open("mysql", "root:zy@2021@tcp(121.4.147.189:3306)/flutter_app")
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `tweet_status` where  user_id =? and tweet_id = ? ",c.User_id,c.tweet_id)
	checkErr(err)

	var data bool
	data = rows.Next()
	//数据库返回参数为空时进入
	if !data {
		if(c.Clicktype == 0){
			insertdata, err := db.Prepare("INSERT tweet_status SET retweet_status=?,user_id=?,tweet_id=?")
			checkErr(err)
			res, err := insertdata.Exec(c.Click,c.User_id,c.tweet_id)
			checkErr(err)
			num, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", num)
			if(c.Click == 0 ){
				updatedata, err := db.Prepare("UPDATE tweet SET retweets=(retweets-1) where id=?")
				checkErr(err)
				res, err := updatedata.Exec(c.tweet_id)
				checkErr(err)
				affect, err := res.RowsAffected()
				checkErr(err)
				fmt.Println("受影响行数:", affect)
			}else {
				updatedata, err := db.Prepare("UPDATE tweet SET retweets=(retweets+1) where id=?")
				checkErr(err)
				res, err := updatedata.Exec(c.tweet_id)
				checkErr(err)
				affect, err := res.RowsAffected()
				checkErr(err)
				fmt.Println("受影响行数:", affect)
			}
		}else {
			insertdata, err := db.Prepare("INSERT tweet_status SET like_status=?,user_id=?,tweet_id=?")
			checkErr(err)
			res, err := insertdata.Exec(c.Click,c.User_id,c.tweet_id)
			checkErr(err)
			num, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", num)
			if(c.Click == 0 ){
				updatedata, err := db.Prepare("UPDATE tweet SET likes=(likes-1) where id=?")
				checkErr(err)
				res, err := updatedata.Exec(c.tweet_id)
				checkErr(err)
				affect, err := res.RowsAffected()
				checkErr(err)
				fmt.Println("受影响行数:", affect)
			}else {
				updatedata, err := db.Prepare("UPDATE tweet SET likes=(likes+1) where id=?")
				checkErr(err)
				res, err := updatedata.Exec(c.tweet_id)
				checkErr(err)
				affect, err := res.RowsAffected()
				checkErr(err)
				fmt.Println("受影响行数:", affect)
			}
		}

	//数据库返回参数不为空时进入
	}else if data {  

	if(c.Clicktype == 0){
		updatedata, err := db.Prepare("UPDATE tweet_status SET retweet_status=? where user_id =? and tweet_id=?")
		checkErr(err)
		res, err := updatedata.Exec(c.Click,c.User_id,c.tweet_id)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println("受影响行数:", affect)
		if(c.Click == 0 ){
			updatedata, err := db.Prepare("UPDATE tweet SET retweets=(retweets-1) where id=?")
			checkErr(err)
			res, err := updatedata.Exec(c.tweet_id)
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", affect)
		}else {
			updatedata, err := db.Prepare("UPDATE tweet SET retweets=(retweets+1) where id=?")
			checkErr(err)
			res, err := updatedata.Exec(c.tweet_id)
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", affect)
		}
	}else{
		updatedata, err := db.Prepare("UPDATE tweet_status SET like_status=? where user_id =? and tweet_id=?")
		checkErr(err)
		res, err := updatedata.Exec(c.Click,c.User_id,c.tweet_id)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println("受影响行数:", affect)
		if(c.Click == 0 ){
			updatedata, err := db.Prepare("UPDATE tweet SET likes=(likes-1) where id=?")
			checkErr(err)
			res, err := updatedata.Exec(c.tweet_id)
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", affect)
		}else {
			updatedata, err := db.Prepare("UPDATE tweet SET likes=(likes+1) where id=?")
			checkErr(err)
			res, err := updatedata.Exec(c.tweet_id)
			checkErr(err)
			affect, err := res.RowsAffected()
			checkErr(err)
			fmt.Println("受影响行数:", affect)
		}
	}

	}
}

//------------------接口配置-----------------------------

func main() {
	http.HandleFunc("/UserRegister", UserRegisterHanler)
	http.HandleFunc("/UserLogin", UserLoginHanler)
	http.HandleFunc("/AddMemo", AddMemoHanler)
	http.HandleFunc("/DeleteMemo", DeleteMemoHanler)
	http.HandleFunc("/UpdateMemo", UpdateMemoHanler)
	http.HandleFunc("/SearchMemo", SearchMemoHanler)
	http.HandleFunc("/SearchBanner", SearchBannerHanler)
	http.HandleFunc("/SearchIcons", SearchIconsHanler)
	http.HandleFunc("/SearchShoppingCart", SearchShoppingCartHanler)
	http.HandleFunc("/UpdateShoppingCartNumber", UpdateShoppingCartNumberHanler)
	http.HandleFunc("/DeleteShoppingCartOrder", DeleteShoppingCartOrderHanler)
	http.HandleFunc("/PagedQueryShopping", PagedQueryShoppingHanler)
	http.HandleFunc("/ProductDetails", ProductDetailsHanler)
	http.HandleFunc("/AddTweet", AddTweetHanler)
	http.HandleFunc("/SearchTweet", SearchTweetHanler)
	http.HandleFunc("/UpdateTweetStatus", UpdateTweetStatusHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
