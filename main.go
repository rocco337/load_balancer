package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {

	nodes := []string{"localhost:8082", "localhost:8081"}
	new(RoundRobinStrategy).Instance().Init(nodes)

	//open proxy server that forwards request to one of active nodes
	server := &http.Server{
		Addr:           ":8090",
		Handler:        http.HandlerFunc(handler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func handler(writer http.ResponseWriter, req *http.Request) {
	nextNode := new(RoundRobinStrategy).Instance().Next()
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Host = nextNode
			req.URL.Scheme = "http"
		},
	}
	fmt.Println("Forwarding request to: " + nextNode)
	proxy.ServeHTTP(writer, req)
}
