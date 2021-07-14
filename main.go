package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql" //匿名导入
)

//.StrictSlash(true) 去掉最后一个斜杠的问题  把会POST请求编程GET请求
//router := mux.NewRouter().StrictSlash(true)
var router = mux.NewRouter()
//设置了变量 结构体是 database/sql 包封装的一个数据库操作对象，包含了操作数据库的基本方法，通常情况下，我们把它理解为 连接池对象。
var db *sql.DB

//func init() {
//	sql.Register("mysql", &MySQLDriver{})
//}

func initDB() {
	var err error
	// 设置数据库连接信息
	config := mysql.Config{
		User: "root",
		Passwd: "root",
		Addr:"127.0.0.1:3306",
		Net:"tcp",
		DBName:"goblog",
		AllowNativePasswords: true,
	}

	//fmt.Println(config.FormatDSN())

	// 准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	//设置最大连接数 设置连接池最大打开数据库连接数，<= 0 表示无限制，默认为 0。
	db.SetMaxIdleConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试链接，失败会报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		//解析错误，这里应该有错误处理
		fmt.Fprint(w, "请提供正确的数据！")
		return
	}

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	//} else if len(title) < 3 || len(title) > 40 {
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if title == "" {
		errors["body"] = "标题不能为空"
	//} else if len(title) < 3 || len(title) > 40 {
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	if len(errors) == 0 {
		fmt.Fprintf(w, "验证通过！<br>")
		fmt.Fprintf(w, "title 的值为：%v <br>", title)
		//fmt.Fprintf(w, "title 的长度为：%v <br>", len(title)) //长度 中文3个
		fmt.Fprintf(w, "title 的长度为：%v <br>", utf8.RuneCountInString(title)) //长度 中文3个
		fmt.Fprintf(w, "body 的值为：%v <br>", body)
		//fmt.Fprintf(w, "body 的长度为：%v <br>", len(body)) //长度 中文3个
		fmt.Fprintf(w, "body 的长度为：%v <br>", utf8.RuneCountInString (body))
	} else {
		storeURL, _ := router.Get("articles.store").URL()

		//用以给模板文件传输变量时使用。
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		//关于模板后缀名 .gohtml ，可以使用任意后缀名，这不会影响代码的运行。常见的 Go 模板后缀名有
		tmpl, err := template.ParseFiles("goblog/resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL	*url.URL
	Errors map[string]string
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

//func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprint(w, "创建博文表单")
//}
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("goblog/resources/views/articles/create.gohtml")

	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, data)
}

func main() {
	initDB()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	//在有正则匹配的情况下，使用 : 区分。第一部分是名称，第二部分是正则表达式
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

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