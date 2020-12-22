package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func loginHanler(w http.ResponseWriter, r *http.Request) {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(r.Body)
	// 用于存放参数key=value数据
	var params map[string]string
	// 解析参数 存入map
	decoder.Decode(&params)
	fmt.Printf("POST json: username=%s, password=%s\n", params["username"], params["password"])
	fmt.Fprintf(w, `{"code":0}`)
}
func main() {
	http.HandleFunc("/login", loginHanler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("panic occur:", err)
	}

}
