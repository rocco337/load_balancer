package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

var serverPort int

func init() {
	flag.IntVar(&serverPort, "p", 8090, "default loadbalancer port")
}

//auto discovery
//configuration - port, autodiscovery
//server health - remove from list if unhealthy
func main() {
	flag.Parse()

	nodes := []string{"localhost:8082", "localhost:8081"}
	new(RoundRobinStrategy).Instance().Init(nodes)

	fmt.Println("Starting loadbalancer...")
	fmt.Printf("Listening to port %v \n", serverPort)

	//open proxy server that forwards request to one of active nodes
	server := &http.Server{
		Addr:           fmt.Sprintf(":%v", serverPort),
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
