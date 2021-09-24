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
	tweet    string
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
	ress, err := insert.Exec(n.User_id,n.tweet,n.Time)
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
	http.HandleFunc("/AddTweet", AddTweetHanler)
	http.HandleFunc("/SearchTweet", SearchTweetHanler)
	http.HandleFunc("/UpdateTweetStatus", UpdateTweetStatusHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
