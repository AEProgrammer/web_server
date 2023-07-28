package web_server

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	
	h := NewHTTPServer()

	h.Get("/user", func(ctx *Context) {
	})

	h.Get("/order/detail", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, order detail~"))
	})

	http.ListenAndServe(":8081", h)
}
