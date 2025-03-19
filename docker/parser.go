package docker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/dockrelix/dockrelix-backend/models/parser"
	"gopkg.in/yaml.v3"
)

func RemoveStackFromName(name, stackName string) string {
	output, _ := strings.CutPrefix(name, stackName+"_")
	return output
}

func getServiceMode(srv swarm.Service) string {
	if srv.Spec.Mode.Replicated != nil {
		return "replicated"
	}
	return "global"
}

func GenerateStackConfig(services []swarm.Service, networks []swarm.Network, volumes []*volume.Volume, secrets []swarm.Secret, configs []swarm.Config, stackName string) ([]byte, error) {
	config := parser.ComposeConfig{
		Version:  "3.8",
		Services: make(map[string]parser.Service),
		Networks: make(map[string]parser.Network),
		Volumes:  make(map[string]parser.Volume),
		Configs:  make(map[string]parser.Config),
		Secrets:  make(map[string]parser.Secret),
	}

	for _, srv := range services {
		service := parser.Service{
			Image: srv.Spec.TaskTemplate.ContainerSpec.Image,
			Deploy: parser.DeployConfig{
				Mode: getServiceMode(srv),
			},
			Networks: []string{},
		}

		if srv.Spec.Mode.Replicated != nil && srv.Spec.Mode.Replicated.Replicas != nil {
			service.Deploy.Replicas = int(*srv.Spec.Mode.Replicated.Replicas)
		}

		for _, port := range srv.Endpoint.Ports {
			service.Ports = append(service.Ports, fmt.Sprintf("%d:%d/%s", port.PublishedPort, port.TargetPort, port.Protocol))
		}

		for _, nw := range srv.Spec.TaskTemplate.Networks {
			for _, net := range networks {
				if nw.Target == net.ID {
					service.Networks = append(service.Networks, RemoveStackFromName(net.Spec.Name, srv.Spec.Name))
				}
			}
		}

		for _, mount := range srv.Spec.TaskTemplate.ContainerSpec.Mounts {
			service.Volumes = append(service.Volumes, fmt.Sprintf("%s:%s", RemoveStackFromName(mount.Source, stackName), mount.Target))
		}

		for _, cfg := range srv.Spec.TaskTemplate.ContainerSpec.Configs {
			service.Configs = append(service.Configs, parser.ConfigRef{
				Source: RemoveStackFromName(cfg.ConfigName, stackName),
				Target: cfg.File.Name,
			})
		}

		for _, secret := range srv.Spec.TaskTemplate.ContainerSpec.Secrets {
			service.Secrets = append(service.Secrets, parser.SecretRef{
				Source: RemoveStackFromName(secret.SecretName, stackName),
				Target: secret.File.Name,
			})
		}

		service.Environment = append(service.Environment, srv.Spec.TaskTemplate.ContainerSpec.Env...)

		if hc := srv.Spec.TaskTemplate.ContainerSpec.Healthcheck; hc != nil {
			service.Healthcheck = &parser.Healthcheck{
				Test:        hc.Test,
				Interval:    hc.Interval.String(),
				Timeout:     hc.Timeout.String(),
				Retries:     int(hc.Retries),
				StartPeriod: hc.StartPeriod.String(),
			}
		}

		if updateConfig := srv.Spec.UpdateConfig; updateConfig != nil {
			service.Deploy.UpdateConfig = &parser.UpdateConfig{
				Parallelism: int(updateConfig.Parallelism),
				Delay:       updateConfig.Delay.String(),
				Order:       string(updateConfig.Order),
			}
		}

		if rollbackConfig := srv.Spec.RollbackConfig; rollbackConfig != nil {
			service.Deploy.RollbackConfig = &parser.UpdateConfig{
				Parallelism: int(rollbackConfig.Parallelism),
				Delay:       rollbackConfig.Delay.String(),
				Order:       string(rollbackConfig.Order),
			}
		}

		if restartPolicy := srv.Spec.TaskTemplate.RestartPolicy; restartPolicy != nil {
			service.Deploy.RestartPolicy = &parser.RestartPolicy{
				Condition:   string(restartPolicy.Condition),
				Delay:       restartPolicy.Delay.String(),
				MaxAttempts: int(*restartPolicy.MaxAttempts),
				Window:      restartPolicy.Window.String(),
			}
		}

		if placement := srv.Spec.TaskTemplate.Placement; placement != nil {
			if len(placement.Constraints) > 0 || len(placement.Preferences) > 0 {
				service.Deploy.Placement = &parser.Placement{
					Constraints: placement.Constraints,
				}
				for _, pref := range placement.Preferences {
					service.Deploy.Placement.Preferences = append(service.Deploy.Placement.Preferences, parser.PlacementPreference{
						Spread: pref.Spread.SpreadDescriptor,
					})
				}
			}
		}

		if resources := srv.Spec.TaskTemplate.Resources; resources != nil {
			if resources.Limits != nil || resources.Reservations != nil {
				service.Deploy.Resources = &parser.Resources{}
				if resources.Limits != nil {
					service.Deploy.Resources.Limits = &parser.ResourceSpec{
						CPUs:   strconv.FormatInt(resources.Limits.NanoCPUs/1e9, 10),
						Memory: strconv.FormatInt(resources.Limits.MemoryBytes/(1024*1024), 10) + "M",
					}
				}
				if resources.Reservations != nil {
					service.Deploy.Resources.Reservations = &parser.ResourceSpec{
						CPUs:   strconv.FormatInt(resources.Reservations.NanoCPUs/1e9, 10),
						Memory: strconv.FormatInt(resources.Reservations.MemoryBytes/(1024*1024), 10) + "M",
					}
				}
			}
		}

		config.Services[srv.Spec.Name] = service
	}

	for _, net := range networks {
		labels := net.Spec.Labels
		delete(labels, "com.docker.stack.namespace")
		config.Networks[RemoveStackFromName(net.Spec.Annotations.Name, stackName)] = parser.Network{
			Driver:     net.Spec.DriverConfiguration.Name,
			External:   !net.Spec.Internal,
			Labels:     labels,
			Attachable: net.Spec.Attachable,
			Internal:   net.Spec.Internal,
		}
	}

	for _, vol := range volumes {
		labels := vol.Labels
		delete(labels, "com.docker.stack.namespace")
		config.Volumes[RemoveStackFromName(vol.Name, stackName)] = parser.Volume{
			Driver:   vol.Driver,
			External: vol.Options["external"] == "true",
			Labels:   labels,
		}
	}

	for _, secret := range secrets {
		labels := secret.Spec.Labels
		delete(labels, "com.docker.stack.namespace")
		config.Secrets[RemoveStackFromName(secret.Spec.Name, stackName)] = parser.Secret{
			File: secret.Spec.Name,
		}
	}

	for _, cfg := range configs {
		labels := cfg.Spec.Labels
		delete(labels, "com.docker.stack.namespace")
		config.Configs[RemoveStackFromName(cfg.Spec.Name, stackName)] = parser.Config{
			File: cfg.Spec.Name,
		}
	}

	return yaml.Marshal(config)
}

func ParseStackConfig(cli *client.Client, stackName string) (parser.ComposeConfig, error) {
	services, err := cli.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "com.docker.stack.namespace="+stackName)),
	})
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	networks, err := cli.NetworkList(context.Background(), network.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "com.docker.stack.namespace="+stackName)),
	})
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	volumes, err := cli.VolumeList(context.Background(), volume.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "com.docker.stack.namespace="+stackName)),
	})
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	secrets, err := cli.SecretList(context.Background(), types.SecretListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "com.docker.stack.namespace="+stackName)),
	})
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	configs, err := cli.ConfigList(context.Background(), types.ConfigListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "com.docker.stack.namespace="+stackName)),
	})
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	var sanitizedServices []swarm.Service
	for _, srv := range services {
		sanitizedSrv := srv
		sanitizedSrv.Spec.Name = RemoveStackFromName(srv.Spec.Name, stackName)
		sanitizedSrv.Spec.TaskTemplate.ContainerSpec.Image = strings.Split(srv.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
		sanitizedServices = append(sanitizedServices, sanitizedSrv)
	}

	var networksList []swarm.Network
	for _, net := range networks {
		networksList = append(networksList, swarm.Network{
			ID: net.ID,
			Spec: swarm.NetworkSpec{
				Annotations: swarm.Annotations{
					Name:   RemoveStackFromName(net.Name, stackName),
					Labels: net.Labels,
				},
				DriverConfiguration: &swarm.Driver{
					Name: net.Driver,
				},
				Scope:      net.Scope,
				Attachable: net.Attachable,
				Internal:   net.Internal,
			},
		})
	}

	var sanitizedVolumes []*volume.Volume
	for _, vol := range volumes.Volumes {
		sanitizedVol := *vol
		sanitizedVol.Name = RemoveStackFromName(vol.Name, stackName)
		sanitizedVolumes = append(sanitizedVolumes, &sanitizedVol)
	}

	var sanitizedSecrets []swarm.Secret
	for _, secret := range secrets {
		sanitizedSecret := secret
		sanitizedSecret.Spec.Name = RemoveStackFromName(secret.Spec.Name, stackName)
		sanitizedSecrets = append(sanitizedSecrets, sanitizedSecret)
	}

	var sanitizedConfigs []swarm.Config
	for _, cfg := range configs {
		sanitizedCfg := cfg
		sanitizedCfg.Spec.Name = RemoveStackFromName(cfg.Spec.Name, stackName)
		sanitizedConfigs = append(sanitizedConfigs, sanitizedCfg)
	}

	configBytes, err := GenerateStackConfig(sanitizedServices, networksList, sanitizedVolumes, sanitizedSecrets, sanitizedConfigs, stackName)
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	var config parser.ComposeConfig
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return parser.ComposeConfig{}, err
	}

	return config, nil
}
