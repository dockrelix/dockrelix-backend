package models

type ServiceData struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Replicas *int64     `json:"replicas,omitempty"`
	Image    string     `json:"image"`
	Ports    []PortData `json:"ports,omitempty"`
}

type PortData struct {
	Target    int64  `json:"target"`
	Published int64  `json:"published"`
	Mode      string `json:"mode"`
	Protocol  string `json:"protocol"`
}

type NetworkData struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Scope      string            `json:"scope"`
	Driver     string            `json:"driver"`
	Containers []string          `json:"containers,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type VolumeData struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Labels     map[string]string `json:"labels,omitempty"`
	Mountpoint string            `json:"mountpoint,omitempty"`
}

type Stack struct {
	Name     string        `json:"name"`
	Services []ServiceData `json:"services"`
	Networks []NetworkData `json:"networks"`
	Volumes  []VolumeData  `json:"volumes"`
}
