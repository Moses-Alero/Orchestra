package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"orchestra/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type ContainerBasicInfo struct{
	ID     string   `json:"id"`
	Name   string
	State  *types.ContainerState
}

type ContainerInfo struct{
	ContainerBasicInfo
	Created         string
	Path            string
	Args            []string
	Image           string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
}

func StartContainer(config *models.Config, respChan chan <- string, port string){
		ctx := context.Background()
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println(err)
			return
		}

		defer dockerClient.Close()
		containerConfig := &container.Config{
			Image: config.Spec.Template.Spec.Containers[0].Image,
			Cmd: []string{},
			Tty: false,
		}

		portBindings := nat.PortMap{
			"80/tcp": []nat.PortBinding{

				{
					HostIP: "0.0.0.0",
					HostPort: port,
				},
			}, 
		}
		
		hostConfig := &container.HostConfig{
			PortBindings: portBindings,
		}
		resp, err := dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Container Response: %v \n", resp)
		if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			fmt.Println(err)
			return
		}
		
		respChan  <- resp.ID
}

func GetContainerInfo(containerId string) (*ContainerBasicInfo, error){
	ctx := context.Background()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	
	defer dockerClient.Close()

	container, err := dockerClient.ContainerInspect(ctx, containerId)
	if err != nil {
		return nil, err
	}

	containerInfo := &ContainerBasicInfo{
		ID: container.ID,
		Name: container.Name,
		State: container.State,
	}

	return containerInfo, nil
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
			details, err :=dockerClient.ContainerInspect(ctx,container.ID)
			if err != nil {
				fmt.Println(err)
			}
			
			container := *&details.ContainerJSONBase

			containerInfo := ContainerBasicInfo{
				ID: container.ID,
				Name: container.Name,
				State: container.State,
			}
			d,_ := json.MarshalIndent(containerInfo, "", "  " )
			fmt.Println(string(d))
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
