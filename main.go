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

var proxyConnections map[string]*httputil.ReverseProxy

//auto discovery
//configuration - port, autodiscovery
//server health - remove from list if unhealthy
func main() {
	proxyConnections = make(map[string]*httputil.ReverseProxy, 0)

	flag.Parse()

	nodes := []string{"localhost:8082", "localhost:8081"}
	new(RoundRobinStrategy).Instance().Init(nodes)
	warmUpConnections(nodes)

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

func warmUpConnections(nodes []string) {
	for _, node := range nodes {
		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Host = node
				req.URL.Scheme = "http"
			},
		}

		proxyConnections[node] = proxy
	}
}

func handler(writer http.ResponseWriter, req *http.Request) {
	nextNode := new(RoundRobinStrategy).Instance().Next()
	fmt.Println("Forwarding request to: " + nextNode)

	proxyConnections[nextNode].ServeHTTP(writer, req)
}
