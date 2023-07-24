package web_server

type router struct {
	// 一个http方法对应一棵树
	// http method -> 路由数根结点
	trees map[string]*node
}

func newRouter() *router {
	return &router{}
}

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {

}

type node struct {
	path string

	// path	到子节点的映射
	children map[string]*node

	// 用户注册的业务逻辑
	handler HandleFunc
}
