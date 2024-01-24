package models

type Config struct {
	ApiVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata,omitempty"`
	Spec       Spec     `yaml:"spec,omitempty"`
}

type Metadata struct {
	Name string `yaml:"name,omitempty"`
}

type Spec struct {
	Replicas int       `yaml:"replicas,omitempty"`
	Selector Seclector `yaml:"selector,omitempty"`
	Template Template  `yaml:"template,omitempty"`
}

type Seclector struct {
	MatchLabels MatchLabels `yaml:"matchLabels,omitempty"`
}

type MatchLabels struct {
	App string `yaml:"app,omitempty"`
}

type Template struct {
	Metadata TMetadata `yaml:"metadata,omitempty"`
	Spec     TSpec     `yaml:"spec,omitempty"`
}

type TMetadata struct {
	Labels MatchLabels `yaml:"labels,omitempty"`
}

type TSpec struct {
	Containers []Container `yaml:"containers,omitempty"`
}

type Container struct {
	Name  string  `yaml:"name,omitempty"`
	Image string  `yaml:"image,omitempty"`
	Ports []Ports `yaml:"ports,omitempty"`
}

type Ports struct {
	ContainerPort int `yaml:"containerPort"`
}
