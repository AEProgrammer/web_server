package web_server

import (
	"fmt"
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
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {

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
	if path != "/" && path[len(path)-1] == '/' {
		return
		//panic("路径不能以/结尾 ")
	}

	// 根结点
	if path == "/" {
		if root.handler != nil {
			panic(any("重复注册，路由冲突[/]"))
		}
		root.handler = handleFunc
		return
	}

	// 切割path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		// 不能有连续的//
		if seg == "" {
			return
		}
		// 递归找准位置 中间节点如果不存在要新建节点
		// 如果当前已经找到*节点
		if root.path == "*" {
			panic(any(fmt.Sprintf("禁止在*之后再注册路由[%s]", seg)))
		}
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic(any(fmt.Sprintf("路径冲突，重复注册[%s]", path)))
	}
	root.handler = handleFunc
}

// findRoute 优先考虑静态匹配 匹配不上就考虑通配符匹配
func (r *router) findRoute(method string, path string) (*node, bool) {
	// 深度查找路由树
	// 如果http method没有被注册
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	return root, root.handler != nil
}

func (n *node) childOrCreate(seg string) *node {

	// 参数路径匹配
	if seg[0] == ':' {
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	// 通配符匹配
	if seg == "*" {
		n.startChild = &node{
			path: seg,
		}
		return n.startChild
	}
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

// childOf 优先考虑静态匹配 匹配不上就考虑通配符匹配
func (n *node) childOf(path string) (*node, bool) {
	// 如果该节点根本没有children节点 则 优先考虑自身节点是不是paramChild再考虑该节点自身是不是startChild
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.startChild, n.startChild != nil
	}
	child, ok := n.children[path]
	// 如果匹配不到静态节点 优先考虑自身节点是不是paramChild再考虑自身节点是不是startChild
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true
		}
		return n.startChild, n.startChild != nil
	}
	return child, ok
}

type node struct {
	path string

	// path	到子节点的映射
	children map[string]*node

	// 通配符匹配
	startChild *node

	// 参数路径匹配
	paramChild *node

	// 用户注册的业务逻辑
	handler HandleFunc
}
