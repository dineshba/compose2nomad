package main

import (
	"strconv"
	"strings"
)

// Group represents nomad group
type Group struct {
	Count int             `json:"count"`
	Task  map[string]Task `json:"task"`
}

// Config represents nomad task config
type Config struct {
	Args    []string  `json:"args,omitempty"`
	Image   string    `json:"image"`
	PortMap []PortMap `json:"port_map,omitempty"`
}

// PortMap represents nomad task port map
type PortMap struct {
	ServicePort int `json:"service_port"`
}

// Service represents nomad task service
type Service struct {
	Name        string   `json:"name,omitempty"`
	Port        string   `json:"port,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	AddressMode string   `json:"address_mode,omitempty"`
}

// Task represents nomad task
type Task struct {
	Config    []Config `json:"config"`
	Driver    string   `json:"driver"`
	Resources []struct {
		CPU     int `json:"cpu"`
		Memory  int `json:"memory"`
		Network []struct {
			Mbits int `json:"mbits"`
			Port  []struct {
				ServicePort []struct {
				} `json:"service_port"`
			} `json:"port"`
		} `json:"network"`
	} `json:"resources,omitempty"`
	Service []Service         `json:"service,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// Job represents nomad job
type Job struct {
	Datacenters []string         `json:"datacenters"`
	Type        string           `json:"type"`
	Group       map[string]Group `json:"group"`
}

func convertToNomadJob(name string, service DockerComposeService) (Job, error) {
	//portMaps, err := getPortMaps(service.Environment)
	//if err != nil {
	//	return Job{}, err
	//}
	//services := make([]Service, len(portMaps))
	//for _, portMap := range portMaps {
	//	services = append(services, Service{Name: name, Port: strconv.Itoa(portMap.ServicePort)})
	//}

	envs := environmentsAsMap(service.Environment)
	task := Task{
		Driver: "docker",
		Config: []Config{
			{
				Image: service.Image,
			},
		},
		Service: []Service{{Name: name, AddressMode: "driver"}},
		Env:     envs,
	}
	group := Group{Count: 1, Task: map[string]Task{name: task}}
	job := Job{Datacenters: []string{"dc1"}, Type: "service", Group: map[string]Group{name: group}}

	return job, nil
}

func environmentsAsMap(envs []string) map[string]string {
	acc := map[string]string{}
	for _, env := range envs {
		temp := strings.Split(env, "=")
		acc[temp[0]] = temp[1]
	}
	return acc
}

func getPortMaps(configs []string) ([]PortMap, error) {
	for _, config := range configs {
		if strings.HasPrefix(config, "PORT=") {
			servicePort, err := strconv.Atoi(strings.TrimPrefix(config, "PORT="))
			if err != nil {
				return nil, err
			}
			return []PortMap{{ServicePort: servicePort}}, nil
		}
	}
	return []PortMap{}, nil
}
