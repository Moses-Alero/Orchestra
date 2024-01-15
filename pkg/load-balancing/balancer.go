package lb

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoadBalancer struct {
	Servers   []*Server
	RoundRobinCount int
	Port      string
}

type Server struct {
  Addr string
	Proxy *httputil.ReverseProxy
}

func (s *Server) Address() string { return s.Addr }
func (s *Server) IsAlive() bool { return true }
func (s *Server) Serve(rw http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(rw, req)
}


func (lb *LoadBalancer) getNextAvailableServer() *Server {
	server := lb.Servers[lb.RoundRobinCount%len(lb.Servers)]
	lb.RoundRobinCount++
	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()

	// could optionally log stuff about the request here!
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())

	// could delete pre-existing X-Forwarded-For header to prevent IP spoofing
	targetServer.Serve(rw, req)
}


func  NewLoadBalancer(port string, servers []*Server) *LoadBalancer{
	return &LoadBalancer{
		Port: port,
		RoundRobinCount: 0,
		Servers: servers,
	}
}

func createContainerServer(addr string) *Server {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Println(err)
	}

	return &Server{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}



func LoadBalance(clusterPort string, ports []string) *LoadBalancer{
	servers := make([]*Server, 0)
	for _,port := range ports{
		server := createContainerServer(port)
		servers = append(servers, server)
	}

  lb := NewLoadBalancer(clusterPort, servers)	

	return lb
}


