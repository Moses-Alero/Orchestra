package cmd

import (
	"fmt"
	"os"
	"strconv"

	"orchestra/models"
	"orchestra/pkg/cluster"
	"orchestra/pkg/docker"
	"orchestra/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var orchestra = &cobra.Command{
	Use:   "start",
	Short: "Start the orchestra",
	Long:  "Initiates the orchestra",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Tuning the orcehstra \n")

		if utils.CheckForOrchestraInfo() {
			fmt.Printf("rerer")
			docker.StopAllContainers()
			utils.RemoveOrchestraInfo()
		}

		filePath := args[0]

		data, err := os.ReadFile(filePath)

		var ymlConfig models.Config

		if err != nil {
			fmt.Println(err)
			return
		}

		err = yaml.Unmarshal(data, &ymlConfig)

		if err != nil {
			fmt.Println(err)
			return
		}

		respIds := make([]string, 0)
		ports := make([]string, 0)
		respChan := make(chan string)

		port := ymlConfig.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort

		for i := 1; i < ymlConfig.Spec.Replicas+1; i++ {
			containerPort := strconv.Itoa(port + i)
			containerPortAddr := "http://localhost:" + containerPort
			go docker.StartContainer(&ymlConfig, respChan, containerPort)
			resp := <-respChan
			ports = append(ports, containerPortAddr)
			respIds = append(respIds, resp)
		}

		clusterName := ymlConfig.Spec.Selector.MatchLabels.App

		close(respChan)

		cluster.StoreClusterInfo(clusterName, respIds, strconv.Itoa(port))
		c := cluster.SetProxy(clusterName, ports)
		c.StartProxy()
		fmt.Printf("The orchestra has started \n")
	},
}

var InspectContainer = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect container",
	Long:  "Inspect Containers By Name",
	Run: func(cmd *cobra.Command, args []string) {
		containerName := args[0]
		info := cluster.GetContainerInfo(containerName)
		fmt.Println(info)
	},
}

var OrchestraInfo = &cobra.Command{
	Use:   "info",
	Short: "Orchestra Info",
	Long:  "Get all the neccessary info about the orchestra(Cluster)",
	Run: func(cmd *cobra.Command, args []string) {
		c := cluster.Orchestra.ClusterInfo()
		fmt.Println(c)
	},
}

var listContainers = &cobra.Command{
	Use:   "list",
	Short: "lists all the containers running",
	Long:  "List all the containers running in the orchestra",
	Run: func(cmd *cobra.Command, args []string) {
		docker.ListContainer()
	},
}

var stopAllContainers = &cobra.Command{
	Use:   "stop",
	Short: "stops containers",
	Long:  "Stops all running containers in the orchestra",
	Run: func(cmd *cobra.Command, args []string) {
		docker.StopAllContainers()
		utils.RemoveOrchestraInfo()
	},
}
