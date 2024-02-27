package main

import (
	simplerouter "GO_project/router"
	"net/http"
)

func test(response simplerouter.HttpResponse, r *http.Request) {
	response.ReturnJson(200, "Success")
}

func err_custom(response simplerouter.HttpResponse, r *http.Request) {
	response.ReturnJson(404, "ご指定のページまたはファイルが見つかりませんでした")
}

func main() {
	router := simplerouter.New()
	router.SetMethodError(err_custom)
	router.POST("/test", test)
	http.ListenAndServe(":8080", router)
}
