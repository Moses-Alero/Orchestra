package cmd

import (
	"fmt"
	"os"
	"strconv"

	"orchestra/models"
	"orchestra/pkg/cluster"
	"orchestra/pkg/docker"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)



var orchestra = &cobra.Command{
	Use: "start",
	Short: "Start the orchestra",
	Long: "Initiates the orchestra",
	Run:func(cmd *cobra.Command, args []string){
		fmt.Printf("Tuning the orcehstra \n")	
			
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

		for i := 1; i <  ymlConfig.Spec.Replicas + 1; i++{ 
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


var listContainers = &cobra.Command{
	Use: "list",
	Short: "lists all the containers running",
	Long: "List all the containers running in the orchestra",
	Run: func(cmd *cobra.Command, args []string){
	  docker.ListContainer()	
	},
}


var stopAllContainers = &cobra.Command{
	Use: "stop",
	Short: "stops containers",
	Long: "Stops all running containers in the orchestra",
	Run: func(cmd *cobra.Command, args []string){
		docker.StopAllContainers()
	},
} 
