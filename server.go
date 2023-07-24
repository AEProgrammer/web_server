package web_server

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保一定实现了Server接口
var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// AddRoute 增加路由注册功能 method是http方法 path是路径
	AddRoute(method string, path string, handleFunc HandleFunc)
}

type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		newRouter(),
	}
}

//func (h *HTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
//	// Context创建
//
//	// 路由匹配
//
//	// 执行业务逻辑
//}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.AddRoute(http.MethodPost, path, handleFunc)
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Context创建
	ctx := &Context{
		Req: request,
		Resp: writer,
	}
	// 路由匹配

	// 执行业务逻辑
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// 查找路由 执行业务逻辑
}

func (h *HTTPServer) Start(addr string) error {
	l ,err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(l, h)
}

