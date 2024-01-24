package cluster

import (
	"encoding/json"
	"fmt"
	"log"
	"orchestra/models"
	"orchestra/pkg/docker"
	"orchestra/pkg/load-balancing"
	"orchestra/utils"
)

var Orchestra *models.Cluster

func StoreClusterInfo(clusterName string, containers []string, port string) {
	fmt.Println("New cluster")

	containerMap := make(map[string]string)
	for _, containerId := range containers {
		containerInfo, err := docker.GetContainerInfo(containerId)
		if err != nil {
			log.Fatal(err)
		}
		containerMap[containerInfo.Name] = containerId
	}

	Orchestra = &models.Cluster{
		Name:         clusterName,
		ContainerIds: containers,
		Port:         port,
		ContainerMap: containerMap,
	}

	utils.StoreOrchestraInfo(Orchestra)

	fmt.Println(Orchestra.ContainerMap)
	fmt.Println("Cluster saved, name: ", clusterName)

}

func GetContainerInfo(containerName string) *docker.ContainerBasicInfo {
	orchestra := utils.GetOrchestraInfo()

	val, ok := orchestra.ContainerMap[containerName]
	if !ok {
		log.Fatal("no container named ", containerName)
	}
	info, err := docker.GetContainerInfo(val)
	if err != nil {
		log.Fatal(err)
	}

	jsonOutput, _ := json.MarshalIndent(info, "", " ")
	fmt.Println(string(jsonOutput))
	return info

}

func SetProxy(clusterName string, ports []string) *models.Cluster {
	//setup port for cluster
	Orchestra.LoadBalancer = lb.LoadBalance(Orchestra.Port, ports)
	return Orchestra
}

func GetOrchestra() *models.Cluster {
	return Orchestra
}
