package main

// 封装自己的web server
import (
	"fmt"
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
	Email             string `json:"email"` // ``中的Tag 在运行时使用reflect
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed-password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}

// SignUp 用户注册
func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}

	// TODO: Context应该在框架层面来创建
	ctx := Context{
		W: w,
		R: r,
	}

	err := ctx.ReadJson(req)
	if err != nil {
		fmt.Fprintf(w, "read body data failed, err: %v", err)
		return
	}

	resp := &commonResponse{
		Data: "123",
	}

	err = ctx.WriteJson(http.StatusOK, resp)
	if err != nil {
		// 如果响应失败，不应再次使用 ResponseWriter 往 response 里写数据
		fmt.Printf("写入响应失败, err: %v", err)
	}
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
