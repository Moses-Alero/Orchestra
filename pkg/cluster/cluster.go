package cluster

import (
	"fmt"
	"log"
	"net/http"
	"orchestra/pkg/docker"
	"orchestra/pkg/load-balancing"
)



type Cluster struct {
	Name   string
	ContainerIds  []string 
	Port   string
	LoadBalancer *lb.LoadBalancer
	ContainerMap  map[string]string
}


var Orchestra *Cluster


func StoreClusterInfo(clusterName string, containers []string, port string) {
	fmt.Println("New cluster")
	
	containerMap := make(map[string]string)
	for _, containerId := range containers{
		containerInfo, err := docker.GetContainerInfo(containerId)
		if err != nil{
			log.Fatal(err)
		}
		containerMap[containerInfo.Name] = containerId
	}


Orchestra = &Cluster{
		Name: clusterName,
		ContainerIds: containers,
		Port: port,
		ContainerMap: containerMap,
	}
	fmt.Println(Orchestra.ContainerMap)
	fmt.Println("Cluster saved, name: ",clusterName)
}


func SetProxy(clusterName string, ports []string) *Cluster{
		//setup port for cluster 
		Orchestra.LoadBalancer = lb.LoadBalance(Orchestra.Port, ports)
	  return 	Orchestra
}

func GetOrchestra() *Cluster{
	return Orchestra
}

//implement more info Display herer
func (c *Cluster) ClusterInfo() *Cluster {
	return c
}

func  GetContainerInfo(containerName string){
		fmt.Println("Map: ",Orchestra.ContainerMap)
		//val, ok := Orchestra.ContainerMap[containerName]	
		//if !ok {
		// log.Fatal("no container named ", containerName)
		//} 
		//info, err := docker.GetContainerInfo(val)
		//if err != nil{
		//	log.Fatal(err)
		//}
		//return info

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
