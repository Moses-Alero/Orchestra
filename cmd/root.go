package cmd

import (
	
	"github.com/spf13/cobra"
)

func init(){
	root.AddCommand(orchestra)
}

var root = &cobra.Command{
	Use: "orchestra",
	Short: "Container Orchestra",
	Long: "Container Orcheestration tool for managing docker container clusters",
}

func Execute() error{
		return root.Execute()
}
