package cmd

import (
	"fmt"
	"os"

	"orchestra/models"
	"orchestra/pkg/docker"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)



var orchestra = &cobra.Command{
	Use: "start",
	Short: "Start the orchestra",
	Long: "Initiatest the orchestra",
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
		
		respIds := make([]interface{}, 0)
    respChan := make(chan interface{})


		for i := 0; i <  ymlConfig.Spec.Replicas; i++{
 			go docker.StartContainer(&ymlConfig, respChan)
			resp := <-respChan
			respIds = append(respIds, resp)
		}

		close(respChan)
		fmt.Println(respIds)

		fmt.Printf("The orchestra has started")
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
