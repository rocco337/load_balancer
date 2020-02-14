package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"net/http/httputil"
	"time"
)

var serverPort int
var balancingStrategy string
var sessionIDCookieName = "SessionId"

const roundRobinBalancingStrategy = "round-robin"
const stickySessionBalancingStrategy = "sticky-session"

func init() {
	flag.IntVar(&serverPort, "p", 8090, "default loadbalancer port")
	flag.StringVar(&balancingStrategy, "s", "round-robin", "balancing stragey: round-robin, sticky-session")
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
	fmt.Printf("Loadbalancing stragey: %s \n", balancingStrategy)

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
	nextNode := getBackendNodeURL(balancingStrategy, req)
	fmt.Println("Forwarding request to: " + nextNode)

	proxyConnections[nextNode].ServeHTTP(writer, req)
}

func getBackendNodeURL(balancingStrategy string, req *http.Request) string {
	var nextNodeURL string

	if balancingStrategy == roundRobinBalancingStrategy {
		nextNodeURL = new(RoundRobinStrategy).Instance().Next()
	} else if balancingStrategy == stickySessionBalancingStrategy {
		sessionIDCookie, err := readSessionCookie(req)
		if err != nil {
			//cannot read cookie, fallback to roundRobin
			nextNodeURL = getBackendNodeURL(roundRobinBalancingStrategy, req)
		}

		nextNodeURL = new(StickySessionStrategy).Instance().Next(sessionIDCookie)
	} else {
		panic("Unknown load balancing strategy" + balancingStrategy)
	}

	return nextNodeURL

}

func readSessionCookie(req *http.Request) (string, error) {
	sessionIDCookie, err := req.Cookie(sessionIDCookieName)
	if err != nil {
		return "", errors.New("")
	}

	fmt.Println("Sticky sessionId: ", sessionIDCookie.Value)
	return sessionIDCookie.Value, nil
}
