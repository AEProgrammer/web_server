package web_server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// 构造路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail/:id",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/*",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/*/abc",
		//},
		//{
		//	method: http.MethodPost,
		//	path : "login",
		//},
		//{
		//	method: http.MethodPost,
		//	path : "login////",
		//},
	}
	mockHandler := func(ctx *Context) {}

	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 断言路由树与预期一摸一样
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						startChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
				startChild: &node{
					path:    "*",
					handler: mockHandler,
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:    "create",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}

	// 断言两者相等
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	//r = newRouter()
	//assert.Panics(t, func() {
	//	r.addRoute(http.MethodGet, "",  mockHandler)
	//})
}

func TestRouter_addRoute_failed(t *testing.T) {

	mockHandler := func(ctx *Context) {}

	r := newRouter()
	r.addRoute(http.MethodGet, "", mockHandler)
	assert.True(t, len(r.trees) == 0, "")

	r = newRouter()
	r.addRoute(http.MethodGet, "login", mockHandler)
	assert.True(t, len(r.trees) == 1, "")
	assert.True(t, r.trees["GET"].path == "/", "")

	r = newRouter()
	r.addRoute(http.MethodGet, "/login///", mockHandler)
	assert.True(t, len(r.trees) == 1, "")
	assert.True(t, r.trees["GET"].path == "/", "")

	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "重复注册，路由冲突[/]")

	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "重复注册，路由冲突[/a/b/c]")

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/*/abc", mockHandler)
	}, "禁止在*之后再注册路由[ab]")
}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http方法"), false
		}
		// v, dst要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	// 检查通配符路径节点是否一致
	if (n.startChild == nil && y.startChild != nil) || (n.startChild != nil && y.startChild == nil) {
		return fmt.Sprintf("通配符节点数量不匹配"), false
	} else if n.startChild != nil && y.startChild != nil {
		msg, ok := n.startChild.equal(y.startChild)
		if !ok {
			return msg, ok
		}
	}

	// 检查参数路径节点是否一致
	if (n.paramChild == nil && y.paramChild != nil) || (n.paramChild != nil && y.paramChild == nil) {
		return fmt.Sprintf("参数路径节点数量不匹配"), false
	} else if n.paramChild != nil && y.paramChild != nil {
		msg, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return msg, ok
		}
	}

	// 比较handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		//{
		//	method: http.MethodGet,
		//	path : "/",
		//},
		//{
		//	method: http.MethodGet,
		//	path : "/user",
		//},
		//{
		//	method: http.MethodGet,
		//	path : "/user/home",
		//},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		//{
		//	method: http.MethodPost,
		//	path : "/order/create",
		//},
	}
	r := newRouter()
	mockHandler := func(ctx *Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCase := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{
			// 方法不存在
			name:      "method notfound",
			method:    http.MethodHead,
			path:      "/order/detail",
			wantFound: false,
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			wantNode: &node{
				path:    "detail",
				handler: mockHandler,
			},
		},
		{
			// 通过通配符命中
			name:      "order *",
			method:    http.MethodGet,
			path:      "/order/abc",
			wantFound: true,
			wantNode: &node{
				path:    "*",
				handler: mockHandler,
			},
		},
		{
			name:      "root node",
			method:    http.MethodGet,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			// 断言两者相等
			msg, ok := tc.wantNode.equal(n)
			assert.True(t, ok, msg)
			//assert.Equal(t, tc.wantNode.path, n.path)
			//assert.Equal(t, tc.wantNode.children, n.children)
			//nHandler := reflect.ValueOf(n.handler)
			//yHandler := reflect.ValueOf(tc.wantNode.handler)
			//assert.True(t, nHandler == yHandler)
		})
	}
}
