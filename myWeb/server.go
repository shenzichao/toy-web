package main

// 封装自己的web server
import (
	"fmt"
	"net/http"
)

// Server 是http server 的顶级抽象
type Server interface {
	// Route 选择路由, 命中会执行handlerFunc
	// 支持 REST API的类型 POST PUT DELETE GET
	Route(method, pattern string, handlerFunc func(ctx *Context))
	// Start 启动 server
	Start(address string) error
}

// 基于 Go的net/http实现 http server
type sdkHttpServer struct {
	// Name server的名称
	Name    string
	handler *HandlerBasedOnMap
}

func (s *sdkHttpServer) Route(method, pattern string,
	handlerFunc func(ctx *Context)) {
	// TODO: 与HandlerBasedOnMap 强耦合 s.handler.handlers意味着需要知道s.handler的底层实现，需要解耦
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handlerFunc
}

func (s *sdkHttpServer) Start(address string) error {
	// 注册根路由的handler
	http.Handle("/", s.handler)
	return http.ListenAndServe(address, nil)
}

type signUpReq struct {
	Email             string `json:"email"` // ``中的Tag 在运行时使用reflect
	Password          string `json:"password"`
	ConformedPassword string `json:"conformed-password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
		handler: &HandlerBasedOnMap{
			handlers: make(map[string]func(ctx *Context)),
		},
	}
}

// SignUp 用户注册
func SignUp(ctx *Context) {
	req := &signUpReq{}

	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequestJson(err)
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
	//server.Route("/", homePage)
	server.Route("POST", "/signup", SignUp)
	server.Start(":8080")
}
