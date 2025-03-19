package parser

type ComposeConfig struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	Networks map[string]Network `yaml:"networks,omitempty"`
	Volumes  map[string]Volume  `yaml:"volumes,omitempty"`
	Configs  map[string]Config  `yaml:"configs,omitempty"`
	Secrets  map[string]Secret  `yaml:"secrets,omitempty"`
}

type Service struct {
	Image       string       `yaml:"image"`
	Ports       []string     `yaml:"ports,omitempty"`
	Networks    []string     `yaml:"networks,omitempty"`
	Deploy      DeployConfig `yaml:"deploy,omitempty"`
	Volumes     []string     `yaml:"volumes,omitempty"`
	Configs     []ConfigRef  `yaml:"configs,omitempty"`
	Secrets     []SecretRef  `yaml:"secrets,omitempty"`
	Environment []string     `yaml:"environment,omitempty"`
	Healthcheck *Healthcheck `yaml:"healthcheck,omitempty"`
}

type Network struct {
	Driver     string            `yaml:"driver,omitempty"`
	External   bool              `yaml:"external,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
	Attachable bool              `yaml:"attachable,omitempty"`
	Internal   bool              `yaml:"internal,omitempty"`
}

type Volume struct {
	Driver   string            `yaml:"driver,omitempty"`
	External bool              `yaml:"external,omitempty"`
	Labels   map[string]string `yaml:"labels,omitempty"`
}

type Config struct {
	File string `yaml:"file,omitempty"`
}

type Secret struct {
	File string `yaml:"file,omitempty"`
}

type DeployConfig struct {
	Mode           string         `yaml:"mode,omitempty"`
	Replicas       int            `yaml:"replicas,omitempty"`
	UpdateConfig   *UpdateConfig  `yaml:"update_config,omitempty"`
	RollbackConfig *UpdateConfig  `yaml:"rollback_config,omitempty"`
	RestartPolicy  *RestartPolicy `yaml:"restart_policy,omitempty"`
	Placement      *Placement     `yaml:"placement,omitempty"`
	Resources      *Resources     `yaml:"resources,omitempty"`
}

type UpdateConfig struct {
	Parallelism int    `yaml:"parallelism,omitempty"`
	Delay       string `yaml:"delay,omitempty"`
	Order       string `yaml:"order,omitempty"`
}

type RestartPolicy struct {
	Condition   string `yaml:"condition,omitempty"`
	Delay       string `yaml:"delay,omitempty"`
	MaxAttempts int    `yaml:"max_attempts,omitempty"`
	Window      string `yaml:"window,omitempty"`
}

type Placement struct {
	Constraints []string              `yaml:"constraints,omitempty"`
	Preferences []PlacementPreference `yaml:"preferences,omitempty"`
}

type PlacementPreference struct {
	Spread string `yaml:"spread,omitempty"`
}

type Resources struct {
	Limits       *ResourceSpec `yaml:"limits,omitempty"`
	Reservations *ResourceSpec `yaml:"reservations,omitempty"`
}

type ResourceSpec struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type ConfigRef struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type SecretRef struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

type Healthcheck struct {
	Test        []string `yaml:"test,omitempty"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}
