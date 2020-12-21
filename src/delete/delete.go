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

	_, err := Db.Exec("delete from user_info where user_id=?", 5)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}

	fmt.Println("delete succ")
}
