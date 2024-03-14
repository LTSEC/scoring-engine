package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Test struct {
	WebScore int    `yaml:"webscore"`
	WebIP    string `yaml:"webIP"`
	WebDir   string `yaml:"webDir"`

	FtpScore int               `yaml:"ftpScore"`
	FtpIP    string            `yaml:"ftpIP"`
	FtpCreds map[string]string `yaml:"ftpCreds"`

	TeamScores map[string]int `yaml:"teamScores"`
}

func main() {

	var yamlPath = "../tests/example.yaml"

	file, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var me Test

	if err := yaml.Unmarshal(file, &me); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", me)
}
