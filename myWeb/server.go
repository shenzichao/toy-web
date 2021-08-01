package main

// 封装自己的web server
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Server 是http server 的顶级抽象
type Server interface {
	// Route 选择路由, 命中会执行handlerFunc
	Route(pattern string, handlerFunc http.HandlerFunc)
	// Start 启动 server
	Start(address string) error
}

// 基于 Go的net/http实现 http server
type sdkHttpServer struct {
	// Name server的名称
	Name string
}

func (s *sdkHttpServer) Route(pattern string, handlerFunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlerFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, nil)
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed-password"`
}

func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}

// SignUp 用户注册
func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}
	// TODO：47-57 的代码每次读取json输入都需要写一遍，branch verson2 引入context优化
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed, error: %v", err)
		return
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Fprintf(w, "deserialied failed, err: %v", err)
		return
	}

	// 返回虚拟id表示注册成功
	fmt.Fprintf(w, "%d", err)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home page \n")
}

func main() {
	server := NewSdkHttpServer("test_server")
	server.Route("/", homePage)
	server.Route("/signup", SignUp)
	server.Start(":8080")
}
