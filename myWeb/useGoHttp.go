package main

import (
	"fmt"
	"io"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home page \n")
	// 请求中的body 只能读取一次
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed, err: %v", err)
	}
	// body 的类型为 byte[]
	fmt.Fprintf(w, "read the data %s \n", string(body))

	// 再次读取body, 不会报错，但是内容为空
	body, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body one more time failed, err: %v", err)
	}
	fmt.Fprintf(w, "read body one more time, data is [%s], the length of data is %d \n", string(body), len(body))
}

// 单纯使用Go的http库实现web
func main1() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":8080", nil)

}
