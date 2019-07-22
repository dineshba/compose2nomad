package main

// DockerComposeService represents a service in compose-file
type DockerComposeService struct {
	Image  string `json:"image"`
	Deploy struct {
		Resources struct {
			Limits struct {
				Cpus   string `json:"cpus"`
				Memory string `json:"memory"`
			} `json:"limits"`
		} `json:"resources"`
	} `json:"deploy"`
	Healthcheck struct {
		Test        []string `json:"test"`
		Interval    string   `json:"interval"`
		Timeout     string   `json:"timeout"`
		Retries     int      `json:"retries"`
		StartPeriod string   `json:"start_period"`
	} `json:"healthcheck"`
	Environment []string `json:"environment"`
}

// DockerComposeFileContent represents the compose-file
type DockerComposeFileContent struct {
	Version  string
	Services map[string]DockerComposeService
}
