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

type Search struct {
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

type Searchlic struct {
	Searchs []Search
}
type UserID struct{
	User_id int
}
func SearchShoppingCartHanler(w http.ResponseWriter, r *http.Request) {
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
		t.Searchs = append(t.Searchs, Search{Store_ID: store_id, Store_Name: store_name,BaoYou:baoyou,ProductList:k.ProductList})
	}
	b, err := json.Marshal(t)
	checkErr(err)
	fmt.Fprintf(w, string(b))

}


//------------------修改购物车订单数量-----------------------------
type ShoppingID struct{
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
	var n ShoppingID
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


//------------------接口配置-----------------------------

func main() {
	http.HandleFunc("/SearchBanner", SearchBannerHanler)
	http.HandleFunc("/SearchIcons", SearchIconsHanler)
	http.HandleFunc("/SearchShoppingCart", SearchShoppingCartHanler)
	http.HandleFunc("/UpdateShoppingCartNumber", UpdateShoppingCartNumberHanler)
	http.HandleFunc("/DeleteShoppingCartOrder", DeleteShoppingCartOrderHanler)
	http.HandleFunc("/PagedQueryShopping", PagedQueryShoppingHanler)
	http.HandleFunc("/ProductDetails", ProductDetailsHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
