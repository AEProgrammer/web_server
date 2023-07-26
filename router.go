package web_server

import (
	"strings"
)

type router struct {
	// 一个http方法对应一棵树
	// http method -> 路由数根结点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// AddRoute 有限制
// 类似 "/user/home"
// path必须以 / 开头 不能以 / 结尾 中间也不能有连续的///
func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {

	if path == "" {
		return
	}

	root, ok := r.trees[method]
	if !ok {
		// 没有根结点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 开头必须以/
	if path[0] != '/' {
		return
		//panic("路径必须以/开头")
	}

	// 结尾
	if path != "/" && path[len(path) - 1] == '/' {
		return
		//panic("路径不能以/结尾 ")
	}

	if path == "/" {
		root.handler = handleFunc
		return
	}

	// 切割path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		// 递归找准位置 中间节点如果不存在要新建节点
		child := root.childOrCreate(seg)
		root = child
	}
	root.handler = handleFunc
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}
	return child
}

type node struct {
	path string

	// path	到子节点的映射
	children map[string]*node

	// 用户注册的业务逻辑
	handler HandleFunc
}


