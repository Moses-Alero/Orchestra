package models

import (
	"fmt"
	"net/http"
	lb "orchestra/pkg/load-balancing"
)

type Cluster struct {
	PID          int
	Name         string
	ContainerIds []string
	Port         string
	LoadBalancer *lb.LoadBalancer
	ContainerMap map[string]string
}

// implement more info Display here
func (c *Cluster) ClusterInfo() *Cluster {
	return c
}

func (c *Cluster) StartProxy() {
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		c.LoadBalancer.ServeProxy(rw, req)
	}

	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'http://localhost:%s'\n", c.Port)
	http.ListenAndServe(":"+c.Port, nil)
}
