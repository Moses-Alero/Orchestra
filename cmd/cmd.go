package cmd

import (
	"fmt"
	"os"
  
	"orchestra/pkg/docker"
	"orchestra/models"

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

		docker.StartContainer(&ymlConfig)

		fmt.Printf("The orchestra has started")
	},
}
