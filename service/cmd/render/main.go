package main

import (
	"fmt"

	"github.com/domsolutions/http2"
	"github.com/valyala/fasthttp"
)

func main() {
	s := &fasthttp.Server{
		Handler: requestHandler,
		Name:    "HTTP2 test",
	}

	http2.ConfigureServer(s, http2.ServerConfig{})

	s.ListenAndServeTLS(":8080", "./cert.pem", "./key.pem")
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	fmt.Printf(" %s \n", &ctx.Request)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
