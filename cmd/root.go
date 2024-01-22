package cmd

import (
	
	"github.com/spf13/cobra"
)

func init(){
	root.AddCommand(orchestra)
	root.AddCommand(InspectContainer)
	root.AddCommand(listContainers)
	root.AddCommand(stopAllContainers)
}

var root = &cobra.Command{
	Use: "orchestra",
	Short: "Container Orchestra",
	Long: "Container Orchestration tool for managing docker container clusters",
}

func Execute() error{
		return root.Execute()
}
