package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	dat, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	var dockerFileContent = DockerComposeFileContent{}
	err = yaml.Unmarshal(dat, &dockerFileContent)
	if err != nil {
		panic(err)
	}
	if len(dockerFileContent.Services) == 0 {
		panic("Empty docker-compose file")
	}
	for name, service := range dockerFileContent.Services {
		job, err := convertToNomadJob(name, service)
		if err != nil {
			panic(err)
		}
		err = printJob(name, map[string]map[string]Job{"job": map[string]Job{name: job}})
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Generated %s.nomad file", name))
	}
}

func printJob(name string, job map[string]map[string]Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("unable to marshal %v", err)
	}

	ast, err := hcl.ParseBytes(data)
	if err != nil {
		return fmt.Errorf("unable to parse as hcl %v", err)
	}

	fo, err := os.Create(fmt.Sprintf("%s.nomad", name))
	if err != nil {
		return err
	}

	err = printer.Fprint(fo, ast.Node)
	if err != nil {
		return fmt.Errorf("unable to write as hcl %v", err)
	}
	return nil
}
