package docker_test

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/dockrelix/dockrelix-backend/docker"
)

type MockClient struct {
	ServiceListFunc func(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error)
	NetworkListFunc func(ctx context.Context, options network.ListOptions) ([]network.Summary, error)
	VolumeListFunc  func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	SecretListFunc  func(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error)
	ConfigListFunc  func(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error)
}

func (m *MockClient) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	return m.ServiceListFunc(ctx, options)
}

func (m *MockClient) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
	return m.NetworkListFunc(ctx, options)
}

func (m *MockClient) VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
	return m.VolumeListFunc(ctx, options)
}

func (m *MockClient) SecretList(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
	return m.SecretListFunc(ctx, options)
}

func (m *MockClient) ConfigList(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
	return m.ConfigListFunc(ctx, options)
}

func TestParseStackConfig(t *testing.T) {
	mockClient := &MockClient{
		ServiceListFunc: func(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
			return []swarm.Service{
				{
					Spec: swarm.ServiceSpec{
						Annotations: swarm.Annotations{Name: "test_stack_service1"},
						TaskTemplate: swarm.TaskSpec{
							ContainerSpec: &swarm.ContainerSpec{
								Image: "nginx:latest",
							},
						},
					},
				},
			}, nil
		},
		NetworkListFunc: func(ctx context.Context, options network.ListOptions) ([]network.Summary, error) {
			return []network.Summary{
				{
					ID: "network_id_1",

					Name: "test_stack_network1",
				},
			}, nil
		},
		VolumeListFunc: func(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error) {
			return volume.ListResponse{
				Volumes: []*volume.Volume{
					{Name: "test_stack_volume1"},
				},
			}, nil
		},
		SecretListFunc: func(ctx context.Context, options types.SecretListOptions) ([]swarm.Secret, error) {
			return []swarm.Secret{
				{
					Spec: swarm.SecretSpec{
						Annotations: swarm.Annotations{Name: "test_stack_secret1"},
					},
				},
			}, nil
		},
		ConfigListFunc: func(ctx context.Context, options types.ConfigListOptions) ([]swarm.Config, error) {
			return []swarm.Config{
				{
					Spec: swarm.ConfigSpec{
						Annotations: swarm.Annotations{Name: "test_stack_config1"},
					},
				},
			}, nil
		},
	}

	stackName := "test_stack"

	config, err := docker.ParseStackConfig(mockClient, stackName)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if _, ok := config.Services["service1"]; !ok {
		t.Errorf("expected service1 to be present in services")
	}

	if _, ok := config.Networks["network1"]; !ok {
		t.Errorf("expected network1 to be present in networks")
	}

	if _, ok := config.Volumes["volume1"]; !ok {
		t.Errorf("expected volume1 to be present in volumes")
	}

	if _, ok := config.Secrets["secret1"]; !ok {
		t.Errorf("expected secret1 to be present in secrets")
	}

	if _, ok := config.Configs["config1"]; !ok {
		t.Errorf("expected config1 to be present in configs")
	}

	if config.Services["service1"].Image != "nginx:latest" {
		t.Errorf("expected image to be 'nginx:latest', got %s", config.Services["service1"].Image)
	}
}
