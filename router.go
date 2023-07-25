package web_server

import "strings"

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

func (r *router) AddRoute(method string, path string, handleFunc HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		// 没有根结点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 切割path
	path = path[1:]
	segs := strings.Split(path, "/")
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


