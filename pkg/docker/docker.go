package docker

import (
	"context"
	"fmt"
	"orchestra/models"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
  "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

)

func StartContainer(config *models.Config, respChan chan <- interface{}){
	ctx := context.Background()
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println(err)
			return
		}

		defer dockerClient.Close()

		containerConfig := &container.Config{
			Domainname: config.Spec.Template.Spec.Containers[0].Name,
			Image: config.Spec.Template.Spec.Containers[0].Image,
			Cmd: []string{},
			Tty: false,
		}

		resp, err := dockerClient.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Continer Response: %v \n", resp)
		if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			fmt.Println(err)
			return
		}
		
		respChan  <- resp

		fmt.Printf("Container ID: %s \n", resp.ID)
  
}

func ListContainer(){
		ctx := context.Background()
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer dockerClient.Close()

		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			fmt.Println(container.ID)
		}

}

func StopAllContainers(){
		ctx := context.Background()
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer dockerClient.Close()

		containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			fmt.Print("Stopping container ", container.ID[:10], "... ")
			noWaitTimeout := 0 // to not wait for the container to exit gracefully
     	stpOtions := containertypes.StopOptions{Timeout: &noWaitTimeout}
			if err := dockerClient.ContainerStop(ctx, container.ID, stpOtions); err != nil {
				panic(err)
			}
			fmt.Println("Success")
		}
}
