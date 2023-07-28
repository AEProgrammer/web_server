package web_server

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保一定实现了Server接口
var _ Server = &hTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// AddRoute 增加路由注册功能 method是http方法 path是路径
	addRoute(method string, path string, handleFunc HandleFunc)
}

type hTTPServer struct {
	*router
}

func NewHTTPServer() *hTTPServer {
	return &hTTPServer{
		newRouter(),
	}
}

//func (h *hTTPServer) AddRoute(method string, path string, handleFunc HandleFunc) {
//	// Context创建
//
//	// 路由匹配
//
//	// 执行业务逻辑
//}

func (h *hTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *hTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

// ServeHTTP 处理请求的入口
func (h *hTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Context创建
	ctx := &Context{
		Req: request,
		Resp: writer,
	}
	// 路由匹配
	// 执行业务逻辑
	h.serve(ctx)
}

func (h *hTTPServer) serve(ctx *Context) {
	// 查找路由 执行业务逻辑
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.handler == nil {
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("not found"))
		return
	}
	n.handler(ctx)
}

func (h *hTTPServer) Start(addr string) error {
	l ,err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(l, h)
}

