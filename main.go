package main

import (
	"ch35/goblog/app/http/middlewares"
	"ch35/goblog/bootstrap"
	"ch35/goblog/config"
	c "ch35/goblog/pkg/config"
	"net/http"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	//http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}