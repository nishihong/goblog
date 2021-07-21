package main

import (
	"ch35/goblog/app/http/middlewares"
	"ch35/goblog/bootstrap"
	"net/http"
)

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}