package main

import (
	"fmt"
	"net/http"

	"common/log"
	"common/url"
	"common/url/test/handler"
)

func NewLogger() (*log.Logger, error) {
	var h log.Handler
	var err error
	h, err = log.NewRotatingFileHandler("test.log", 200000000, 20)
	if err != nil {
		fmt.Printf("new log handler err: %v\n", err.Error())
		return nil, err
	}
	logger := log.NewDefault(h)
	logger.SetLevel(0)
	return logger, nil
}
func main() {
	logger, err := NewLogger()
	if err != nil {
		fmt.Printf("the err is %s:", err.Error())
		return
	}
	r := &url.RestUrlRouter{
		Logger: logger,
	}
	r.Init("/{version}/regions/{region}/")
	// snapshot
	snHandler := &handler.SnapshotHandler{}
	snHandler.Register(r)

	http.Handle("/", r)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("[Init] http server exit, error: %s", err.Error())
	}
}

//func sayhelloName(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()       //解析参数，默认是不会解析的
//	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
//	fmt.Println("path", r.URL.Path)
//	fmt.Println("scheme", r.URL.Scheme)
//	fmt.Println(r.Form["url_long"])
//	for k, v := range r.Form {
//		fmt.Println("key:", k)
//		fmt.Println("val:", strings.Join(v, ""))
//	}
//	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
//}

//func main() {
//	http.HandleFunc("/", sayhelloName)       //设置访问的路由
//	err := http.ListenAndServe(":9090", nil) //设置监听的端口
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//}
