//一个标准的可执行的 Go 程序必须有 package main 的声明。
package main

//我们使用 import 关键词用以引入程序所需的 Go 包。在 goblog 中，我们引入了两个 Go 标准库的包。
import (

	//标准库 net/http 提供了 HTTP 编程有关的接口，内部封装了 TCP 连接和报文解析的复杂琐碎细节。http 提供了 HTTP 客户端和服务器实现。
	//HTTP 客户端可用以发送请求到第三方 API 或者请求网页，以获取所需数据，类似于 curl 或者 wget 。
	//HTTP 服务器用以提供 HTTP 服务器来处理 HTTP 请求，此处我们使用了此功能：
	"fmt"
	"net/http"
)

//http.ResponseWriter 是返回用户的响应，一般用 w 作为简写。
//返回 500 状态码 w.WriteHeader(http.StatusInternalServerError)
//设置返回标头 w.Header().Set("name", "my name is smallsoup")

//http.Request 是用户的请求信息，一般用 r 作为简写。
// r.URL.Query() 获取用户参数
//获取客户端信息 r.Header.Get("User-Agent")
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	//fmt 的 Fprint 函数会将内容输出到实现了 io.Writer 接口类型的变量 w 中，我们通常用这个函数往文件中写入内容。注意，只要满足 io.Writer 接口的类型都支持写入。在我们的代码中 w 是 http.ResponseWriter 的实例，已经实现了 io.Writer 接口。
	//fmt.Fprint(w, "<h1>Hello，这里是 goblog</h1>")

	//fmt.Fprint(w, "请求路径为："+r.URL.Path)

	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello,这里是 goblogdf</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "此微博是用以记录变成笔记，如您有反馈或建议，请联系 " + "<a href=\"mailto:summer@example.com\">summer@example.com</a>")
	} else {
		fmt.Fprint(w, "<h1>请求页面未找到 ：(</h1>" + "<p>如有疑惑，请联系我们。</p>")
	}
}

func main() {
	//http.HandleFunc 用以指定处理 HTTP 请求的函数，此函数允许我们只写一个 handler（在此例子中 handlerFunc，可任意命名），请求会通过参数传递进来，使用者只需与 http.Request 和 http.ResponseWriter 两个对象交互即可。
	//http.HandleFunc 里传参的 / 意味着 任意路径。
	http.HandleFunc("/", handlerFunc)
	//http.ListenAndServe 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口，如 baidu.com:80，另一个 Web 常用是 HTTPS 的 443 端口，如 baidu.com:443。当我们监听本地端口时，可使用 localhost 加上端口号来访问，如以下代码：
	http.ListenAndServe(":3000", nil)
}