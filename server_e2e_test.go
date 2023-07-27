package web_server

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	var s Server = &HTTPServer{}
	http.ListenAndServe(":8081", s)
	
	s.addRoute(http.MethodGet, "get", func(ctx *Context) {
		
	})
	
	h := &HTTPServer{}
	h.Get("/user", func(ctx *Context) {
		
	})
}
