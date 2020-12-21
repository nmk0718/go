package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	UserId   int    `db:"user_id"`
	Username string `db:"user_name"`
	PassWord string `db:"user_pass"`
	Email    string `db:"user_mail"`
}

var Db *sqlx.DB

func init() {
	database, err := sqlx.Open("mysql", "hospitalTest:Liangjian123360@8899@tcp(192.168.50.57:3306)/flutter_app")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database
}

func main() {
	r, err := Db.Exec("insert into user_info(user_name, user_pass, user_mail)values(?, ?, ?)", "吴彦祖", "888888", "wuyanzu@qq.com")
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}

	fmt.Println("insert succ:", id)
}
