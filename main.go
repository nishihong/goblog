package main

import (
	"ch35/goblog/bootstrap"
	"ch35/goblog/pkg/database"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql" //匿名导入
)

//.StrictSlash(true) 去掉最后一个斜杠的问题  把会POST请求编程GET请求
//router := mux.NewRouter().StrictSlash(true)
var router *mux.Router
//设置了变量 结构体是 database/sql 包封装的一个数据库操作对象，包含了操作数据库的基本方法，通常情况下，我们把它理解为 连接池对象。
var db *sql.DB

//func init() {
//	sql.Register("mysql", &MySQLDriver{})
//}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func forceHTMLMiddleware(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2.继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 路由解析之前 就将后面的 / 去掉     直接访问 /会报错
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 除首页以外，移除所有请求路径后面的斜杆
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		// 2. 将请求传递下去
		next.ServeHTTP(w, r)
	})
}

func main() {
	database.Initialize()
	db = database.DB

	router = bootstrap.SetupRoute()

	// 中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL)

	//http.ListenAndServe(":3000", router)
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}