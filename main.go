package main

import (
	"container/ring"
	"net/http"
	"net/http/httputil"
	"time"
)

var (
	nodesList *ring.Ring
)

func main() {

	nodes := []string{"localhost:8082", "localhost:8081"}
	nodesList = ring.New(len(nodes))
	for i := 0; i < nodesList.Len(); i++ {
		nodesList.Value = nodes[i]
		nodesList = nodesList.Next()
	}

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
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Host = nextNode()
			req.URL.Scheme = "http"
		},
	}

	proxy.ServeHTTP(writer, req)
}

func nextNode() string {
	nextValue := nodesList.Value.(string)
	nodesList = nodesList.Next()

	return nextValue
}
