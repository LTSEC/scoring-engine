package config

import (
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type Yaml struct {
	WebScore int    `yaml:"webscore"`
	WebIP    string `yaml:"webIP"`
	WebDir   string `yaml:"webDir"`

	FtpScore int               `yaml:"ftpScore"`
	FtpIP    string            `yaml:"ftpIP"`
	FtpCreds map[string]string `yaml:"ftpCreds"`

	TeamScores map[string]int `yaml:"teamScores"`

	SshIP    string            `yaml:"sshIP"`
	PortNum  int               `yaml:"portNum"`
	SshCreds map[string]string `yaml:"sshCreds"`
}

func Parse() Yaml {

	var yamlPath = "../tests/example.yaml"

	file, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var config Yaml

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
