package cluster

import (
	"errors"
	"fmt"
	"net/http"
	"orchestra/pkg/load-balancing"
)



type Cluster struct {
	Name   string
	ContainerIds  []string 
	Port   string
	LoadBalancer *lb.LoadBalancer
}


var clusters []*Cluster

func StoreClusterInfo(clusterName string, containers []string) {
	fmt.Println("New cluster")
	cluster := Cluster{
		Name: clusterName,
		ContainerIds: containers,
	}

	clusters = append(clusters, &cluster)
	fmt.Println("Cluster saved, name: ",clusterName)
}

func FindCluster(clusterName string) (*Cluster, error){
	for _, cluster := range clusters {
		if clusterName == cluster.Name{
			return cluster, nil	
		}
	}
		
	return nil, errors.New("Cluster not found")

} 

func SetProxy(clusterName string, ports []string) *Cluster{
	 	cluster, _ := FindCluster(clusterName)
		//setup port for cluster 
		cluster.LoadBalancer = lb.LoadBalance(cluster.Port, ports)
	  return 	cluster
}



func (c *Cluster) ClusterInfo() *Cluster {
	return c
}

func (c *Cluster) StartProxy(){
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		c.LoadBalancer.ServeProxy(rw, req)
	}

	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'http://localhost:%s'\n", c.Port)
	http.ListenAndServe(":"+c.Port, nil)
}
