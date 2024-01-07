package docker

import (
	"context"
	"fmt"
	"orchestra/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

)

func StartContainer(config *models.Config){
	ctx := context.Background()
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println(err)
			return
		}

		containerConfig := &container.Config{
			Hostname: config.Spec.Template.Spec.Containers[0].Name,
			Image: config.Spec.Template.Spec.Containers[0].Image,
			Cmd: []string{},
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

		fmt.Printf("Container ID: %s \n", resp.ID)

}
