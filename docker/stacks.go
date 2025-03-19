package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/dockrelix/dockrelix-backend/database"
	"github.com/dockrelix/dockrelix-backend/models"
)

func ListStacks(cli *client.Client) []models.Stack {
	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		log.Fatalf("Error fetching services: %v", err)
	}

	networks, err := cli.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		panic(err)
	}
	volumes, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching volumes: %v", err)
	}

	stacks := make(map[string]*models.Stack)

	for _, service := range services {
		stackName := service.Spec.Labels["com.docker.stack.namespace"]
		if stackName == "" {
			continue
		}

		if _, exists := stacks[stackName]; !exists {
			stacks[stackName] = &models.Stack{
				Name:     stackName,
				Services: []models.ServiceData{},
				Networks: []models.NetworkData{},
				Volumes:  []models.VolumeData{},
			}
		}

		serviceData := models.ServiceData{
			ID:    service.ID,
			Name:  service.Spec.Name,
			Image: service.Spec.TaskTemplate.ContainerSpec.Image,
		}

		if service.Spec.Mode.Replicated != nil {
			replicas := int64(*service.Spec.Mode.Replicated.Replicas)
			serviceData.Replicas = &replicas
		}

		if service.Endpoint.Ports != nil {
			for _, port := range service.Endpoint.Ports {
				portData := models.PortData{
					Target:    int64(port.TargetPort),
					Published: int64(port.PublishedPort),
					Mode:      string(port.PublishMode),
					Protocol:  string(port.Protocol),
				}
				serviceData.Ports = append(serviceData.Ports, portData)
			}
		}

		stacks[stackName].Services = append(stacks[stackName].Services, serviceData)
	}

	for _, network := range networks {
		stackName := network.Labels["com.docker.stack.namespace"]
		if stackName != "" {
			if _, exists := stacks[stackName]; exists {
				networkData := models.NetworkData{
					ID:     network.ID,
					Name:   network.Name,
					Scope:  network.Scope,
					Driver: network.Driver,
					Labels: network.Labels,
				}
				stacks[stackName].Networks = append(stacks[stackName].Networks, networkData)
			}
		}
	}

	for _, volume := range volumes.Volumes {
		stackName := volume.Labels["com.docker.stack.namespace"]
		if stackName != "" {
			if _, exists := stacks[stackName]; exists {
				volumeData := models.VolumeData{
					Name:       volume.Name,
					Driver:     volume.Driver,
					Labels:     volume.Labels,
					Mountpoint: volume.Mountpoint,
				}
				stacks[stackName].Volumes = append(stacks[stackName].Volumes, volumeData)
			}
		}
	}

	var result []models.Stack
	for _, stack := range stacks {
		result = append(result, *stack)
	}

	return result
}

func SaveDraft(stackDraft models.StackDraft) error {
	if err := database.DB.Create(&models.StackDraft{Name: stackDraft.Name, Data: stackDraft.Data}).Error; err != nil {
		return err
	}

	return nil
}

func GetDrafts() []models.StackDraft {
	var drafts []models.StackDraft
	database.DB.Find(&drafts)
	return drafts
}
